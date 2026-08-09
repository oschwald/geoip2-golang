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

release_commit=$(git rev-parse --verify HEAD)

release_heading=$(awk '/^## / { print; exit }' CHANGELOG.md)
if [ "$release_heading" = "## Unreleased" ]; then
    die "replace the Unreleased heading with a version and release date"
fi

semver_number='(0|[1-9][0-9]*)'
semver_identifier='(0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)'
heading_pattern="^## $semver_number\.$semver_number\.$semver_number(-$semver_identifier(\.$semver_identifier)*)? - [0-9]{4}-[0-9]{2}-[0-9]{2}$"
if [ -z "$release_heading" ] || [[ ! $release_heading =~ $heading_pattern ]]; then
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

if ! gh auth status >/dev/null 2>&1; then
    die "gh is not authenticated"
fi

repository=$(
    gh repo view "$(git remote get-url origin)" --json nameWithOwner --jq .nameWithOwner
) || die "could not resolve the origin GitHub repository"

if git show-ref --verify --quiet "refs/tags/$tag"; then
    die "tag $tag already exists locally"
fi

remote_tag=$(git ls-remote origin "refs/tags/$tag") ||
    die "could not check whether tag $tag exists on origin"
if [ -n "$remote_tag" ]; then
    die "tag $tag already exists on origin"
fi

release_response=$(
    gh api --include --silent "repos/$repository/releases/tags/$tag" 2>/dev/null || true
)
release_http_status=$(printf '%s\n' "$release_response" | awk 'NR == 1 { print $2 }')
case "$release_http_status" in
    200) die "release $tag already exists" ;;
    404) ;;
    *) die "could not verify whether release $tag exists" ;;
esac

echo "Running release checks..."
if [ "$(go env GOHOSTOS)" != "linux" ] || [ "$(go env GOHOSTARCH)" != "amd64" ]; then
    die "release checks require a Linux/amd64 host"
fi
if [ "$(go env GOOS)" != "linux" ] || [ "$(go env GOARCH)" != "amd64" ]; then
    die "release checks require GOOS=linux and GOARCH=amd64"
fi
if [ "$(go env CGO_ENABLED)" != "1" ]; then
    die "release checks require CGO_ENABLED=1"
fi
golangci-lint fmt
go mod tidy
go mod verify
go generate ./...
go test ./...
go test -race ./...
CGO_ENABLED=0 GOOS=linux GOARCH=386 go test ./...
golangci-lint run

if [ -n "$(git status --porcelain)" ]; then
    echo "Release checks modified the working directory:" >&2
    git status --short >&2
    git diff >&2
    exit 1
fi

if [ "$(git rev-parse --verify HEAD)" != "$release_commit" ]; then
    die "HEAD changed during release checks"
fi

echo
echo "Version: $version"
echo
echo "Release notes:"
printf '%s\n' "$notes"
echo

read -r -p "Push $current_branch and create $tag in $repository? [y/N] " should_release
if [[ ! $should_release =~ ^[Yy]$ ]]; then
    echo "Aborting"
    exit 1
fi

notes_file=$(mktemp "${TMPDIR:-/tmp}/geoip2-release-notes.XXXXXX")
trap 'rm -f -- "$notes_file"' EXIT
printf '%s\n' "$notes" >"$notes_file"

git push origin "$release_commit:refs/heads/$current_branch"

remote_tag=$(git ls-remote origin "refs/tags/$tag") ||
    die "could not recheck whether tag $tag exists on origin"
if [ -n "$remote_tag" ]; then
    die "tag $tag appeared on origin during release checks"
fi

gh release create \
    --repo "$repository" \
    --target "$release_commit" \
    --title "$version" \
    --notes-file "$notes_file" \
    "$tag"
