import type { Metadata } from 'next';
import Link from 'next/link';
import { LegalPage } from '@/components/legal-page';
import { createPageMetadata } from '@/lib/metadata';

export const metadata: Metadata = createPageMetadata({
  title: 'Privacy Policy',
  description: 'How Freshdesk CLI handles credentials, API requests, and the public documentation site.',
  path: '/privacy',
});

export default function PrivacyPage() {
  return (
    <LegalPage title="Privacy Policy" effective="October 29, 2026">
      <p className="legal-page__lede">
        Freshdesk CLI is an independent, unofficial open-source program that runs on your computer.
        The CLI has no maintainer-operated API backend, analytics, telemetry, or advertising.
      </p>
      <h2>Credentials and local files</h2>
      <p>
        API keys are stored in the OS keychain by default. Choosing <code>--insecure-storage</code>
        explicitly saves a plaintext secrets file, with mode <code>0600</code> on Unix.
        Windows relies on the user&apos;s filesystem ACLs. Profile metadata is stored separately in
        <code> ~/.config/freshdesk-cli/config.yaml</code>. Environment authentication does not require
        saving a credential. Configuration paths respect <code>XDG_CONFIG_HOME</code> and <code>FRESHDESK_CONFIG</code>.
      </p>
      <h2>Network requests</h2>
      <p>
        Commands send requests directly to the helpdesk origin you configure, using your key to
        access records or perform the actions you request. The key owner is the acting Freshdesk agent.
        Your organization&apos;s Freshdesk configuration may trigger automations, notifications, or integrations.
      </p>
      <p>
        Explicit update checks and self-updates contact GitHub release pages and download hosts.
        The installer also downloads from GitHub. There are no background update checks.
        Offline discovery and dry-run previews do not contact Freshdesk.
      </p>
      <h2>Records and output</h2>
      <p>
        API records are returned to your terminal or the files you choose through shell redirection.
        The CLI does not send those records to the maintainer. You control access to your terminal,
        shell history, logs, exported files, and any agent or CI system running the commands.
      </p>
      <h2>Documentation website</h2>
      <p>
        This website serves static public documentation through Cloudflare. It does not request a
        Freshdesk API key, connect to your helpdesk, or include a third-party analytics SDK.
        Theme preferences are stored locally in your browser. Hosting providers may process request
        metadata under their own privacy policies.
      </p>
      <h2>Third parties and contact</h2>
      <p>
        Use of Freshdesk, GitHub, and Cloudflare is also subject to those providers&apos; policies.
        Report suspected vulnerabilities through the repository&apos;s private reporting channel.
        For other questions, see the <Link href="/contact">contact page</Link>.
      </p>
    </LegalPage>
  );
}
