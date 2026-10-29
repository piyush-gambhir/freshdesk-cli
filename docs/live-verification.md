# Live verification

Live checks were run against one authorized Freshdesk tenant on October 29, 2026. Credentials and customer records are not included in this repository. The initial read checks enforced `--read-only` and `--no-input`. A subsequent, explicitly approved ticket-lifecycle test used ten scoped mutations and verified cleanup.

The first pass completed 82 successful read checks across tickets, contacts, companies, agent/group administration, forms and fields, knowledge base, community discussions, canned responses, custom-object schemas, automation rules, SLAs, email settings, and account configuration. This is a check count, including some repeated operations; it is not a claim that all API operations were verified live.

## Compatibility corrections

- The admin ticket-field list rejects both `page` and `per_page` with HTTP 400 and `invalid_field`. A request without those parameters returned HTTP 200. The named command now omits pagination and does not advertise `--all`.
- The CSAT survey API can return an `http://` pagination link in an HTTPS response. The CLI validates that the host and path match, extracts only the cursor token, and constructs the next request on the original HTTPS origin. It never follows the HTTP URL or downgrades transport security. Single-page and bounded two-page live checks succeeded after this correction.

Regression tests cover both behaviors. Cross-origin cursor links and non-advancing cursors remain rejected.

## Approved ticket-lifecycle test

One uniquely marked sample ticket was created with an internal requester, Pending status, and Low priority. Every mutation used a request preview and its `--expect-plan` hash. IDs were accepted only from records created during the run, and the ticket marker was checked before subsequent mutations.

All ten mutation requests succeeded:

1. Create the sample ticket.
2. Change its priority to Medium.
3. Restore its original priority.
4. Add a private note with an empty notification-recipient list.
5. Edit that note.
6. Delete that note.
7. Upload a harmless text attachment through a second private note.
8. Delete the attachment.
9. Delete the second note.
10. Soft-delete the sample ticket.

Thirteen read-back assertions verified the ticket marker/requester, priority changes, note privacy and content, uploaded attachment metadata, and removal of the test notes, attachment, and active ticket. No step failed and no mutation was retried. The ticket was soft-deleted, not permanently purged.

Existing tenant automations remained active. The user was informed that ticket creation triggers a backend webhook before authorizing the run. This verifies the scoped CLI flow; it does not establish that production-side integrations produced no events. Public replies, notifications to explicitly selected recipients, and account-setting changes were outside this test.

## Incomplete coverage

Three checks returned HTTP 403: agent skills, ticket summaries, and knowledge-base subfolders. Those responses do not distinguish all possible account entitlement and permission causes. Associated-ticket and prime-association reads returned HTTP 404 for the sampled ticket. A later CSAT-response list check also returned HTTP 404 for a sampled survey, so those operations remain unverified.

Additional checks succeeded for admin ticket-field detail, custom-object record lists/counts/filtering, and followed/participated discussion topics. Custom-object record responses exposed `display_id` rather than an internal `id`; record-detail lookup was not inferred from that response.

Some resources had no suitable existing records for a detail check. Import/export job creation, public replies, account configuration changes, and writes beyond the ticket-lifecycle test were not performed. Additional live write testing requires separately approved records and mutation scope.
