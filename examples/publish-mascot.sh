#!/usr/bin/env bash
# Publish only generated artwork, leaving the checked-out development branch alone.
set -eu
: "${SVG_PATH:?Set SVG_PATH to the generated SVG}"
: "${BADGE_REPO:?Set BADGE_REPO to a new temporary worktree path}"

if git show-ref --verify --quiet refs/remotes/origin/badge-artifact; then
  git worktree add --detach "$BADGE_REPO" origin/badge-artifact
else
  git worktree add --detach "$BADGE_REPO"
  git -C "$BADGE_REPO" switch --orphan badge-artifact
fi
cp "$SVG_PATH" "$BADGE_REPO/mascot.svg"
cd "$BADGE_REPO"
git add mascot.svg
if ! git diff --cached --quiet; then
  git -c user.name='github-actions[bot]' \
    -c user.email='41898282+github-actions[bot]@users.noreply.github.com' \
    commit -m 'Update build mascot [skip ci]'
  git push origin HEAD:refs/heads/badge-artifact
fi
