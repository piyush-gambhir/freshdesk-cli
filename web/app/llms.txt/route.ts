import { source } from '@/lib/source';
import { llms } from 'fumadocs-core/source';
import { site } from '@/lib/site';
import { siteUrl } from '@/lib/shared';
import { getOtherSuiteProjects } from '@/lib/suite';

export const revalidate = false;

export async function GET() {
  // index() returns a Promise since fumadocs-core 16.15.17.
  const index = (await llms(source).index()).replace(/\]\((\/[^)]+)\)/g, (_match, path: string) => `](${siteUrl}${path})`);
  const relatedSites = getOtherSuiteProjects(site.repo)
    .map(({ name, href }) => `- ${name}: ${href}`)
    .join('\n');
  const intro =
    'Freshdesk CLI (binary: freshdesk) is an independent, unofficial command-line interface for the Freshdesk REST API. Named commands cover 266 documented API operation variants across tickets, conversations, contacts, companies, knowledge bases, forums, and administration. Agents should use --read-only --no-input -o json for reads, discover capabilities with freshdesk commands -o json, preview writes with --dry-run, and execute only authorized changes with --expect-plan HASH --yes. The approved live write test covered one marked ticket; other writes remain fixture-tested. It is not affiliated with Freshworks.';
  return new Response(
    `${intro}\n\n${index}\n\n## Related CLI sites\n\n${relatedSites}\n`,
  );
}
