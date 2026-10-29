# Maintaining the repository

## Required checks and repository settings

The default branch is `main`. Its ruleset requires pull requests, signed commits, linear history, resolved review threads, and these checks against the current base:

- `test (ubuntu-24.04)`
- `test (macos-latest)`
- `test (windows-latest)`
- `release-check`
- `web`
- `analyze (go)`
- `analyze (javascript-typescript)`
- `analyze (python)`
- `analyze (actions)`

The ruleset blocks branch deletion and force pushes. Repository administrators retain the explicit bypass used by the other CLI repositories. Release tags matching `v*` cannot be updated or deleted under the tag ruleset. Merge commits are disabled; squash and rebase merges are enabled. Auto-merge and automatic deletion of merged branches are enabled.

GitHub Actions defaults to read permissions and cannot approve pull requests. All actions must be pinned to full commit SHAs. Secret scanning, push protection, Dependabot alerts/security updates, and private vulnerability reporting are enabled. Dependabot schedules Go and Actions updates daily and website dependency updates weekly. The auto-merge workflow only enables merging minor and patch dependency updates; all required checks still apply. CodeQL scans Go, JavaScript/TypeScript, Python, and Actions. Scorecard runs on main and weekly.

## Releases

Change `cli-go/VERSION` in a pull request. After merge, the Release workflow validates the version, tests the source, creates the tag at that commit, generates six platform archives and their SBOMs, uploads a draft release, attests the checksum subjects, and publishes it. Version 0.1.2 and later releases are immutable. Earlier releases are not retroactively locked.

Do not push version tags manually. If a release fails, rerun the original workflow run so its commit stays fixed. The workflow skips an already published version, refuses a tag that points to another commit, and avoids marking an older release or a prerelease as latest. Manual dispatch is available on `main` for recovery.

Before merging, run the affected checks and `make test vet build docs`. CI also runs staticcheck, govulncheck, three-platform race tests, source-install verification, registry regeneration, website checks, and snapshot release packaging.

## Website

The public site is https://projects.piyushgambhir.com/freshdesk-cli. It uses the approved CLI-suite template with a coral accent, Next.js static export, Fumadocs, and a dedicated Cloudflare Worker. Its only route is `projects.piyushgambhir.com/freshdesk-cli*`; other project routes are unchanged.

Guide content lives in `web/content/docs/`. Builds regenerate reference pages from `docs/commands.md`, `docs/api-coverage.md`, `docs/live-verification.md`, and this file. Run these checks in `web/`:

```sh
pnpm install --frozen-lockfile
pnpm types:check
pnpm build:cloudflare
pnpm test:search
```

Verify desktop and mobile pages, both themes, navigation, search, copy actions, 404s, metadata, agent text exports, and reduced-motion behavior before publishing visual changes. From an up-to-date `main`, deploy with `scripts/deploy-docs.sh`. Cloudflare credentials stay in the local Wrangler credential store or a deployment environment; they are not committed to the repository.

The site does not load a Freshdesk API key or customer records. It serves public documentation only. Do not add live tenant data to examples, screenshots, fixtures, or build inputs.
