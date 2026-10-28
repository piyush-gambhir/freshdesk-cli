# Freshdesk API coverage

The reviewed registry covers **266 documented operation variants**, including translations and field-service variants that reuse an endpoint. This is a documentation coverage count, not a count of unique routes or live-verified features.

Source: [Freshdesk API reference](https://developers.freshdesk.com/api/), retrieved 2026-10-29.
The endpoint-facts snapshot is in `cli-go/internal/registry/testdata/freshdesk-api.json`.
Registry SHA-256: `0e0fc88b737264c39c7c9fedbd84c56d62d031c467c58196124c7447719cac70`.

All rows have named commands and pass request-routing fixture tests. No live Freshdesk tenant was used. Plan entitlements and agent permissions still apply.

| Command | Verb | Path | Effect | Reference |
| --- | --- | --- | --- | --- |
| `tickets create` | POST | `/api/v2/tickets` | write | [reference](https://developers.freshdesk.com/api/#create_ticket) |
| `tickets outbound-email` | POST | `/api/v2/tickets/outbound_email` | write | [reference](https://developers.freshdesk.com/api/#create_outbound_email) |
| `tickets get` | GET | `/api/v2/tickets/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_a_ticket) |
| `tickets associated` | GET | `/api/v2/tickets/{id}/associated_tickets` | read | [reference](https://developers.freshdesk.com/api/#get_associated_tickets) |
| `tickets prime-association` | GET | `/api/v2/tickets/{id}/prime_association` | read | [reference](https://developers.freshdesk.com/api/#ticket_prime_association) |
| `tickets list` | GET | `/api/v2/tickets` | read | [reference](https://developers.freshdesk.com/api/#list_all_tickets) |
| `tickets search` | GET | `/api/v2/search/tickets` | read | [reference](https://developers.freshdesk.com/api/#filter_tickets) |
| `tickets update` | PUT | `/api/v2/tickets/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_ticket) |
| `tickets bulk-update` | POST | `/api/v2/tickets/bulk_update` | write | [reference](https://developers.freshdesk.com/api/#bulk_update_tickets) |
| `tickets forward` | POST | `/api/v2/tickets/{id}/forward` | write | [reference](https://developers.freshdesk.com/api/#forward_tickets) |
| `tickets merge` | PUT | `/api/v2/tickets/merge` | write | [reference](https://developers.freshdesk.com/api/#ticket_merge) |
| `tickets watchers list` | GET | `/api/v2/tickets/{id}/watchers` | read | [reference](https://developers.freshdesk.com/api/#list_all_watchers) |
| `tickets watch` | POST | `/api/v2/tickets/{id}/watch` | write | [reference](https://developers.freshdesk.com/api/#add_watcher) |
| `tickets unwatch` | PUT | `/api/v2/tickets/{id}/unwatch` | write | [reference](https://developers.freshdesk.com/api/#remove_watcher) |
| `tickets bulk-watch` | PUT | `/api/v2/tickets/bulk_watch` | write | [reference](https://developers.freshdesk.com/api/#add_watcher_multiple_tickets) |
| `tickets bulk-unwatch` | PUT | `/api/v2/tickets/bulk_unwatch` | write | [reference](https://developers.freshdesk.com/api/#remove_watcher_multiple_tickets) |
| `tickets delete` | DELETE | `/api/v2/tickets/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_a_ticket) |
| `tickets bulk-delete` | POST | `/api/v2/tickets/bulk_delete` | write | [reference](https://developers.freshdesk.com/api/#bulk_delete_tickets) |
| `attachments delete` | DELETE | `/api/v2/attachments/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_an_attachment) |
| `tickets restore` | PUT | `/api/v2/tickets/{id}/restore` | write | [reference](https://developers.freshdesk.com/api/#restore_a_ticket) |
| `tickets fields` | GET | `/api/v2/ticket_fields` | read | [reference](https://developers.freshdesk.com/api/#list_all_ticket_fields) |
| `tickets conversations` | GET | `/api/v2/tickets/{id}/conversations` | read | [reference](https://developers.freshdesk.com/api/#list_all_ticket_notes) |
| `tickets time-entries` | GET | `/api/v2/tickets/{id}/time_entries` | read | [reference](https://developers.freshdesk.com/api/#list_all_ticket_timeentries) |
| `tickets ratings` | GET | `/api/v2/tickets/{ticket_id}/satisfaction_ratings` | read | [reference](https://developers.freshdesk.com/api/#view_ticket_satisfaction_ratings) |
| `tickets summary get` | GET | `/api/v2/tickets/{id}/summary` | read | [reference](https://developers.freshdesk.com/api/#view_ticket_summary) |
| `tickets summary update` | PUT | `/api/v2/tickets/{id}/summary` | write | [reference](https://developers.freshdesk.com/api/#update_ticket_summary) |
| `tickets summary delete` | DELETE | `/api/v2/tickets/{id}/summary` | write | [reference](https://developers.freshdesk.com/api/#delete_ticket_summary) |
| `tickets archived get` | GET | `/api/v2/tickets/archived/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_archive_ticket) |
| `tickets archived delete` | DELETE | `/api/v2/tickets/archived/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_archive_ticket) |
| `tickets archived conversations` | GET | `/api/v2/tickets/archived/{id}/conversations` | read | [reference](https://developers.freshdesk.com/api/#list_all_conversation_of_archive_tickets) |
| `tickets access create` | POST | `/api/v2/tickets/{id}/accesses` | write | [reference](https://developers.freshdesk.com/api/#create_user_accesses) |
| `tickets access update` | PATCH | `/api/v2/tickets/{id}/accesses` | write | [reference](https://developers.freshdesk.com/api/#update_user_accesses) |
| `tickets access list` | GET | `/api/v2/tickets/{id}/accesses` | read | [reference](https://developers.freshdesk.com/api/#get_user_accesses) |
| `tickets access delete` | DELETE | `/api/v2/tickets/{id}/accesses` | write | [reference](https://developers.freshdesk.com/api/#delete_user_accesses) |
| `ticket-fields list` | GET | `/api/v2/admin/ticket_fields` | read | [reference](https://developers.freshdesk.com/api/#list_all_ticket_fields) |
| `ticket-fields create` | POST | `/api/v2/admin/ticket_fields` | write | [reference](https://developers.freshdesk.com/api/#create_ticket_field) |
| `ticket-fields get` | GET | `/api/v2/admin/ticket_fields/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_ticket_field) |
| `ticket-fields update` | PUT | `/api/v2/admin/ticket_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_ticket_field) |
| `ticket-fields delete` | DELETE | `/api/v2/admin/ticket_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_ticket_field) |
| `ticket-fields sections list` | GET | `/api/v2/admin/ticket_fields/{id}/sections` | read | [reference](https://developers.freshdesk.com/api/#list_all_section_fields) |
| `ticket-fields sections create` | POST | `/api/v2/admin/ticket_fields/{id}/sections` | write | [reference](https://developers.freshdesk.com/api/#create_section) |
| `ticket-fields sections get` | GET | `/api/v2/admin/ticket_fields/{id}/sections/{section_id}` | read | [reference](https://developers.freshdesk.com/api/#list_specific_section) |
| `ticket-fields sections update` | PUT | `/api/v2/admin/ticket_fields/{id}/sections/{section_id}` | write | [reference](https://developers.freshdesk.com/api/#update_section) |
| `ticket-fields sections delete` | DELETE | `/api/v2/admin/ticket_fields/{id}/sections/{section_id}` | write | [reference](https://developers.freshdesk.com/api/#delete_section) |
| `ticket-forms list` | GET | `/api/v2/ticket-forms` | read | [reference](https://developers.freshdesk.com/api/#list_all_ticket_forms) |
| `ticket-forms create` | POST | `/api/v2/ticket-forms` | write | [reference](https://developers.freshdesk.com/api/#create_ticket_form) |
| `ticket-forms get` | GET | `/api/v2/ticket-forms/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_ticket_form) |
| `ticket-forms update` | PUT | `/api/v2/ticket-forms/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_ticket_form) |
| `ticket-forms delete` | DELETE | `/api/v2/ticket-forms/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_ticket_form) |
| `ticket-forms clone` | GET | `/api/v2/ticket-forms/{id}/clone` | write | [reference](https://developers.freshdesk.com/api/#clone_ticket_form) |
| `ticket-forms fields get` | GET | `/api/v2/ticket-forms/{form_id}/fields/{field_id}` | read | [reference](https://developers.freshdesk.com/api/#view_ticket_form_field) |
| `ticket-forms fields update` | POST | `/api/v2/ticket-forms/{form_id}/fields/{field_id}` | write | [reference](https://developers.freshdesk.com/api/#update_ticket_form_field) |
| `ticket-forms fields delete` | DELETE | `/api/v2/ticket-forms/{form_id}/fields/{field_id}` | write | [reference](https://developers.freshdesk.com/api/#delete_ticket_form_field) |
| `tickets reply` | POST | `/api/v2/tickets/{id}/reply` | write | [reference](https://developers.freshdesk.com/api/#reply_ticket) |
| `tickets note` | POST | `/api/v2/tickets/{ticket_id}/notes` | write | [reference](https://developers.freshdesk.com/api/#add_note_to_a_ticket) |
| `conversations update` | PUT | `/api/v2/conversations/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_conversation) |
| `conversations delete` | DELETE | `/api/v2/conversations/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_conversation) |
| `tickets reply-to-forward` | POST | `/api/v2/tickets/{id}/reply_to_forward` | write | [reference](https://developers.freshdesk.com/api/#reply_to_forward) |
| `outbound-messages send` | POST | `/api/v2/channels/outbound-messages` | write | [reference](https://developers.freshdesk.com/api/#send_outbound_message) |
| `outbound-messages get` | GET | `/api/v2/channels/outbound-messages/{message_id}` | read | [reference](https://developers.freshdesk.com/api/#retrieve_outbound_message) |
| `account get` | GET | `/api/v2/account` | read | [reference](https://developers.freshdesk.com/api/#view_account) |
| `account export` | POST | `/api/v2/account/export` | write | [reference](https://developers.freshdesk.com/api/#generate_export) |
| `jobs get` | GET | `/api/v2/jobs/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_job) |
| `contacts create` | POST | `/api/v2/contacts` | write | [reference](https://developers.freshdesk.com/api/#create_contact) |
| `contacts get` | GET | `/api/v2/contacts/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_contact) |
| `contacts list` | GET | `/api/v2/contacts` | read | [reference](https://developers.freshdesk.com/api/#list_all_contacts) |
| `contacts autocomplete` | GET | `/api/v2/contacts/autocomplete` | read | [reference](https://developers.freshdesk.com/api/#search_contacts) |
| `contacts search` | GET | `/api/v2/search/contacts` | read | [reference](https://developers.freshdesk.com/api/#filter_contacts) |
| `contacts update` | PUT | `/api/v2/contacts/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_contact) |
| `contacts delete` | DELETE | `/api/v2/contacts/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_contact) |
| `contacts hard-delete` | DELETE | `/api/v2/contacts/{id}/hard_delete` | write | [reference](https://developers.freshdesk.com/api/#hard_delete_contact) |
| `contacts make-agent` | PUT | `/api/v2/contacts/{id}/make_agent` | write | [reference](https://developers.freshdesk.com/api/#make_agent) |
| `contacts restore` | PUT | `/api/v2/contacts/{id}/restore` | write | [reference](https://developers.freshdesk.com/api/#restore_contact) |
| `contacts send-invite` | PUT | `/api/v2/contacts/{id}/send_invite` | write | [reference](https://developers.freshdesk.com/api/#send_invite) |
| `contacts merge` | POST | `/api/v2/contacts/merge` | write | [reference](https://developers.freshdesk.com/api/#merge_contact) |
| `contacts exports create` | POST | `/api/v2/contacts/export` | write | [reference](https://developers.freshdesk.com/api/#export_contact) |
| `contacts exports get` | GET | `/api/v2/contacts/export/{id}` | read | [reference](https://developers.freshdesk.com/api/#export_contact) |
| `contacts imports list` | GET | `/api/v2/contacts/imports` | read | [reference](https://developers.freshdesk.com/api/#import_contact) |
| `contacts imports get` | GET | `/api/v2/contacts/imports/{id}` | read | [reference](https://developers.freshdesk.com/api/#import_contact) |
| `contacts imports create` | POST | `/api/v2/contacts/imports` | write | [reference](https://developers.freshdesk.com/api/#import_contact) |
| `contacts imports cancel` | PUT | `/api/v2/contacts/imports/{id}/cancel` | write | [reference](https://developers.freshdesk.com/api/#import_contact) |
| `contact-fields list` | GET | `/api/v2/contact_fields` | read | [reference](https://developers.freshdesk.com/api/#list_all_contact_fields) |
| `contact-fields get` | GET | `/api/v2/contact_fields/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_contact_field) |
| `contact-fields create` | POST | `/api/v2/contact_fields` | write | [reference](https://developers.freshdesk.com/api/#create_contact_field) |
| `contact-fields update` | PUT | `/api/v2/contact_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_contact_field) |
| `contact-fields delete` | DELETE | `/api/v2/contact_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_contact_field) |
| `agents get` | GET | `/api/v2/agents/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_agent) |
| `agents list` | GET | `/api/v2/agents` | read | [reference](https://developers.freshdesk.com/api/#list_all_agents) |
| `agents create` | POST | `/api/v2/agents/` | write | [reference](https://developers.freshdesk.com/api/#create_agent) |
| `agents update` | PUT | `/api/v2/agents/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_agent) |
| `agents delete` | DELETE | `/api/v2/agents/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_agent) |
| `agents me` | GET | `/api/v2/agents/me` | read | [reference](https://developers.freshdesk.com/api/#me) |
| `agents bulk-create` | POST | `/api/v2/agents/bulk` | write | [reference](https://developers.freshdesk.com/api/#create_multiple_agents) |
| `agents autocomplete` | GET | `/api/v2/agents/autocomplete` | read | [reference](https://developers.freshdesk.com/api/#search_agents) |
| `agents availability list` | GET | `/api/v2/agents` | read | [reference](https://developers.freshdesk.com/api/#list_all_agents_availability) |
| `agents availability get` | GET | `/api/v2/agents/{id}/availability` | read | [reference](https://developers.freshdesk.com/api/#view_agent_availability) |
| `agents availability update` | PUT | `/api/v2/agents/{id}/availability` | write | [reference](https://developers.freshdesk.com/api/#update_agent_availability) |
| `skills get` | GET | `/api/v2/admin/skills/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_skill) |
| `skills list` | GET | `/api/v2/admin/skills/` | read | [reference](https://developers.freshdesk.com/api/#list_all_skills) |
| `skills create` | POST | `/api/v2/admin/skills` | write | [reference](https://developers.freshdesk.com/api/#create_skill) |
| `skills delete` | DELETE | `/api/v2/admin/skills/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_skill) |
| `skills update` | PUT | `/api/v2/admin/skills/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_skill) |
| `roles get` | GET | `/api/v2/roles/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_role) |
| `roles list` | GET | `/api/v2/roles` | read | [reference](https://developers.freshdesk.com/api/#list_all_roles) |
| `groups create` | POST | `/api/v2/groups` | write | [reference](https://developers.freshdesk.com/api/#create_group) |
| `groups get` | GET | `/api/v2/groups/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_group) |
| `groups list` | GET | `/api/v2/groups` | read | [reference](https://developers.freshdesk.com/api/#list_all_groups) |
| `groups update` | PUT | `/api/v2/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_group) |
| `groups delete` | DELETE | `/api/v2/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_group) |
| `admin-groups create` | POST | `/api/v2/admin/groups` | write | [reference](https://developers.freshdesk.com/api/#create_admin_group) |
| `admin-groups get` | GET | `/api/v2/admin/groups/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_admin_group) |
| `admin-groups list` | GET | `/api/v2/admin/groups` | read | [reference](https://developers.freshdesk.com/api/#list_all_admin_groups) |
| `admin-groups update` | PUT | `/api/v2/admin/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_admin_group) |
| `admin-groups delete` | DELETE | `/api/v2/admin/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_admin_group) |
| `admin-groups agents list` | GET | `/api/v2/admin/groups/{id}/agents` | read | [reference](https://developers.freshdesk.com/api/#list_all_agents_in_a_group) |
| `admin-groups agents update` | PATCH | `/api/v2/admin/groups/{id}/agents` | write | [reference](https://developers.freshdesk.com/api/#add_or_remove_agent_in_a_group) |
| `companies create` | POST | `/api/v2/companies` | write | [reference](https://developers.freshdesk.com/api/#create_company) |
| `companies get` | GET | `/api/v2/companies/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_company) |
| `companies list` | GET | `/api/v2/companies` | read | [reference](https://developers.freshdesk.com/api/#list_all_companies) |
| `companies autocomplete` | GET | `/api/v2/companies/autocomplete` | read | [reference](https://developers.freshdesk.com/api/#search_companies) |
| `companies search` | GET | `/api/v2/search/companies` | read | [reference](https://developers.freshdesk.com/api/#filter_companies) |
| `companies update` | PUT | `/api/v2/companies/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_company) |
| `companies delete` | DELETE | `/api/v2/companies/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_company) |
| `companies exports create` | POST | `/api/v2/companies/export` | write | [reference](https://developers.freshdesk.com/api/#export_company) |
| `companies exports get` | GET | `/api/v2/companies/export/{id}` | read | [reference](https://developers.freshdesk.com/api/#export_company) |
| `companies imports list` | GET | `/api/v2/companies/imports` | read | [reference](https://developers.freshdesk.com/api/#import_company) |
| `companies imports get` | GET | `/api/v2/companies/imports/{id}` | read | [reference](https://developers.freshdesk.com/api/#import_company) |
| `companies imports create` | POST | `/api/v2/companies/imports` | write | [reference](https://developers.freshdesk.com/api/#import_company) |
| `companies imports cancel` | PUT | `/api/v2/companies/imports/{id}/cancel` | write | [reference](https://developers.freshdesk.com/api/#import_company) |
| `company-fields list` | GET | `/api/v2/company_fields` | read | [reference](https://developers.freshdesk.com/api/#list_all_company_fields) |
| `company-fields get` | GET | `/api/v2/company_fields/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_company_field) |
| `company-fields create` | POST | `/api/v2/company_fields` | write | [reference](https://developers.freshdesk.com/api/#create_company_field) |
| `company-fields update` | PUT | `/api/v2/company_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_company_field) |
| `company-fields delete` | DELETE | `/api/v2/company_fields/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_company_field) |
| `custom-objects list` | GET | `/api/v2/custom_objects/schemas` | read | [reference](https://developers.freshdesk.com/api/#retrieve_existing_custom_objects) |
| `custom-objects get` | GET | `/api/v2/custom_objects/schemas/{schema_id}` | read | [reference](https://developers.freshdesk.com/api/#retrieve_specific_custom_object) |
| `custom-objects records create` | POST | `/api/v2/custom_objects/schemas/{schema_id}/records/` | write | [reference](https://developers.freshdesk.com/api/#create_a_new_custom_object_record) |
| `custom-objects records get` | GET | `/api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | read | [reference](https://developers.freshdesk.com/api/#retrieve_a_custom_object_record) |
| `custom-objects records update` | PUT | `/api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | write | [reference](https://developers.freshdesk.com/api/#update_a_custom_object_record) |
| `custom-objects records delete` | DELETE | `/api/v2/custom_objects/schemas/{schema_id}/records/{record_id}` | write | [reference](https://developers.freshdesk.com/api/#delete_a_custom_object_record) |
| `custom-objects records list` | GET | `/api/v2/custom_objects/schemas/{schema_id}/records` | read | [reference](https://developers.freshdesk.com/api/#retrieve_all_the_records_of_a_custom_object) |
| `custom-objects records count` | GET | `/api/v2/custom_objects/schemas/{schema_id}/records/count` | read | [reference](https://developers.freshdesk.com/api/#retrieve_the_count_of_records_custom_objects) |
| `custom-objects records filter` | GET | `/api/v2/custom_objects/schemas/{schema_id}/records` | read | [reference](https://developers.freshdesk.com/api/#filter_records_of_custom_objects) |
| `canned-responses create` | POST | `/api/v2/canned_responses` | write | [reference](https://developers.freshdesk.com/api/#create_canned_response) |
| `canned-responses get` | GET | `/api/v2/canned_responses/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_canned_response) |
| `canned-responses update` | PUT | `/api/v2/canned_responses/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_canned_response) |
| `canned-responses bulk-create` | POST | `/api/v2/canned_responses/bulk` | write | [reference](https://developers.freshdesk.com/api/#bulk_create_canned_responses) |
| `canned-responses folders create` | POST | `/api/v2/canned_response_folders` | write | [reference](https://developers.freshdesk.com/api/#create_canned_response_folder) |
| `canned-responses folders update` | PUT | `/api/v2/canned_response_folders/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_canned_response_folder) |
| `canned-responses folders list` | GET | `/api/v2/canned_response_folders` | read | [reference](https://developers.freshdesk.com/api/#list_all_canned_response_folders) |
| `canned-responses folders get` | GET | `/api/v2/canned_response_folders/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_canned_response_folders) |
| `canned-responses folders responses` | GET | `/api/v2/canned_response_folders/{id}/responses` | read | [reference](https://developers.freshdesk.com/api/#get_details_of_canned_responses_in_folders) |
| `discussions categories create` | POST | `/api/v2/discussions/categories` | write | [reference](https://developers.freshdesk.com/api/#create_category) |
| `discussions categories get` | GET | `/api/v2/discussions/categories/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_category) |
| `discussions categories list` | GET | `/api/v2/discussions/categories` | read | [reference](https://developers.freshdesk.com/api/#list_all_categories) |
| `discussions forums list` | GET | `/api/v2/discussions/categories/{id}/forums` | read | [reference](https://developers.freshdesk.com/api/#list_all_forums) |
| `discussions categories update` | PUT | `/api/v2/discussions/categories/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_category) |
| `discussions categories delete` | DELETE | `/api/v2/discussions/categories/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_category) |
| `discussions forums create` | POST | `/api/v2/discussions/categories/{category_id}/forums` | write | [reference](https://developers.freshdesk.com/api/#create_forum) |
| `discussions forums get` | GET | `/api/v2/discussions/forums/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_forum) |
| `discussions topics list` | GET | `/api/v2/discussions/forums/{id}/topics` | read | [reference](https://developers.freshdesk.com/api/#list_all_topics) |
| `discussions forums update` | PUT | `/api/v2/discussions/forums/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_forum) |
| `discussions forums delete` | DELETE | `/api/v2/discussions/forums/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_forum) |
| `discussions forums follow` | POST | `/api/v2/discussions/forums/{id}/follow` | write | [reference](https://developers.freshdesk.com/api/#monitor_forum) |
| `discussions forums unfollow` | DELETE | `/api/v2/discussions/forums/{forum_id}/follow` | write | [reference](https://developers.freshdesk.com/api/#unmonitor_forum) |
| `discussions forums follow-status` | GET | `/api/v2/discussions/forums/{id}/follow` | read | [reference](https://developers.freshdesk.com/api/#view_monitor_status_for_forum) |
| `discussions topics create` | POST | `/api/v2/discussions/forums/{forum_id}/topics` | write | [reference](https://developers.freshdesk.com/api/#create_topic) |
| `discussions topics get` | GET | `/api/v2/discussions/topics/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_topic) |
| `discussions comments list` | GET | `/api/v2/discussions/topics/{id}/comments` | read | [reference](https://developers.freshdesk.com/api/#list_all_comments) |
| `discussions topics update` | PUT | `/api/v2/discussions/topics/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_topic) |
| `discussions topics delete` | DELETE | `/api/v2/discussions/topics/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_topic) |
| `discussions topics follow` | POST | `/api/v2/discussions/topics/{id}/follow` | write | [reference](https://developers.freshdesk.com/api/#monitor_topic) |
| `discussions topics unfollow` | DELETE | `/api/v2/discussions/topics/{topic_id}/follow` | write | [reference](https://developers.freshdesk.com/api/#unmonitor_topic) |
| `discussions topics follow-status` | GET | `/api/v2/discussions/topics/{id}/follow` | read | [reference](https://developers.freshdesk.com/api/#view_monitor_status_for_topic) |
| `discussions topics followed` | GET | `/api/v2/discussions/topics/followed_by` | read | [reference](https://developers.freshdesk.com/api/#list_all_monitored_topics) |
| `discussions topics participated` | GET | `/api/v2/discussions/topics/participated_by` | read | [reference](https://developers.freshdesk.com/api/#list_all_participated_topics) |
| `discussions comments create` | POST | `/api/v2/discussions/topics/{topic_id}/comments` | write | [reference](https://developers.freshdesk.com/api/#create_comment) |
| `discussions comments update` | PUT | `/api/v2/discussions/comments/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_comment) |
| `discussions comments delete` | DELETE | `/api/v2/discussions/comments/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_comment) |
| `solutions categories create` | POST | `/api/v2/solutions/categories` | write | [reference](https://developers.freshdesk.com/api/#create_solution_category) |
| `solutions categories translations create` | POST | `/api/v2/solutions/categories/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#create_solution_category) |
| `solutions categories update` | PUT | `/api/v2/solutions/categories/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_category) |
| `solutions categories translations update` | PUT | `/api/v2/solutions/categories/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_category) |
| `solutions categories get` | GET | `/api/v2/solutions/categories/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_category) |
| `solutions categories translations get` | GET | `/api/v2/solutions/categories/{id}/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_category) |
| `solutions categories list` | GET | `/api/v2/solutions/categories` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_categories) |
| `solutions categories translations list` | GET | `/api/v2/solutions/categories/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_categories) |
| `solutions categories delete` | DELETE | `/api/v2/solutions/categories/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_solution_category) |
| `solutions folders create` | POST | `/api/v2/solutions/categories/{id}/folders` | write | [reference](https://developers.freshdesk.com/api/#create_solution_folder) |
| `solutions folders translations create` | POST | `/api/v2/solutions/folders/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#create_solution_folder) |
| `solutions folders update` | PUT | `/api/v2/solutions/folders/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_folder) |
| `solutions folders translations update` | PUT | `/api/v2/solutions/folders/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_folder) |
| `solutions folders get` | GET | `/api/v2/solutions/folders/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_folder) |
| `solutions folders translations get` | GET | `/api/v2/solutions/folders/{id}/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_folder) |
| `solutions folders list` | GET | `/api/v2/solutions/categories/{category_id}/folders` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_folders) |
| `solutions folders translations list` | GET | `/api/v2/solutions/categories/{category_id}/folders/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_folders) |
| `solutions folders subfolders` | GET | `/api/v2/solutions/folders/{folder_id}/subfolders` | read | [reference](https://developers.freshdesk.com/api/#list_all_sub_folders) |
| `solutions folders translations subfolders` | GET | `/api/v2/solutions/folders/{folder_id}/subfolders/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#list_all_sub_folders) |
| `solutions folders delete` | DELETE | `/api/v2/solutions/folders/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_solution_folder) |
| `solutions articles create` | POST | `/api/v2/solutions/folders/{id}/articles` | write | [reference](https://developers.freshdesk.com/api/#create_solution_article) |
| `solutions articles translations create` | POST | `/api/v2/solutions/articles/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#create_solution_article) |
| `solutions articles update` | PUT | `/api/v2/solutions/articles/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_article) |
| `solutions articles translations update` | PUT | `/api/v2/solutions/articles/{id}/{language_code}` | write | [reference](https://developers.freshdesk.com/api/#update_solution_article) |
| `solutions articles get` | GET | `/api/v2/solutions/articles/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_article) |
| `solutions articles translations get` | GET | `/api/v2/solutions/articles/{id}/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#view_solution_article) |
| `solutions articles list` | GET | `/api/v2/solutions/folders/{id}/articles` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_articles) |
| `solutions articles translations list` | GET | `/api/v2/solutions/folders/{id}/articles/{language_code}` | read | [reference](https://developers.freshdesk.com/api/#list_all_solution_articles) |
| `solutions articles delete` | DELETE | `/api/v2/solutions/articles/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_solution_article) |
| `solutions articles search` | GET | `/api/v2/search/solutions` | read | [reference](https://developers.freshdesk.com/api/#search_solution_article) |
| `csat surveys list` | GET | `/api/v2/customer-satisfaction/surveys` | read | [reference](https://developers.freshdesk.com/api/#csat_surveys) |
| `csat responses create` | POST | `/api/v2/customer-satisfaction/surveys/{survey_id}/responses` | write | [reference](https://developers.freshdesk.com/api/#create_survey_responses) |
| `csat responses list` | GET | `/api/v2/customer-satisfaction/surveys/{survey_id}/responses` | read | [reference](https://developers.freshdesk.com/api/#view_survey_responses) |
| `surveys list` | GET | `/api/v2/surveys` | read | [reference](https://developers.freshdesk.com/api/#list_all_survey) |
| `tickets ratings create` | POST | `/api/v2/tickets/{ticket_id}/satisfaction_ratings` | write | [reference](https://developers.freshdesk.com/api/#create_satisfaction_rating) |
| `surveys ratings list` | GET | `/api/v2/surveys/satisfaction_ratings` | read | [reference](https://developers.freshdesk.com/api/#view_all_satisfaction_ratings) |
| `service tasks create` | POST | `/api/v2/tickets` | write | [reference](https://developers.freshdesk.com/api/#create_service_task) |
| `service tasks get` | GET | `/api/v2/tickets/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_service_task) |
| `service tasks update` | PUT | `/api/v2/tickets/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_service_task) |
| `service tasks delete` | DELETE | `/api/v2/tickets/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_service_task) |
| `service groups create` | POST | `/api/v2/groups` | write | [reference](https://developers.freshdesk.com/api/#create_service_group) |
| `service groups get` | GET | `/api/v2/groups/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_service_group) |
| `service groups update` | PUT | `/api/v2/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_service_group) |
| `service groups delete` | DELETE | `/api/v2/groups/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_service_group) |
| `service make-technician` | PUT | `/api/v2/contacts/{id}/make_agent` | write | [reference](https://developers.freshdesk.com/api/#make_field_technician) |
| `time-entries create` | POST | `/api/v2/tickets/{ticket_id}/time_entries` | write | [reference](https://developers.freshdesk.com/api/#create_time_entry) |
| `time-entries list` | GET | `/api/v2/time_entries` | read | [reference](https://developers.freshdesk.com/api/#list_all_time_entries) |
| `time-entries update` | PUT | `/api/v2/time_entries/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_time_entry) |
| `time-entries toggle` | PUT | `/api/v2/time_entries/{time_entry_id}/toggle_timer` | write | [reference](https://developers.freshdesk.com/api/#toggle_time_entry) |
| `time-entries delete` | DELETE | `/api/v2/time_entries/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_time_entry) |
| `email-configs get` | GET | `/api/v2/email_configs/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_email_config) |
| `email-configs list` | GET | `/api/v2/email_configs` | read | [reference](https://developers.freshdesk.com/api/#list_all_email_configs) |
| `email mailboxes create` | POST | `/api/v2/email/mailboxes` | write | [reference](https://developers.freshdesk.com/api/#create_email_mailbox) |
| `email mailboxes get` | GET | `/api/v2/email/mailboxes/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_email_mailbox) |
| `email mailboxes list` | GET | `/api/v2/email/mailboxes` | read | [reference](https://developers.freshdesk.com/api/#list_all_email_mailboxes) |
| `email mailboxes update` | PUT | `/api/v2/email/mailboxes/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_email_mailbox) |
| `email mailboxes delete` | DELETE | `/api/v2/email/mailboxes/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_email_mailbox) |
| `email settings update` | PUT | `/api/v2/email/settings` | write | [reference](https://developers.freshdesk.com/api/#update_mailbox_settings) |
| `email settings get` | GET | `/api/v2/email/settings` | read | [reference](https://developers.freshdesk.com/api/#view_mailbox_settings) |
| `email bcc update` | PUT | `/api/v2/notifications/email/bcc` | write | [reference](https://developers.freshdesk.com/api/#update_automatic_bcc_emails) |
| `email bcc get` | GET | `/api/v2/notifications/email/bcc` | read | [reference](https://developers.freshdesk.com/api/#view_automatic_bcc_emails) |
| `products get` | GET | `/api/v2/products/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_product) |
| `products list` | GET | `/api/v2/products` | read | [reference](https://developers.freshdesk.com/api/#list_all_products) |
| `business-hours get` | GET | `/api/v2/business_hours/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_business_hour) |
| `business-hours list` | GET | `/api/v2/business_hours` | read | [reference](https://developers.freshdesk.com/api/#list_all_business_hours) |
| `scenario-automations list` | GET | `/api/v2/scenario_automations` | read | [reference](https://developers.freshdesk.com/api/#list_all_scenario_automations) |
| `sla-policies create` | POST | `/api/v2/sla_policies` | write | [reference](https://developers.freshdesk.com/api/#create_sla_policy) |
| `sla-policies list` | GET | `/api/v2/sla_policies` | read | [reference](https://developers.freshdesk.com/api/#list_all_sla_policies) |
| `sla-policies update` | PUT | `/api/v2/sla_policies/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_sla_policy) |
| `activities publish` | POST | `/api/v2/contact-activities` | write | [reference](https://developers.freshdesk.com/api/#publish_omnichannel_activity) |
| `automations list` | GET | `/api/v2/automations/{automation_type_id}/rules` | read | [reference](https://developers.freshdesk.com/api/#list_all_automation_rules) |
| `automations get` | GET | `/api/v2/automations/{automation_type_id}/rules/{id}` | read | [reference](https://developers.freshdesk.com/api/#view_an_automation_rule) |
| `automations delete` | DELETE | `/api/v2/automations/{automation_type_id}/rules/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_an_automation_rule) |
| `automations create` | POST | `/api/v2/automations/{automation_type_id}/rules` | write | [reference](https://developers.freshdesk.com/api/#create_an_automation_rule) |
| `automations update` | PUT | `/api/v2/automations/{automation_type_id}/rules/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_an_automation_rule) |
| `settings get` | GET | `/api/v2/settings/helpdesk` | read | [reference](https://developers.freshdesk.com/api/#settings_helpdesk) |
| `threads create` | POST | `/api/v2/collaboration/threads` | write | [reference](https://developers.freshdesk.com/api/#create_a_thread) |
| `threads get` | GET | `/api/v2/collaboration/threads/{id}` | read | [reference](https://developers.freshdesk.com/api/#get_a_thread) |
| `threads update` | PUT | `/api/v2/collaboration/threads/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_a_thread) |
| `threads delete` | DELETE | `/api/v2/collaboration/threads/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_a_thread) |
| `threads messages create` | POST | `/api/v2/collaboration/messages` | write | [reference](https://developers.freshdesk.com/api/#create_message_for_thread) |
| `threads messages get` | GET | `/api/v2/collaboration/messages/{id}` | read | [reference](https://developers.freshdesk.com/api/#get_message_for_thread) |
| `threads messages update` | PUT | `/api/v2/collaboration/messages/{id}` | write | [reference](https://developers.freshdesk.com/api/#update_message_for_thread) |
| `threads messages delete` | DELETE | `/api/v2/collaboration/messages/{id}` | write | [reference](https://developers.freshdesk.com/api/#delete_message_for_thread) |
| `threads messages quote` | GET | `/api/v2/collaboration/messages/{id}/generate-quote` | read | [reference](https://developers.freshdesk.com/api/#get_quoted_text) |
| `audit-log exports create` | POST | `/api/v2/audit_log/export` | write | [reference](https://developers.freshdesk.com/api/#export_audit_log) |
| `audit-log exports get` | GET | `/api/v2/audit_log/export/{id}` | read | [reference](https://developers.freshdesk.com/api/#export_audit_log) |

## Reference discrepancies

These are deliberate choices from the public reference. They require tenant verification before claiming live compatibility.

- `ticket-forms clone`: The reference's curl example includes /clone; the heading omits it. GET has a write effect.

- `ticket-forms fields update`: Reference specifies POST for updating a form field.

- `contact-fields delete`: The curl example uses plural contact_fields; the heading is singular.

- `skills update`: Reference labels the verb 'update'; its curl example specifies PUT.

- `company-fields update`: The curl example uses plural company_fields; the heading is singular.

- `company-fields delete`: Use plural company_fields, consistent with the reference's curl example.

- `canned-responses bulk-create`: Reference heading specifies /canned_responses/bulk; its curl example has a duplicated /api/v2 and a different suffix. This command uses the heading; live verification is pending.

- `canned-responses folders create`: The curl example uses plural canned_response_folders; the heading is singular.

- `canned-responses folders update`: The curl example uses plural canned_response_folders; the heading is singular.

- `scenario-automations list`: The heading shows a legacy .json route; the curl example uses /api/v2/scenario_automations.

## Boundaries

- The reference documents creating a custom-object schema in the admin UI, but provides no schema-creation REST endpoint. Record CRUD, schema reads, and record filtering are included.
- Roles, products, business hours, and helpdesk settings expose read-only operations in this reference.
- Analytics dashboards, portal themes, billing, Freddy configuration, and app deployment are outside this REST API inventory. FDK remains the app-development tool.
- Custom-object APIs use their own page-size/cursor protocol. Use `--param` and the returned native metadata; `--all` is not offered for these operations.
- The new CSAT APIs use `/customer-satisfaction`; the legacy survey APIs remain separate. Some contact/company search APIs are marked beta upstream.
- JSON and `--set` can express endpoint payloads; the server validates account-specific fields and business rules. There is no claim of complete offline validation for every payload.
- Bulk and export/import requests can return asynchronous jobs. Inspect the returned job/export/import ID with the corresponding get command.
- Attachment upload is included through `--attach` and `--upload`. Download URLs are native fields in ticket/conversation responses; download them with an HTTPS client. The CLI does not save remote attachment files automatically.
