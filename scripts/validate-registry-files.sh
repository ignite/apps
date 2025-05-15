#!/bin/bash

set -euo pipefail

FOLDER="_registry" # Change this to your target folder

# Check if ignite is installed
if ! command -v ignite >/dev/null 2>&1; then
  curl https://get.ignite.com/cli\! | bash
fi

if ! ignite app list | grep -q 'appregistry'; then
  ignite app install -g ./appregistry
fi

# Find and validate each JSON file
find "$FOLDER" -type f -name '*.json' | while read -r file; do
  if [[ $(basename "$file") != "registry.json" ]]; then
    echo "Validating $file..."
    ignite appregistry validate "$file" --branch "$(git rev-parse --abbrev-ref HEAD)"
  fi
done

echo "✅ All files validated."
