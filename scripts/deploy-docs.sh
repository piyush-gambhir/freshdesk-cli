#!/usr/bin/env bash
# Publish the checked-in docs from an up-to-date main checkout.
set -euo pipefail

cd "$(dirname "$0")/.."
if [ "$(git branch --show-current)" != main ]; then
  echo "Deploy the documentation from main after required checks pass." >&2
  exit 1
fi
if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
  echo "Commit tracked changes before deploying the documentation." >&2
  exit 1
fi
git fetch origin main
if [ "$(git rev-parse HEAD)" != "$(git rev-parse origin/main)" ]; then
  echo "Update main to origin/main before deploying the documentation." >&2
  exit 1
fi
cd web
pnpm install --frozen-lockfile
pnpm types:check
pnpm build:cloudflare
pnpm test:search
pnpm deploy:cloudflare
