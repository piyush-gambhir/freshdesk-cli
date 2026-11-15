import {
  Bot,
  Ticket,
  KeyRound,
  ListChecks,
  BookOpen,
  ContactRound,
  type LucideIcon,
} from 'lucide-react';

export interface Feature {
  icon: LucideIcon;
  title: string;
  body: string;
  docsLink?: {
    label: string;
    href: string;
  };
}

export interface SiteConfig {
  name: string;
  binary: string;
  repo: string;
  tagline: string;
  description: string;
  installCommand: string;
  features: Feature[];
  exampleTitle: string;
  example: string;
  compatible?: string[];
  featuresTitle?: string;
  featuresSubtitle?: string;
  ctaBody?: string;
  accent?: string;
  accentName?: string;
  accentHex?: string;
}

export const site: SiteConfig = {
  name: 'Freshdesk CLI',
  binary: 'freshdesk',
  repo: 'piyush-gambhir/freshdesk-cli',
  tagline: 'Freshdesk CLI',
  description:
    'An independent, unofficial command-line interface for Freshdesk. Work with tickets, customer records, knowledge bases, and account administration using structured output and explicit write controls.',
  accent: 'oklch(0.71 0.16 30)',
  accentName: 'coral',
  accentHex: '#f57663',
  installCommand:
    'curl -fsSL https://raw.githubusercontent.com/piyush-gambhir/freshdesk-cli/main/install.sh | sh',
  featuresTitle: 'Commands and safeguards',
  featuresSubtitle: '266 documented API operation variants, with coverage and live verification recorded separately.',
  features: [
    {
      icon: Ticket,
      title: 'Tickets and conversations',
      body: 'List, search, and inspect tickets. Preview ticket changes, private notes, replies, and attachments before sending them.',
      docsLink: { label: 'List, search, and inspect tickets', href: '/docs/commands/tickets' },
    },
    {
      icon: ContactRound,
      title: 'Customer records',
      body: 'Read contacts and companies, inspect their fields, and send native JSON payloads for account-specific data.',
      docsLink: { label: 'Read contacts and companies', href: '/docs/commands/customers' },
    },
    {
      icon: BookOpen,
      title: 'Knowledge base and forums',
      body: 'Navigate categories, folders, articles, translations, and community forums through named commands.',
      docsLink: { label: 'Navigate categories, folders, articles', href: '/docs/commands/knowledge-base' },
    },
    {
      icon: ListChecks,
      title: 'Account administration',
      body: 'Inspect agents, groups, ticket fields, automations, SLAs, and other administration APIs available to your account.',
      docsLink: { label: 'Inspect agents, groups, ticket fields', href: '/docs/commands/administration' },
    },
    {
      icon: KeyRound,
      title: 'Controlled writes',
      body: 'Read-only mode, offline request previews, and hashes that bind approval to the exact payload and attachment contents.',
      docsLink: { label: 'Read-only mode, offline request previews', href: '/docs/safety' },
    },
    {
      icon: Bot,
      title: 'Agent discovery',
      body: 'Discover commands offline, resolve arguments, and use JSON output without loading configuration or touching the keychain.',
      docsLink: { label: 'Discover commands offline', href: '/docs/agents' },
    },
  ],
  exampleTitle: 'Read and preview',
  example: [
    '# Save a verified key in the OS keychain',
    'freshdesk auth login --domain your-team',
    '# Read without allowing writes',
    'freshdesk --read-only agents me -o json',
    'freshdesk --read-only tickets list --per-page 10 -o json',
    '# Inspect a request without sending it',
    'freshdesk tickets update 123 --status pending --dry-run -o json',
  ].join('\n'),
  compatible: ['Tickets', 'Conversations', 'Contacts', 'Companies', 'Knowledge base', 'Forums', 'Automations', 'Administration'],
  ctaBody: 'Install a checksummed binary for macOS or Linux, or download a Windows archive from GitHub Releases.',
};
