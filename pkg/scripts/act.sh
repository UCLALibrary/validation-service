#! /bin/bash

#
# A script to use ACT to run (and test) our GitHub Actions.
#
# Requirements:
#   * ACT (https://nektosact.com/)
#   * git (https://docs.github.com/en/get-started/git-basics/set-up-git)
#   * yq (https://mikefarah.gitbook.io/yq)
#

# Warn about one of our system requirements not being able to be found
function warn_not_found {
  echo "Warning: $1 was not found, but must be installed to use this script"
  exit 1
}

# For releases, we need to have a tag in the local repo and to generate a tag event
function get_tag_event_file {
  EVENT_FILE="/tmp/tag_event.json"

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

  if ! git describe --tags "$(git rev-list --tags --max-count=1)" > /dev/null 2>&1; then
    echo "No git tags exist in your local repo; either fetch them or create one"
    exit 1
  fi

  # Get the latest commit's SHA to run a release
  LATEST_TAG=$(git describe --tags "$(git rev-list --tags --max-count=1)")
  COMMIT_SHA=$(git rev-list -n 1 "$LATEST_TAG")

  # Define the two configuration files we require for the script to work
  FILES=(
    "$HOME/.act-secrets"
    "$HOME/.act-variables"
  )

  for FILE in "${FILES[@]}"; do
    if [ ! -e "$FILE" ]; then
      echo "Configuration file \"$FILE\" was not found; please create it before running this script"
      exit 1
    fi
  done

  # Get the GITHUB_USER value
  source "$HOME/.act-variables"

  # Create a tag event that can be supplied to the ACT release job
  yq ".ref = \"refs/tags/$LATEST_TAG\" | .after = \"$COMMIT_SHA\" | .pusher.name = \"$GITHUB_USER\"" \
    testdata/tag_event.json > "$EVENT_FILE"

  # Return the location of the newly created tag_event.json
  echo "$EVENT_FILE"
}

# Check that there is a requested CI action to run
if [ -z "$1" ]; then
  echo "Supply either 'release', 'build', or 'nightly' to run a CI action"
  exit 1
fi

# Confirm the argument is one of our supported CI actions
if [ "$1" = "release" ]; then
  ACTION="release"
elif [ "$1" = "build" ]; then
  ACTION="build"
elif [ "$1" = "nightly" ]; then
  ACTION="nightly"
else
  echo "Invalid argument: $1. Please use 'release', 'build', or 'nightly'"
  exit 1
fi

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

# If we're running a release we need to generate a tag event, otherwise we run a basic action
if [ "$ACTION" = "release" ]; then
  EVENT_FILE="$(get_tag_event_file)"
  $ACT --secret-file ~/.act-secrets --var-file ~/.act-variables -e "$EVENT_FILE" -j $ACTION
else
  $ACT --secret-file ~/.act-secrets --var-file ~/.act-variables -j $ACTION
fi
