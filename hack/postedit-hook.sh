#!/bin/bash
set -euo pipefail

if [ $# -eq 0 ]; then
    echo "Error: No file path provided."
    exit 1
fi

file_path="$1"
if [ ! -f "$file_path" ]; then
    echo "Error: File '$file_path' does not exist."
    exit 1
fi
file_path_abs=$(realpath "$file_path")

if [ -z "${H132_LWS_DIR:-}" ]; then
    echo "Error: H132_LWS_DIR environment variable is not set."
    exit 1
fi
H132_LWS_DIR=$(realpath "$H132_LWS_DIR")

# Ensure the file is within the git repository
if [[ "$file_path_abs" != "$H132_LWS_DIR"* ]]; then
    echo "Error: '$file_path' is not within the git repository at '$H132_LWS_DIR'."
    exit 1
fi

git add "${file_path_abs#$H132_LWS_DIR/}"
git commit -m "Automated commit by h132 post-edit hook script."

if git push; then
    echo "Changes have been pushed to the remote repository."
else
    echo "Error: Failed to push changes to the remote repository."
    exit 1
fi

exit 0
