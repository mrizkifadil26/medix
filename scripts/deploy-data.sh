#!/bin/bash
set -e

BRANCH=main
FOLDER=data

echo "🚀 Deploying updated JSON files..."

git add $FOLDER/*.json
git commit -m "Update data: $(date +'%F %T')" || echo "No changes to commit"
git push origin $BRANCH

echo "✅ Pushed to GitHub."
