# Freshdesk command reference

Generated with `make docs`. Resource IDs appear as positional arguments in route order.
Use `freshdesk COMMAND --help` for individual flags.

## Global flags

```text
      --api-key string       API key override (prefer FRESHDESK_API_KEY or auth login)
      --domain string        Account slug or HTTPS hostname (or FRESHDESK_DOMAIN)
      --dry-run              Validate and print an API request plan without credentials or network
      --expect-plan string   Execute only if the rebuilt request has this SHA-256 plan hash
      --no-input             Never prompt; writes require --yes
  -o, --output string        Output format: table, json, yaml, csv (default "table")
      --profile string       Named profile (or FRESHDESK_PROFILE)
  -q, --quiet                Suppress informational stderr output
      --read-only            Block API writes, profile changes, and self-update
      --timeout duration     HTTP request timeout (default 30s)
  -v, --verbose              Log HTTP methods, paths, and status to stderr
  -y, --yes                  Confirm API writes and self-update without prompting
```

## Shared flags: tickets list

All corresponding API commands accept these flags. Fields that are not relevant to an endpoint should be supplied only when its reference documents them.

```text
      --all                    Collect pages, stopping at --max-pages or the API limit
      --company-id string      Query parameter company_id
      --email string           Query parameter email
      --filter string          Query parameter filter
      --include string         Query parameter include
      --max-pages int          Maximum pages to collect with --all, 1-500 (default 100)
      --next-token string      Query parameter next_token
      --order-by string        Query parameter order_by
      --order-type string      Query parameter order_type
      --page int               Page number (starts at 1) (default 1)
      --param stringArray      Query parameter KEY=VALUE; repeat for multiple values
      --per-page int           Page size, 1-100 (search APIs use 30) (default 30)
      --prev-token string      Query parameter prev_token
      --query string           Query parameter query
      --requester-id string    Query parameter requester_id
      --since string           Query parameter since
      --term string            Query parameter term
      --until string           Query parameter until
      --updated-since string   Query parameter updated_since
```

## Shared flags: tickets create

All corresponding API commands accept these flags. Fields that are not relevant to an endpoint should be supplied only when its reference documents them.

```text
      --attach stringArray     Attach a local file as attachments[]; repeat for more files
      --billable string        Body field billable
      --body string            Body field body
      --body-file string       Read the body field as UTF-8 text or HTML from FILE (or -)
      --company-id string      Body field company_id
      --content-html string    Body field content_html
      --description string     Body field description
      --email string           Body field email
      --filter string          Query parameter filter
      --group-id string        Body field group_id
      --if-match string        Send an If-Match precondition (when supported by the endpoint)
      --include string         Query parameter include
      --input-json string      JSON object, @FILE, or - for stdin
      --mobile string          Body field mobile
      --multipart              Send multipart form data, even without attachments
      --name string            Body field name
      --next-token string      Query parameter next_token
      --note string            Body field note
      --order-by string        Query parameter order_by
      --order-type string      Query parameter order_type
      --param stringArray      Query parameter KEY=VALUE; repeat for multiple values
      --phone string           Body field phone
      --prev-token string      Query parameter prev_token
      --priority string        Body field priority
      --private string         Body field private
      --query string           Query parameter query
      --requester-id string    Body field requester_id
      --responder-id string    Body field responder_id
      --set stringArray        Body FIELD=VALUE (JSON value or text); dotted fields build nested objects
      --since string           Query parameter since
      --source string          Body field source
      --status string          Body field status
      --subject string         Body field subject
      --term string            Query parameter term
      --time-spent string      Body field time_spent
      --title string           Body field title
      --type string            Body field type
      --until string           Query parameter until
      --updated-since string   Query parameter updated_since
      --upload stringArray     Upload FIELD=FILE for import CSVs, avatars, or other named file fields
      --user-id string         Body field user_id
```

## API commands

JSON input accepts the endpoint's native fields, including custom fields and nested structures.
`--param KEY=VALUE` supplies query parameters not listed as shortcuts.
Every write supports `--dry-run`, `--expect-plan HASH`, and `--yes`.

| Command | HTTP | Effect | Automatic pagination |
| --- | --- | --- | --- |
| `freshdesk tickets create [flags]` | `POST /api/v2/tickets` | write | false |
| `freshdesk tickets outbound-email [flags]` | `POST /api/v2/tickets/outbound_email` | write | false |
| `freshdesk tickets get ID [flags]` | `GET /api/v2/tickets/{id}` | read | false |
| `freshdesk tickets associated ID [flags]` | `GET /api/v2/tickets/{id}/associated_tickets` | read | true |
| `freshdesk tickets prime-association ID [flags]` | `GET /api/v2/tickets/{id}/prime_association` | read | false |
| `freshdesk tickets list [flags]` | `GET /api/v2/tickets` | read | true |
| `freshdesk tickets search [flags]` | `GET /api/v2/search/tickets` | read | true |
| `freshdesk tickets update ID [flags]` | `PUT /api/v2/tickets/{id}` | write | false |
| `freshdesk tickets bulk-update [flags]` | `POST /api/v2/tickets/bulk_update` | write | false |
| `freshdesk tickets forward ID [flags]` | `POST /api/v2/tickets/{id}/forward` | write | false |
| `freshdesk tickets merge [flags]` | `PUT /api/v2/tickets/merge` | write | false |
| `freshdesk tickets watchers list ID [flags]` | `GET /api/v2/tickets/{id}/watchers` | read | true |
| `freshdesk tickets watch ID [flags]` | `POST /api/v2/tickets/{id}/watch` | write | false |
| `freshdesk tickets unwatch ID [flags]` | `PUT /api/v2/tickets/{id}/unwatch` | write | false |
| `freshdesk tickets bulk-watch [flags]` | `PUT /api/v2/tickets/bulk_watch` | write | false |
| `freshdesk tickets bulk-unwatch [flags]` | `PUT /api/v2/tickets/bulk_unwatch` | write | false |
| `freshdesk tickets delete ID [flags]` | `DELETE /api/v2/tickets/{id}` | write | false |
| `freshdesk tickets bulk-delete [flags]` | `POST /api/v2/tickets/bulk_delete` | write | false |
| `freshdesk attachments delete ID [flags]` | `DELETE /api/v2/attachments/{id}` | write | false |
| `freshdesk tickets restore ID [flags]` | `PUT /api/v2/tickets/{id}/restore` | write | false |
| `freshdesk tickets fields [flags]` | `GET /api/v2/ticket_fields` | read | true |
| `freshdesk tickets conversations ID [flags]` | `GET /api/v2/tickets/{id}/conversations` | read | true |
| `freshdesk tickets time-entries ID [flags]` | `GET /api/v2/tickets/{id}/time_entries` | read | true |
| `freshdesk tickets ratings TICKET-ID [flags]` | `GET /api/v2/tickets/{ticket_id}/satisfaction_ratings` | read | true |
| `freshdesk tickets summary get ID [flags]` | `GET /api/v2/tickets/{id}/summary` | read | false |
| `freshdesk tickets summary update ID [flags]` | `PUT /api/v2/tickets/{id}/summary` | write | false |
| `freshdesk tickets summary delete ID [flags]` | `DELETE /api/v2/tickets/{id}/summary` | write | false |
| `freshdesk tickets archived get ID [flags]` | `GET /api/v2/tickets/archived/{id}` | read | false |
| `freshdesk tickets archived delete ID [flags]` | `DELETE /api/v2/tickets/archived/{id}` | write | false |
| `freshdesk tickets archived conversations ID [flags]` | `GET /api/v2/tickets/archived/{id}/conversations` | read | true |
| `freshdesk tickets access create ID [flags]` | `POST /api/v2/tickets/{id}/accesses` | write | false |
| `freshdesk tickets access update ID [flags]` | `PATCH /api/v2/tickets/{id}/accesses` | write | false |
| `freshdesk tickets access list ID [flags]` | `GET /api/v2/tickets/{id}/accesses` | read | true |
| `freshdesk tickets access delete ID [flags]` | `DELETE /api/v2/tickets/{id}/accesses` | write | false |
| `freshdesk ticket-fields list [flags]` | `GET /api/v2/admin/ticket_fields` | read | true |
| `freshdesk ticket-fields create [flags]` | `POST /api/v2/admin/ticket_fields` | write | false |
| `freshdesk ticket-fields get ID [flags]` | `GET /api/v2/admin/ticket_fields/{id}` | read | false |
| `freshdesk ticket-fields update ID [flags]` | `PUT /api/v2/admin/ticket_fields/{id}` | write | false |
| `freshdesk ticket-fields delete ID [flags]` | `DELETE /api/v2/admin/ticket_fields/{id}` | write | false |
| `freshdesk ticket-fields sections list ID [flags]` | `GET /api/v2/admin/ticket_fields/{id}/sections` | read | true |
| `freshdesk ticket-fields sections create ID [flags]` | `POST /api/v2/admin/ticket_fields/{id}/sections` | write | false |
| `freshdesk ticket-fields sections get ID SECTION-ID [flags]` | `GET /api/v2/admin/ticket_fields/{id}/sections/{section_id}` | read | false |
| `freshdesk ticket-fields sections update ID SECTION-ID [flags]` | `PUT /api/v2/admin/ticket_fields/{id}/sections/{section_id}` | write | false |
| `freshdesk ticket-fields sections delete ID SECTION-ID [flags]` | `DELETE /api/v2/admin/ticket_fields/{id}/sections/{section_id}` | write | false |
| `freshdesk ticket-forms list [flags]` | `GET /api/v2/ticket-forms` | read | true |
| `freshdesk ticket-forms create [flags]` | `POST /api/v2/ticket-forms` | write | false |
| `freshdesk ticket-forms get ID [flags]` | `GET /api/v2/ticket-forms/{id}` | read | false |
| `freshdesk ticket-forms update ID [flags]` | `PUT /api/v2/ticket-forms/{id}` | write | false |
| `freshdesk ticket-forms delete ID [flags]` | `DELETE /api/v2/ticket-forms/{id}` | write | false |
| `freshdesk ticket-forms clone ID [flags]` | `GET /api/v2/ticket-forms/{id}/clone` | write | false |
| `freshdesk ticket-forms fields get FORM-ID FIELD-ID [flags]` | `GET /api/v2/ticket-forms/{form_id}/fields/{field_id}` | read | false |
| `freshdesk ticket-forms fields update FORM-ID FIELD-ID [flags]` | `POST /api/v2/ticket-forms/{form_id}/fields/{field_id}` | write | false |
| `freshdesk ticket-forms fields delete FORM-ID FIELD-ID [flags]` | `DELETE /api/v2/ticket-forms/{form_id}/fields/{field_id}` | write | false |
| `freshdesk tickets reply ID [flags]` | `POST /api/v2/tickets/{id}/reply` | write | false |
| `freshdesk tickets note TICKET-ID [flags]` | `POST /api/v2/tickets/{ticket_id}/notes` | write | false |
| `freshdesk conversations update ID [flags]` | `PUT /api/v2/conversations/{id}` | write | false |
| `freshdesk conversations delete ID [flags]` | `DELETE /api/v2/conversations/{id}` | write | false |
| `freshdesk tickets reply-to-forward ID [flags]` | `POST /api/v2/tickets/{id}/reply_to_forward` | write | false |
| `freshdesk outbound-messages send [flags]` | `POST /api/v2/channels/outbound-messages` | write | false |
| `freshdesk outbound-messages get MESSAGE-ID [flags]` | `GET /api/v2/channels/outbound-messages/{message_id}` | read | false |
| `freshdesk account get [flags]` | `GET /api/v2/account` | read | false |
| `freshdesk account export [flags]` | `POST /api/v2/account/export` | write | false |
| `freshdesk jobs get ID [flags]` | `GET /api/v2/jobs/{id}` | read | false |
| `freshdesk contacts create [flags]` | `POST /api/v2/contacts` | write | false |
| `freshdesk contacts get ID [flags]` | `GET /api/v2/contacts/{id}` | read | false |
| `freshdesk contacts list [flags]` | `GET /api/v2/contacts` | read | true |
| `freshdesk contacts autocomplete [flags]` | `GET /api/v2/contacts/autocomplete` | read | false |
| `freshdesk contacts search [flags]` | `GET /api/v2/search/contacts` | read | true |
| `freshdesk contacts update ID [flags]` | `PUT /api/v2/contacts/{id}` | write | false |
| `freshdesk contacts delete ID [flags]` | `DELETE /api/v2/contacts/{id}` | write | false |
| `freshdesk contacts hard-delete ID [flags]` | `DELETE /api/v2/contacts/{id}/hard_delete` | write | false |
| `freshdesk contacts make-agent ID [flags]` | `PUT /api/v2/contacts/{id}/make_agent` | write | false |
| `freshdesk contacts restore ID [flags]` | `PUT /api/v2/contacts/{id}/restore` | write | false |
| `freshdesk contacts send-invite ID [flags]` | `PUT /api/v2/contacts/{id}/send_invite` | write | false |
| `freshdesk contacts merge [flags]` | `POST /api/v2/contacts/merge` | write | false |
| `freshdesk contacts exports create [flags]` | `POST /api/v2/contacts/export` | write | false |
| `freshdesk contacts exports get ID [flags]` | `GET /api/v2/contacts/export/{id}` | read | false |
| `freshdesk contacts imports list [flags]` | `GET /api/v2/contacts/imports` | read | true |
| `freshdesk contacts imports get ID [flags]` | `GET /api/v2/contacts/imports/{id}` | read | false |
| `freshdesk contacts imports create [flags]` | `POST /api/v2/contacts/imports` | write | false |
| `freshdesk contacts imports cancel ID [flags]` | `PUT /api/v2/contacts/imports/{id}/cancel` | write | false |
| `freshdesk contact-fields list [flags]` | `GET /api/v2/contact_fields` | read | true |
| `freshdesk contact-fields get ID [flags]` | `GET /api/v2/contact_fields/{id}` | read | false |
| `freshdesk contact-fields create [flags]` | `POST /api/v2/contact_fields` | write | false |
| `freshdesk contact-fields update ID [flags]` | `PUT /api/v2/contact_fields/{id}` | write | false |
| `freshdesk contact-fields delete ID [flags]` | `DELETE /api/v2/contact_fields/{id}` | write | false |
| `freshdesk agents get ID [flags]` | `GET /api/v2/agents/{id}` | read | false |
| `freshdesk agents list [flags]` | `GET /api/v2/agents` | read | true |
| `freshdesk agents create [flags]` | `POST /api/v2/agents/` | write | false |
| `freshdesk agents update ID [flags]` | `PUT /api/v2/agents/{id}` | write | false |
| `freshdesk agents delete ID [flags]` | `DELETE /api/v2/agents/{id}` | write | false |
| `freshdesk agents me [flags]` | `GET /api/v2/agents/me` | read | false |
| `freshdesk agents bulk-create [flags]` | `POST /api/v2/agents/bulk` | write | false |
| `freshdesk agents autocomplete [flags]` | `GET /api/v2/agents/autocomplete` | read | false |
| `freshdesk agents availability list [flags]` | `GET /api/v2/agents` | read | true |
| `freshdesk agents availability get ID [flags]` | `GET /api/v2/agents/{id}/availability` | read | false |
| `freshdesk agents availability update ID [flags]` | `PUT /api/v2/agents/{id}/availability` | write | false |
| `freshdesk skills get ID [flags]` | `GET /api/v2/admin/skills/{id}` | read | false |
| `freshdesk skills list [flags]` | `GET /api/v2/admin/skills/` | read | true |
| `freshdesk skills create [flags]` | `POST /api/v2/admin/skills` | write | false |
| `freshdesk skills delete ID [flags]` | `DELETE /api/v2/admin/skills/{id}` | write | false |
| `freshdesk skills update ID [flags]` | `PUT /api/v2/admin/skills/{id}` | write | false |
| `freshdesk roles get ID [flags]` | `GET /api/v2/roles/{id}` | read | false |
| `freshdesk roles list [flags]` | `GET /api/v2/roles` | read | true |
| `freshdesk groups create [flags]` | `POST /api/v2/groups` | write | false |
| `freshdesk groups get ID [flags]` | `GET /api/v2/groups/{id}` | read | false |
| `freshdesk groups list [flags]` | `GET /api/v2/groups` | read | true |
| `freshdesk groups update ID [flags]` | `PUT /api/v2/groups/{id}` | write | false |
| `freshdesk groups delete ID [flags]` | `DELETE /api/v2/groups/{id}` | write | false |
| `freshdesk admin-groups create [flags]` | `POST /api/v2/admin/groups` | write | false |
| `freshdesk admin-groups get ID [flags]` | `GET /api/v2/admin/groups/{id}` | read | false |
| `freshdesk admin-groups list [flags]` | `GET /api/v2/admin/groups` | read | true |
| `freshdesk admin-groups update ID [flags]` | `PUT /api/v2/admin/groups/{id}` | write | false |
| `freshdesk admin-groups delete ID [flags]` | `DELETE /api/v2/admin/groups/{id}` | write | false |
| `freshdesk admin-groups agents list ID [flags]` | `GET /api/v2/admin/groups/{id}/agents` | read | true |
| `freshdesk admin-groups agents update ID [flags]` | `PATCH /api/v2/admin/groups/{id}/agents` | write | false |
| `freshdesk companies create [flags]` | `POST /api/v2/companies` | write | false |
| `freshdesk companies get ID [flags]` | `GET /api/v2/companies/{id}` | read | false |
| `freshdesk companies list [flags]` | `GET /api/v2/companies` | read | true |
| `freshdesk companies autocomplete [flags]` | `GET /api/v2/companies/autocomplete` | read | false |
| `freshdesk companies search [flags]` | `GET /api/v2/search/companies` | read | true |
| `freshdesk companies update ID [flags]` | `PUT /api/v2/companies/{id}` | write | false |
| `freshdesk companies delete ID [flags]` | `DELETE /api/v2/companies/{id}` | write | false |
| `freshdesk companies exports create [flags]` | `POST /api/v2/companies/export` | write | false |
| `freshdesk companies exports get ID [flags]` | `GET /api/v2/companies/export/{id}` | read | false |
| `freshdesk companies imports list [flags]` | `GET /api/v2/companies/imports` | read | true |
| `freshdesk companies imports get ID [flags]` | `GET /api/v2/companies/imports/{id}` | read | false |
| `freshdesk companies imports create [flags]` | `POST /api/v2/companies/imports` | write | false |
| `freshdesk companies imports cancel ID [flags]` | `PUT /api/v2/companies/imports/{id}/cancel` | write | false |
| `freshdesk company-fields list [flags]` | `GET /api/v2/company_fields` | read | true |
| `freshdesk company-fields get ID [flags]` | `GET /api/v2/company_fields/{id}` | read | false |
| `freshdesk company-fields create [flags]` | `POST /api/v2/company_fields` | write | false |
| `freshdesk company-fields update ID [flags]` | `PUT /api/v2/company_fields/{id}` | write | false |
| `freshdesk company-fields delete ID [flags]` | `DELETE /api/v2/company_fields/{id}` | write | false |
| `freshdesk custom-objects list [flags]` | `GET /api/v2/custom_objects/schemas` | read | false |
| `freshdesk custom-objects get SCHEMA-ID [flags]` | `GET /api/v2/custom_objects/schemas/{schema_id}` | read | false |
| `freshdesk custom-objects records create SCHEMA-ID [flags]` | `POST /api/v2/custom_objects/schemas/{schema_id}/records/` | write | false |
| `freshdesk custom-objects records get SCHEMA-ID RECORD-ID [flags]` | `GET /api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | read | false |
| `freshdesk custom-objects records update SCHEMA-ID RECORD-ID [flags]` | `PUT /api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | write | false |
| `freshdesk custom-objects records delete SCHEMA-ID RECORD-ID [flags]` | `DELETE /api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | write | false |
| `freshdesk custom-objects records list SCHEMA-ID [flags]` | `GET /api/v2/custom_objects/schemas/{schema_id}/records` | read | false |
| `freshdesk custom-objects records count SCHEMA-ID [flags]` | `GET /api/v2/custom_objects/schemas/{schema_id}/records/count` | read | false |
| `freshdesk custom-objects records filter SCHEMA-ID [flags]` | `GET /api/v2/custom_objects/schemas/{schema_id}/records` | read | false |
| `freshdesk canned-responses create [flags]` | `POST /api/v2/canned_responses` | write | false |
| `freshdesk canned-responses get ID [flags]` | `GET /api/v2/canned_responses/{id}` | read | false |
| `freshdesk canned-responses update ID [flags]` | `PUT /api/v2/canned_responses/{id}` | write | false |
| `freshdesk canned-responses bulk-create [flags]` | `POST /api/v2/canned_responses/bulk` | write | false |
| `freshdesk canned-responses folders create [flags]` | `POST /api/v2/canned_response_folders` | write | false |
| `freshdesk canned-responses folders update ID [flags]` | `PUT /api/v2/canned_response_folders/{id}` | write | false |
| `freshdesk canned-responses folders list [flags]` | `GET /api/v2/canned_response_folders` | read | true |
| `freshdesk canned-responses folders get ID [flags]` | `GET /api/v2/canned_response_folders/{id}` | read | false |
| `freshdesk canned-responses folders responses ID [flags]` | `GET /api/v2/canned_response_folders/{id}/responses` | read | true |
| `freshdesk discussions categories create [flags]` | `POST /api/v2/discussions/categories` | write | false |
| `freshdesk discussions categories get ID [flags]` | `GET /api/v2/discussions/categories/{id}` | read | false |
| `freshdesk discussions categories list [flags]` | `GET /api/v2/discussions/categories` | read | true |
| `freshdesk discussions forums list ID [flags]` | `GET /api/v2/discussions/categories/{id}/forums` | read | true |
| `freshdesk discussions categories update ID [flags]` | `PUT /api/v2/discussions/categories/{id}` | write | false |
| `freshdesk discussions categories delete ID [flags]` | `DELETE /api/v2/discussions/categories/{id}` | write | false |
| `freshdesk discussions forums create CATEGORY-ID [flags]` | `POST /api/v2/discussions/categories/{category_id}/forums` | write | false |
| `freshdesk discussions forums get ID [flags]` | `GET /api/v2/discussions/forums/{id}` | read | false |
| `freshdesk discussions topics list ID [flags]` | `GET /api/v2/discussions/forums/{id}/topics` | read | true |
| `freshdesk discussions forums update ID [flags]` | `PUT /api/v2/discussions/forums/{id}` | write | false |
| `freshdesk discussions forums delete ID [flags]` | `DELETE /api/v2/discussions/forums/{id}` | write | false |
| `freshdesk discussions forums follow ID [flags]` | `POST /api/v2/discussions/forums/{id}/follow` | write | false |
| `freshdesk discussions forums unfollow FORUM-ID [flags]` | `DELETE /api/v2/discussions/forums/{forum_id}/follow` | write | false |
| `freshdesk discussions forums follow-status ID [flags]` | `GET /api/v2/discussions/forums/{id}/follow` | read | false |
| `freshdesk discussions topics create FORUM-ID [flags]` | `POST /api/v2/discussions/forums/{forum_id}/topics` | write | false |
| `freshdesk discussions topics get ID [flags]` | `GET /api/v2/discussions/topics/{id}` | read | false |
| `freshdesk discussions comments list ID [flags]` | `GET /api/v2/discussions/topics/{id}/comments` | read | true |
| `freshdesk discussions topics update ID [flags]` | `PUT /api/v2/discussions/topics/{id}` | write | false |
| `freshdesk discussions topics delete ID [flags]` | `DELETE /api/v2/discussions/topics/{id}` | write | false |
| `freshdesk discussions topics follow ID [flags]` | `POST /api/v2/discussions/topics/{id}/follow` | write | false |
| `freshdesk discussions topics unfollow TOPIC-ID [flags]` | `DELETE /api/v2/discussions/topics/{topic_id}/follow` | write | false |
| `freshdesk discussions topics follow-status ID [flags]` | `GET /api/v2/discussions/topics/{id}/follow` | read | false |
| `freshdesk discussions topics followed [flags]` | `GET /api/v2/discussions/topics/followed_by` | read | true |
| `freshdesk discussions topics participated [flags]` | `GET /api/v2/discussions/topics/participated_by` | read | true |
| `freshdesk discussions comments create TOPIC-ID [flags]` | `POST /api/v2/discussions/topics/{topic_id}/comments` | write | false |
| `freshdesk discussions comments update ID [flags]` | `PUT /api/v2/discussions/comments/{id}` | write | false |
| `freshdesk discussions comments delete ID [flags]` | `DELETE /api/v2/discussions/comments/{id}` | write | false |
| `freshdesk solutions categories create [flags]` | `POST /api/v2/solutions/categories` | write | false |
| `freshdesk solutions categories translations create ID LANGUAGE-CODE [flags]` | `POST /api/v2/solutions/categories/{id}/{language_code}` | write | false |
| `freshdesk solutions categories update ID [flags]` | `PUT /api/v2/solutions/categories/{id}` | write | false |
| `freshdesk solutions categories translations update ID LANGUAGE-CODE [flags]` | `PUT /api/v2/solutions/categories/{id}/{language_code}` | write | false |
| `freshdesk solutions categories get ID [flags]` | `GET /api/v2/solutions/categories/{id}` | read | false |
| `freshdesk solutions categories translations get ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/categories/{id}/{language_code}` | read | false |
| `freshdesk solutions categories list [flags]` | `GET /api/v2/solutions/categories` | read | true |
| `freshdesk solutions categories translations list LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/categories/{language_code}` | read | true |
| `freshdesk solutions categories delete ID [flags]` | `DELETE /api/v2/solutions/categories/{id}` | write | false |
| `freshdesk solutions folders create ID [flags]` | `POST /api/v2/solutions/categories/{id}/folders` | write | false |
| `freshdesk solutions folders translations create ID LANGUAGE-CODE [flags]` | `POST /api/v2/solutions/folders/{id}/{language_code}` | write | false |
| `freshdesk solutions folders update ID [flags]` | `PUT /api/v2/solutions/folders/{id}` | write | false |
| `freshdesk solutions folders translations update ID LANGUAGE-CODE [flags]` | `PUT /api/v2/solutions/folders/{id}/{language_code}` | write | false |
| `freshdesk solutions folders get ID [flags]` | `GET /api/v2/solutions/folders/{id}` | read | false |
| `freshdesk solutions folders translations get ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/folders/{id}/{language_code}` | read | false |
| `freshdesk solutions folders list CATEGORY-ID [flags]` | `GET /api/v2/solutions/categories/{category_id}/folders` | read | true |
| `freshdesk solutions folders translations list CATEGORY-ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/categories/{category_id}/folders/{language_code}` | read | true |
| `freshdesk solutions folders subfolders FOLDER-ID [flags]` | `GET /api/v2/solutions/folders/{folder_id}/subfolders` | read | true |
| `freshdesk solutions folders translations subfolders FOLDER-ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/folders/{folder_id}/subfolders/{language_code}` | read | true |
| `freshdesk solutions folders delete ID [flags]` | `DELETE /api/v2/solutions/folders/{id}` | write | false |
| `freshdesk solutions articles create ID [flags]` | `POST /api/v2/solutions/folders/{id}/articles` | write | false |
| `freshdesk solutions articles translations create ID LANGUAGE-CODE [flags]` | `POST /api/v2/solutions/articles/{id}/{language_code}` | write | false |
| `freshdesk solutions articles update ID [flags]` | `PUT /api/v2/solutions/articles/{id}` | write | false |
| `freshdesk solutions articles translations update ID LANGUAGE-CODE [flags]` | `PUT /api/v2/solutions/articles/{id}/{language_code}` | write | false |
| `freshdesk solutions articles get ID [flags]` | `GET /api/v2/solutions/articles/{id}` | read | false |
| `freshdesk solutions articles translations get ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/articles/{id}/{language_code}` | read | false |
| `freshdesk solutions articles list ID [flags]` | `GET /api/v2/solutions/folders/{id}/articles` | read | true |
| `freshdesk solutions articles translations list ID LANGUAGE-CODE [flags]` | `GET /api/v2/solutions/folders/{id}/articles/{language_code}` | read | true |
| `freshdesk solutions articles delete ID [flags]` | `DELETE /api/v2/solutions/articles/{id}` | write | false |
| `freshdesk solutions articles search [flags]` | `GET /api/v2/search/solutions` | read | false |
| `freshdesk csat surveys list [flags]` | `GET /api/v2/customer-satisfaction/surveys` | read | true |
| `freshdesk csat responses create SURVEY-ID [flags]` | `POST /api/v2/customer-satisfaction/surveys/{survey_id}/responses` | write | false |
| `freshdesk csat responses list SURVEY-ID [flags]` | `GET /api/v2/customer-satisfaction/surveys/{survey_id}/responses` | read | true |
| `freshdesk surveys list [flags]` | `GET /api/v2/surveys` | read | true |
| `freshdesk tickets ratings create TICKET-ID [flags]` | `POST /api/v2/tickets/{ticket_id}/satisfaction_ratings` | write | false |
| `freshdesk surveys ratings list [flags]` | `GET /api/v2/surveys/satisfaction_ratings` | read | true |
| `freshdesk service tasks create [flags]` | `POST /api/v2/tickets` | write | false |
| `freshdesk service tasks get ID [flags]` | `GET /api/v2/tickets/{id}` | read | false |
| `freshdesk service tasks update ID [flags]` | `PUT /api/v2/tickets/{id}` | write | false |
| `freshdesk service tasks delete ID [flags]` | `DELETE /api/v2/tickets/{id}` | write | false |
| `freshdesk service groups create [flags]` | `POST /api/v2/groups` | write | false |
| `freshdesk service groups get ID [flags]` | `GET /api/v2/groups/{id}` | read | false |
| `freshdesk service groups update ID [flags]` | `PUT /api/v2/groups/{id}` | write | false |
| `freshdesk service groups delete ID [flags]` | `DELETE /api/v2/groups/{id}` | write | false |
| `freshdesk service make-technician ID [flags]` | `PUT /api/v2/contacts/{id}/make_agent` | write | false |
| `freshdesk time-entries create TICKET-ID [flags]` | `POST /api/v2/tickets/{ticket_id}/time_entries` | write | false |
| `freshdesk time-entries list [flags]` | `GET /api/v2/time_entries` | read | true |
| `freshdesk time-entries update ID [flags]` | `PUT /api/v2/time_entries/{id}` | write | false |
| `freshdesk time-entries toggle TIME-ENTRY-ID [flags]` | `PUT /api/v2/time_entries/{time_entry_id}/toggle_timer` | write | false |
| `freshdesk time-entries delete ID [flags]` | `DELETE /api/v2/time_entries/{id}` | write | false |
| `freshdesk email-configs get ID [flags]` | `GET /api/v2/email_configs/{id}` | read | false |
| `freshdesk email-configs list [flags]` | `GET /api/v2/email_configs` | read | true |
| `freshdesk email mailboxes create [flags]` | `POST /api/v2/email/mailboxes` | write | false |
| `freshdesk email mailboxes get ID [flags]` | `GET /api/v2/email/mailboxes/{id}` | read | false |
| `freshdesk email mailboxes list [flags]` | `GET /api/v2/email/mailboxes` | read | true |
| `freshdesk email mailboxes update ID [flags]` | `PUT /api/v2/email/mailboxes/{id}` | write | false |
| `freshdesk email mailboxes delete ID [flags]` | `DELETE /api/v2/email/mailboxes/{id}` | write | false |
| `freshdesk email settings update [flags]` | `PUT /api/v2/email/settings` | write | false |
| `freshdesk email settings get [flags]` | `GET /api/v2/email/settings` | read | false |
| `freshdesk email bcc update [flags]` | `PUT /api/v2/notifications/email/bcc` | write | false |
| `freshdesk email bcc get [flags]` | `GET /api/v2/notifications/email/bcc` | read | false |
| `freshdesk products get ID [flags]` | `GET /api/v2/products/{id}` | read | false |
| `freshdesk products list [flags]` | `GET /api/v2/products` | read | true |
| `freshdesk business-hours get ID [flags]` | `GET /api/v2/business_hours/{id}` | read | false |
| `freshdesk business-hours list [flags]` | `GET /api/v2/business_hours` | read | true |
| `freshdesk scenario-automations list [flags]` | `GET /api/v2/scenario_automations` | read | true |
| `freshdesk sla-policies create [flags]` | `POST /api/v2/sla_policies` | write | false |
| `freshdesk sla-policies list [flags]` | `GET /api/v2/sla_policies` | read | true |
| `freshdesk sla-policies update ID [flags]` | `PUT /api/v2/sla_policies/{id}` | write | false |
| `freshdesk activities publish [flags]` | `POST /api/v2/contact-activities` | write | false |
| `freshdesk automations list AUTOMATION-TYPE-ID [flags]` | `GET /api/v2/automations/{automation_type_id}/rules` | read | true |
| `freshdesk automations get AUTOMATION-TYPE-ID ID [flags]` | `GET /api/v2/automations/{automation_type_id}/rules/{id}` | read | false |
| `freshdesk automations delete AUTOMATION-TYPE-ID ID [flags]` | `DELETE /api/v2/automations/{automation_type_id}/rules/{id}` | write | false |
| `freshdesk automations create AUTOMATION-TYPE-ID [flags]` | `POST /api/v2/automations/{automation_type_id}/rules` | write | false |
| `freshdesk automations update AUTOMATION-TYPE-ID ID [flags]` | `PUT /api/v2/automations/{automation_type_id}/rules/{id}` | write | false |
| `freshdesk settings get [flags]` | `GET /api/v2/settings/helpdesk` | read | false |
| `freshdesk threads create [flags]` | `POST /api/v2/collaboration/threads` | write | false |
| `freshdesk threads get ID [flags]` | `GET /api/v2/collaboration/threads/{id}` | read | false |
| `freshdesk threads update ID [flags]` | `PUT /api/v2/collaboration/threads/{id}` | write | false |
| `freshdesk threads delete ID [flags]` | `DELETE /api/v2/collaboration/threads/{id}` | write | false |
| `freshdesk threads messages create [flags]` | `POST /api/v2/collaboration/messages` | write | false |
| `freshdesk threads messages get ID [flags]` | `GET /api/v2/collaboration/messages/{id}` | read | false |
| `freshdesk threads messages update ID [flags]` | `PUT /api/v2/collaboration/messages/{id}` | write | false |
| `freshdesk threads messages delete ID [flags]` | `DELETE /api/v2/collaboration/messages/{id}` | write | false |
| `freshdesk threads messages quote ID [flags]` | `GET /api/v2/collaboration/messages/{id}/generate-quote` | read | false |
| `freshdesk audit-log exports create [flags]` | `POST /api/v2/audit_log/export` | write | false |
| `freshdesk audit-log exports get ID [flags]` | `GET /api/v2/audit_log/export/{id}` | read | false |

## Local and discovery commands

### freshdesk api methods

List the reviewed API operations offline

```text
freshdesk api methods
```

### freshdesk api request

Call an API route with the same preview and write controls

```text
freshdesk api request METHOD PATH [flags]
```

Known routes use the registry's semantic effect, including GET /ticket-forms/ID/clone.
Unknown routes are treated as writes, including GET. Paths stay on the configured helpdesk origin.
Pagination is manual for raw calls; use --param page=N. Named commands support --all where documented.

```text
      --attach stringArray     Attach a local file as attachments[]; repeat for more files
      --billable string        Body field billable
      --body string            Body field body
      --body-file string       Read the body field as UTF-8 text or HTML from FILE (or -)
      --company-id string      Body field company_id
      --content-html string    Body field content_html
      --description string     Body field description
      --email string           Body field email
      --filter string          Query parameter filter
      --group-id string        Body field group_id
      --if-match string        Send an If-Match precondition (when supported by the endpoint)
      --include string         Query parameter include
      --input-json string      JSON object, @FILE, or - for stdin
      --mobile string          Body field mobile
      --multipart              Send multipart form data, even without attachments
      --name string            Body field name
      --next-token string      Query parameter next_token
      --note string            Body field note
      --order-by string        Query parameter order_by
      --order-type string      Query parameter order_type
      --param stringArray      Query parameter KEY=VALUE; repeat for multiple values
      --phone string           Body field phone
      --prev-token string      Query parameter prev_token
      --priority string        Body field priority
      --private string         Body field private
      --query string           Query parameter query
      --requester-id string    Body field requester_id
      --responder-id string    Body field responder_id
      --set stringArray        Body FIELD=VALUE (JSON value or text); dotted fields build nested objects
      --since string           Query parameter since
      --source string          Body field source
      --status string          Body field status
      --subject string         Body field subject
      --term string            Query parameter term
      --time-spent string      Body field time_spent
      --title string           Body field title
      --type string            Body field type
      --until string           Query parameter until
      --updated-since string   Query parameter updated_since
      --upload stringArray     Upload FIELD=FILE for import CSVs, avatars, or other named file fields
      --user-id string         Body field user_id
```

### freshdesk auth list

List saved profiles without credentials

```text
freshdesk auth list
```

### freshdesk auth login

Verify and save an API key in the OS keychain

```text
freshdesk auth login [flags]
```

Get your API key from Freshdesk > profile picture > Profile settings.
Login checks GET /api/v2/agents/me before saving. Use --no-verify to save without a request.
The default store is the OS keychain. --insecure-storage opts into a plaintext 0600 file.

```text
      --insecure-storage   Use an owner-only plaintext secrets file instead of the keychain
      --key-stdin          Read the API key from stdin
      --no-verify          Save without calling Freshdesk
      --save-read-only     Save a read-only restriction in this profile
```

### freshdesk auth logout

Remove the selected profile and saved key from this machine

```text
freshdesk auth logout
```

Removes local credentials only. To revoke the key remotely, regenerate it in Freshdesk Profile settings.

### freshdesk auth status

Inspect credential source without revealing or loading the key

```text
freshdesk auth status [flags]
```

```text
      --verify   Verify access with one GET /agents/me request
```

### freshdesk auth use

Select the default saved profile

```text
freshdesk auth use NAME
```

### freshdesk commands

Discover commands, flags, and effects offline

```text
freshdesk commands
```

### freshdesk commands describe

Describe a stable command ID offline

```text
freshdesk commands describe ID
```

### freshdesk commands resolve

Resolve arguments with the real parser, without executing

```text
freshdesk commands resolve -- ARGUMENTS
```

### freshdesk completion

Generate shell completion

```text
freshdesk completion [bash|zsh|fish|powershell]
```

### freshdesk config list-profiles

List saved profiles without credentials

```text
freshdesk config list-profiles
```

### freshdesk config show

Inspect credential source without revealing or loading the key

```text
freshdesk config show [flags]
```

```text
      --verify   Verify access with one GET /agents/me request
```

### freshdesk config use-profile

Select the default saved profile

```text
freshdesk config use-profile NAME
```

### freshdesk doctor

Check local configuration, optionally verify API access

```text
freshdesk doctor [flags]
```

```text
      --online   Verify credentials with GET /agents/me
```

### freshdesk update

Check or install a checksummed GitHub release

```text
freshdesk update [flags]
```

Downloads the release for this platform and verifies its SHA-256 checksum before replacing the binary.
Source installs in a Go bin directory print the source update command instead.
There is no background release check. --check never writes a cache or replaces files.

```text
      --check   Only check the latest release; no file writes
```

### freshdesk version

Print version, commit, and build date

```text
freshdesk version
```
