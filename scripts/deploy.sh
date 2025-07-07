#!/bin/bash
set -e

export $(cat .env.github | xargs)

# Step 1: Build the static site
go run cmd/builder/main.go

# Step 2: Prepare temporary gh-pages worktree
rm -rf /tmp/gh-pages
git worktree add /tmp/gh-pages gh-pages

# Step 3: Clear existing contents
rm -rf /tmp/gh-pages/*

# Step 4: Copy static files
cp -r dist/* /tmp/gh-pages

# Step 5: Commit and push
cd /tmp/gh-pages
git add .
git commit -m "Deploy to GitHub Pages" || echo "No changes to commit"
git push origin gh-pages

# Step 6: Cleanup
cd -
git worktree remove /tmp/gh-pages
