# freshdesk-cli docs site

The site at [projects.piyushgambhir.com/freshdesk-cli](https://projects.piyushgambhir.com/freshdesk-cli):
Next.js with [Fumadocs](https://fumadocs.dev), exported as static files and served by a Cloudflare Worker.

```bash
pnpm install
pnpm dev                 # http://localhost:3000/freshdesk-cli
pnpm types:check
pnpm build:cloudflare    # static export into .cloudflare/assets/freshdesk-cli
pnpm test:search         # checks the exported search index
```

- Guide pages: `content/docs/`. Site name, hero, features, and accent: `lib/site.ts`.
- `content/docs/reference/` is generated from the repository's `docs/` by `scripts/sync-reference.mjs`,
  which every `dev`, `build`, and `types:check` runs first. Edit `docs/` (or the Go command help) instead.
- `lib/suite.ts` lists related project sites for machine-readable metadata and `llms.txt`.
- Deploy with `../scripts/deploy-docs.sh` using the existing Wrangler session or `CLOUDFLARE_API_TOKEN`.
