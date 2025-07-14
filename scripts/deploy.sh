#!/bin/bash
set -e

export $(cat .env.github | xargs)

# Step 1: Build the static site
go run cmd/webgen/main.go

# Step 2: Prepare temporary gh-pages worktree
rm -rf /tmp/gh-pages
git worktree add /tmp/gh-pages gh-pages

# Step 3: Clear existing contents
rm -rf /tmp/gh-pages/*

# Step 4: Copy static files
cp -r dist/* /tmp/gh-pages

# Step 5: Commit and push
cd /tmp/gh-pages

COMMIT_HASH=$(git rev-parse --short HEAD)
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")

if git diff --cached --quiet; then
  echo "⚠️  No changes to commit"
else
  git commit -m "chore: deploy site from ${COMMIT_HASH} (${TIMESTAMP}) [skip ci]"
  git push origin gh-pages
fi

# Step 6: Cleanup
cd -
git worktree remove /tmp/gh-pages
