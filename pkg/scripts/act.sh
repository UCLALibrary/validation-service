#! /bin/bash

#
# A script to use ACT to run (and test) our GitHub Actions.
#
# Requirements:
#   * ACT (https://nektosact.com/)
#   * git (https://docs.github.com/en/get-started/git-basics/set-up-git)
#   * yq (https://mikefarah.gitbook.io/yq)
#

# Confirm the first argument is one of our supported CI actions
case "$1" in
  build|prerelease|release|nightly)
    ACTION="$1"
    ;;
  *)
    echo "Invalid argument: $1. Please supply 'JOB=build', 'JOB=prerelease', 'JOB=release', or 'JOB=nightly'"
    exit 1
    ;;
esac

# Confirm that there were two arguments passed in; the second should be from the Makefile itself
if [ -n "$2" ]; then
  SERVICE_NAME="$2"
else
  echo "Script is confused; the second arg should have been SERVICE_NAME (supplied via Makefile)"
  exit 1
fi

# Define the two configuration files we require for the script to work
FILES=(
  "$HOME/.act-secrets"
  "$HOME/.act-variables"
)

# Confirm the config files exist so that we can read from them
for FILE in "${FILES[@]}"; do
  if [ ! -e "$FILE" ]; then
    echo "Configuration file \"$FILE\" was not found; please create it before running this script"
    exit 1
  fi
done

# Warn about one of our system requirements not being able to be found
function warn_not_found {
  echo "Warning: $1 was not found, but must be installed to use this script"
  exit 1
}

# For releases, we need to have a tag in the local repo and to generate a (pre)release event
function release_event {
  EVENT_FILE="/tmp/release_event.json"
  PRERELEASE=$([ -n "$1" ] && echo true || echo false)

  # Check to confirm yq is installed on our system
  if ! command -v yq > /dev/null 2>&1; then
    # shellcheck disable=SC2016
    warn_not_found '`yq`'
  fi

  # Check to confirm git is installed on our system
  if ! command -v git > /dev/null 2>&1; then
    # shellcheck disable=SC2016
    warn_not_found '`git`'
  fi

  # Check that there are tags in our local git repo
  if ! git describe --tags "$(git rev-list --tags --max-count=1)" > /dev/null 2>&1; then
    echo "No git tags exist in your local repo; either fetch them or create one"
    exit 1
  fi

  # Get the latest commit's tag and SHA to perform a release
  LATEST_TAG=$(git describe --tags "$(git rev-list --tags --max-count=1)")
  COMMIT_SHA=$(git rev-list -n 1 "$LATEST_TAG")

  # Get the GITHUB_USER and DOCKER_REGISTRY_ACCOUNT values
  source "$HOME/.act-variables"

  TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  # Create a (pre)release event that can be supplied to the ACT release job
  yq ".release.tag_name = \"$LATEST_TAG\" | .release.target_commitish = \"$COMMIT_SHA\" |
      .sender.login = \"$GITHUB_USER\" | .release.author.login = \"$GITHUB_USER\" | .release.name = \"v$LATEST_TAG\" |
      .repository.name = \"$SERVICE_NAME\" | .repository.full_name = \"$DOCKER_REGISTRY_ACCOUNT/$SERVICE_NAME\" |
      .repository.owner.login = \"$DOCKER_REGISTRY_ACCOUNT\" | .release.body = \"Automated release of v$LATEST_TAG\" |
      .release.created_at = \"$TIMESTAMP\" | .release.published_at = \"$TIMESTAMP\" |
      .configmap.data.VERSION = \"$LATEST_TAG\" |
      .release.prerelease = $PRERELEASE" testdata/release_event.json > "$EVENT_FILE"

  # Return the location of the newly created release_event.json
  echo "$EVENT_FILE"
}

# Check to confirm ACT is installed on our system
if command -v act > /dev/null 2>&1; then
  ACT='act'
else
  if command -v gh > /dev/null 2>&1; then
    if gh extension list |grep -q 'nektos/gh-act'; then
      ACT='gh act'
    else
      # shellcheck disable=SC2016
      not_found '`gh act`'
    fi
  else
    # shellcheck disable=SC2016
    not_found '`act` or `gh act`'
  fi
fi

# If we're running a (pre)release we need to generate a release event, otherwise we run a basic action
if [ "$ACTION" = "release" ]; then
  $ACT --secret-file ~/.act-secrets --var-file ~/.act-variables -e "$(release_event)" release
elif [ "$ACTION" = "prerelease" ]; then
  $ACT --secret-file ~/.act-secrets --var-file ~/.act-variables -e "$(release_event \"pre\")" release
else
  $ACT --secret-file ~/.act-secrets --var-file ~/.act-variables -j "$ACTION"
fi
