#!/usr/bin/env bash

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

die() {
    echo "Error: $*" >&2
    exit 1
}

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        die "$1 is not installed or not in PATH"
    fi
}

for command in gh git go golangci-lint; do
    require_command "$command"
done

current_branch=$(git branch --show-current)
if [ -z "$current_branch" ]; then
    die "releases cannot be created from a detached HEAD"
fi

echo "Fetching from origin..."
git fetch origin
git submodule update --init --recursive

if ! git merge-base --is-ancestor origin/main HEAD; then
    die "the current branch does not contain the latest origin/main"
fi

if [ -n "$(git status --porcelain)" ]; then
    die "the working directory is not clean"
fi

if grep -q '^## Unreleased$' CHANGELOG.md; then
    die "replace the Unreleased heading with a version and release date"
fi

heading_pattern='^## [0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z][0-9A-Za-z.-]*)? - [0-9]{4}-[0-9]{2}-[0-9]{2}$'
release_heading=$(grep -m1 -E "$heading_pattern" CHANGELOG.md || true)
if [ -z "$release_heading" ]; then
    die "CHANGELOG.md does not start with a release heading such as ## 2.3.0 - 2026-08-08"
fi

heading_regex='^##[[:space:]]([^[:space:]]+)[[:space:]]-[[:space:]]([0-9]{4}-[0-9]{2}-[0-9]{2})$'
if [[ ! $release_heading =~ $heading_regex ]]; then
    die "could not parse release heading: $release_heading"
fi

version=${BASH_REMATCH[1]}
release_date=${BASH_REMATCH[2]}
today=$(date +%Y-%m-%d)
tag="v$version"

if [ "$release_date" != "$today" ]; then
    die "release date $release_date is not today ($today)"
fi

notes=$(
    awk -v heading="$release_heading" '
        $0 == heading { found = 1; next }
        found && /^## / { exit }
        found { print }
    ' CHANGELOG.md | sed '/./,$!d'
)
if [ -z "$notes" ]; then
    die "the $version changelog section has no release notes"
fi

module_path=$(go list -m -f '{{.Path}}')
major=${version%%.*}
case "$major" in
    0 | 1)
        if [[ $module_path =~ /v[0-9]+$ ]]; then
            die "tag $tag does not match module path $module_path"
        fi
        ;;
    *)
        if [[ $module_path != */v"$major" ]]; then
            die "tag $tag does not match module path $module_path"
        fi
        ;;
esac

if git show-ref --verify --quiet "refs/tags/$tag"; then
    die "tag $tag already exists locally"
fi
if git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1; then
    die "tag $tag already exists on origin"
fi

if ! gh auth status >/dev/null 2>&1; then
    die "gh is not authenticated"
fi

echo "Running release checks..."
golangci-lint fmt
go mod verify
go generate ./...
go test ./...
go test -race ./...
golangci-lint run

if [ -n "$(git status --porcelain)" ]; then
    echo "Release checks modified the working directory:" >&2
    git status --short >&2
    git diff >&2
    exit 1
fi

release_commit=$(git rev-parse HEAD)

echo
echo "Version: $version"
echo
echo "Release notes:"
echo "$notes"
echo

read -r -p "Create $tag from $current_branch and push to origin? [y/N] " should_release
if [[ ! $should_release =~ ^[Yy]$ ]]; then
    echo "Aborting"
    exit 1
fi

git push --atomic origin \
    "$release_commit:refs/heads/$current_branch" \
    "$release_commit:refs/tags/$tag"

notes_file=$(mktemp "${TMPDIR:-/tmp}/geoip2-release-notes.XXXXXX")
trap 'rm -f -- "$notes_file"' EXIT
printf '%s\n' "$notes" >"$notes_file"

gh release create \
    --verify-tag \
    --title "$version" \
    --notes-file "$notes_file" \
    "$tag"
