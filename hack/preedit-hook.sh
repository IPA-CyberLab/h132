#!/bin/bash
set -euo pipefail

if [ $# -eq 0 ]; then
    echo "Error: No file path provided."
    exit 1
fi
file_path="$1"
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

if [ ! -f "$file_path_abs" ]; then
    echo "Info: File '$file_path' does not exist yet. Skipping git checks."
    exit 0
fi

# Get the file path relative to the git repository
file_rel_path="${file_path_abs#$H132_LWS_DIR/}"

cd "$H132_LWS_DIR"

if ! git ls-files --error-unmatch "$file_rel_path" > /dev/null 2>&1; then
    echo "Error: '$file_rel_path' is not tracked by git in repository '$H132_LWS_DIR'."
    exit 1
fi

if [[ -n $(git status --porcelain "$file_rel_path") ]]; then
    echo "Error: '$file_rel_path' has uncommitted changes. Please commit your changes before proceeding."
    exit 1
fi

exit 0
