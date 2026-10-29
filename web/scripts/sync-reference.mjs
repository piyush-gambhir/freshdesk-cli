// Generates the reference pages from the repository's docs/ (written by `make docs`),
// so the site always matches the CLI. The output is gitignored; builds run this first.
import { mkdir, readFile, writeFile } from 'node:fs/promises';
import path from 'node:path';

const root = path.resolve(import.meta.dirname, '..');
const blob = 'https://github.com/piyush-gambhir/freshdesk-cli/blob/main';
const outDir = path.join(root, 'content/docs/reference');

// Repository docs that have their own page on this site, with the link text to show.
const siteLinks = {
  'docs/api-coverage.md': ['/docs/reference/api-coverage', 'API coverage'],
  'docs/commands.md': ['/docs/reference/commands', 'Command reference'],
  'docs/live-verification.md': ['/docs/reference/live-verification', 'Live verification'],
  'docs/maintaining.md': ['/docs/reference/maintaining', 'Maintaining the repository'],
};

const pages = [
  { from: 'docs/commands.md', slug: 'commands', title: 'Command reference', description: 'Every Freshdesk command, positional ID, and flag, generated from the command tree.' },
  { from: 'docs/api-coverage.md', slug: 'api-coverage', title: 'API coverage', description: 'The pinned Freshdesk operation inventory, command mapping, effects, and reference discrepancies.' },
  { from: 'docs/live-verification.md', slug: 'live-verification', title: 'Live verification', description: 'Read verification and the explicitly approved one-ticket lifecycle test, separated from fixture-only coverage.' },
  { from: 'docs/maintaining.md', slug: 'maintaining', title: 'Maintaining the repository', description: 'Required checks, repository safeguards, automatic releases, and website deployment.' },
];

// Points a relative link at this site's page for that doc, or at the file on GitHub.
function rewriteLink(text, target, fromFile) {
  if (/^[a-z]+:/i.test(target) || target.startsWith('#') || target.startsWith('/')) return `[${text}](${target})`;
  const [file, hash = ''] = target.split('#');
  const resolved = path.posix.normalize(path.posix.join(path.posix.dirname(fromFile), file));
  const anchor = hash ? `#${hash}` : '';
  const page = siteLinks[resolved];
  if (!page) return `[${text}](${blob}/${resolved}${anchor})`;
  return `[${text.endsWith('.md') ? page[1] : text}](${page[0]}${anchor})`;
}

// MDX treats { } and < as syntax in prose; code spans and fences stay literal.
function escapeProse(text) {
  return text.replace(/[{}<]/g, (c) => (c === '<' ? '&lt;' : `\\${c}`));
}

function convert(markdown, fromFile) {
  const out = [];
  let fenced = false;
  for (const line of markdown.split('\n')) {
    if (/^\s*```/.test(line)) {
      fenced = !fenced;
      out.push(line);
      continue;
    }
    if (fenced) {
      out.push(line);
      continue;
    }
    const linked = line.replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_m, text, target) => rewriteLink(text, target, fromFile));
    out.push(linked.split(/(`[^`]*`)/).map((part, i) => (i % 2 ? part : escapeProse(part))).join(''));
  }
  return out.join('\n');
}

await mkdir(outDir, { recursive: true });
for (const page of pages) {
  let body = await readFile(path.join(root, '..', page.from), 'utf8');
  body = body.replace(/^# .*\n+/, '');
  const frontmatter = `---\ntitle: ${JSON.stringify(page.title)}\ndescription: ${JSON.stringify(page.description)}\n---\n\n`;
  await writeFile(path.join(outDir, `${page.slug}.mdx`), frontmatter + convert(body, page.from));
}
await writeFile(
  path.join(outDir, 'meta.json'),
  `${JSON.stringify({ title: 'Reference', pages: pages.map((p) => p.slug) }, null, 2)}\n`,
);
console.log(`Synced ${pages.length} reference pages from docs/`);
