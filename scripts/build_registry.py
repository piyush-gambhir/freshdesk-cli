#!/usr/bin/env python3
"""Build the reviewed CLI registry from the pinned endpoint-facts snapshot."""
import json
import re
from pathlib import Path
from urllib.parse import parse_qsl

ROOT = Path(__file__).resolve().parents[1]
BASE = ROOT / "cli-go/internal/registry"

# One line per reference entry, in reference order. Translation variants and
# service-management variants deliberately have their own commands.
COMMANDS = """
tickets create
tickets outbound-email
tickets get
tickets associated
tickets prime-association
tickets list
tickets search
tickets update
tickets bulk-update
tickets forward
tickets merge
tickets watchers list
tickets watch
tickets unwatch
tickets bulk-watch
tickets bulk-unwatch
tickets delete
tickets bulk-delete
attachments delete
tickets restore
tickets fields
tickets conversations
tickets time-entries
tickets ratings
tickets summary get
tickets summary update
tickets summary delete
tickets archived get
tickets archived delete
tickets archived conversations
tickets access create
tickets access update
tickets access list
tickets access delete
ticket-fields list
ticket-fields create
ticket-fields get
ticket-fields update
ticket-fields delete
ticket-fields sections list
ticket-fields sections create
ticket-fields sections get
ticket-fields sections update
ticket-fields sections delete
ticket-forms list
ticket-forms create
ticket-forms get
ticket-forms update
ticket-forms delete
ticket-forms clone
ticket-forms fields get
ticket-forms fields update
ticket-forms fields delete
tickets reply
tickets note
conversations update
conversations delete
tickets reply-to-forward
outbound-messages send
outbound-messages get
account get
account export
jobs get
contacts create
contacts get
contacts list
contacts autocomplete
contacts search
contacts update
contacts delete
contacts hard-delete
contacts make-agent
contacts restore
contacts send-invite
contacts merge
contacts exports create
contacts exports get
contacts imports list
contacts imports get
contacts imports create
contacts imports cancel
contact-fields list
contact-fields get
contact-fields create
contact-fields update
contact-fields delete
agents get
agents list
agents create
agents update
agents delete
agents me
agents bulk-create
agents autocomplete
agents availability list
agents availability get
agents availability update
skills get
skills list
skills create
skills delete
skills update
roles get
roles list
groups create
groups get
groups list
groups update
groups delete
admin-groups create
admin-groups get
admin-groups list
admin-groups update
admin-groups delete
admin-groups agents list
admin-groups agents update
companies create
companies get
companies list
companies autocomplete
companies search
companies update
companies delete
companies exports create
companies exports get
companies imports list
companies imports get
companies imports create
companies imports cancel
company-fields list
company-fields get
company-fields create
company-fields update
company-fields delete
custom-objects list
custom-objects get
custom-objects records create
custom-objects records get
custom-objects records update
custom-objects records delete
custom-objects records list
custom-objects records count
custom-objects records filter
canned-responses create
canned-responses get
canned-responses update
canned-responses bulk-create
canned-responses folders create
canned-responses folders update
canned-responses folders list
canned-responses folders get
canned-responses folders responses
discussions categories create
discussions categories get
discussions categories list
discussions forums list
discussions categories update
discussions categories delete
discussions forums create
discussions forums get
discussions topics list
discussions forums update
discussions forums delete
discussions forums follow
discussions forums unfollow
discussions forums follow-status
discussions topics create
discussions topics get
discussions comments list
discussions topics update
discussions topics delete
discussions topics follow
discussions topics unfollow
discussions topics follow-status
discussions topics followed
discussions topics participated
discussions comments create
discussions comments update
discussions comments delete
solutions categories create
solutions categories translations create
solutions categories update
solutions categories translations update
solutions categories get
solutions categories translations get
solutions categories list
solutions categories translations list
solutions categories delete
solutions folders create
solutions folders translations create
solutions folders update
solutions folders translations update
solutions folders get
solutions folders translations get
solutions folders list
solutions folders translations list
solutions folders subfolders
solutions folders translations subfolders
solutions folders delete
solutions articles create
solutions articles translations create
solutions articles update
solutions articles translations update
solutions articles get
solutions articles translations get
solutions articles list
solutions articles translations list
solutions articles delete
solutions articles search
csat surveys list
csat responses create
csat responses list
surveys list
tickets ratings create
surveys ratings list
service tasks create
service tasks get
service tasks update
service tasks delete
service groups create
service groups get
service groups update
service groups delete
service make-technician
time-entries create
time-entries list
time-entries update
time-entries toggle
time-entries delete
email-configs get
email-configs list
email mailboxes create
email mailboxes get
email mailboxes list
email mailboxes update
email mailboxes delete
email settings update
email settings get
email bcc update
email bcc get
products get
products list
business-hours get
business-hours list
scenario-automations list
sla-policies create
sla-policies list
sla-policies update
activities publish
automations list
automations get
automations delete
automations create
automations update
settings get
threads create
threads get
threads update
threads delete
threads messages create
threads messages get
threads messages update
threads messages delete
threads messages quote
audit-log exports create
audit-log exports get
""".strip().splitlines()

# Corrections are based on the same reference's curl examples. Where those
# disagree as well, retain the route heading and expose the discrepancy.
CORRECTIONS = {
    "clone_ticket_form": ("/api/v2/ticket-forms/{id}/clone", "The reference's curl example includes /clone; the heading omits it. GET has a write effect."),
    "delete_contact_field": ("/api/v2/contact_fields/{id}", "The curl example uses plural contact_fields; the heading is singular."),
    "update_company_field": ("/api/v2/company_fields/{id}", "The curl example uses plural company_fields; the heading is singular."),
    "delete_company_field": ("/api/v2/company_fields/{id}", "Use plural company_fields, consistent with the reference's curl example."),
    "create_canned_response_folder": ("/api/v2/canned_response_folders", "The curl example uses plural canned_response_folders; the heading is singular."),
    "update_canned_response_folder": ("/api/v2/canned_response_folders/{id}", "The curl example uses plural canned_response_folders; the heading is singular."),
    "list_all_scenario_automations": ("/api/v2/scenario_automations", "The heading shows a legacy .json route; the curl example uses /api/v2/scenario_automations."),
}


def main():
    snapshot = json.loads((BASE / "testdata/freshdesk-api.json").read_text())
    assert len(COMMANDS) == len(snapshot["operations"]), "Review the command mapping after refreshing the snapshot"
    registry = []
    for command, source in zip(COMMANDS, snapshot["operations"]):
        path, _, query = source["path"].partition("?")
        path = "/" + path.lstrip("/")
        if path.startswith("/canned_response"):
            path = "/api/v2" + path
        path = re.sub(r"\[([^]]+)\]", r"{\1}", path)
        path = path.replace("/schema-id/", "/{schema_id}/").replace("/schema_id/", "/{schema_id}/")
        path = re.sub(r"\{([^}]+)\}", lambda m: "{" + m[1].replace("-", "_") + "}", path)
        note = ""
        if source["anchor"] in CORRECTIONS:
            path, note = CORRECTIONS[source["anchor"]]
        if source["anchor"] == "bulk_create_canned_responses":
            note = "Reference heading specifies /canned_responses/bulk; its curl example has a duplicated /api/v2 and a different suffix. This command uses the heading; live verification is pending."
        if source["anchor"] == "update_skill":
            note = "Reference labels the verb 'update'; its curl example specifies PUT."
        if source["anchor"] == "update_ticket_form_field":
            note = "Reference specifies POST for updating a form field."
        defaults, required_query = {}, []
        for key, val in parse_qsl(query):
            if "[" in val or "{" in val:
                required_query.append(key)
            else:
                defaults[key] = val
        params = re.findall(r"\{([^}]+)\}", path)
        effect = "write" if source["method"] != "GET" or command == "ticket-forms clone" else "read"
        action = command.split()[-1]
        paginated = effect == "read" and action in {"list", "conversations", "time-entries", "ratings", "fields", "responses", "associated", "subfolders", "search", "followed", "participated", "filter"}
        # These APIs have a different cursor/page_size protocol. Raw query
        # parameters remain available, but --all must not guess that protocol.
        if command.startswith("custom-objects") or "autocomplete" in command or command == "solutions articles search":
            paginated = False
        # The admin ticket-field endpoint rejects page and per_page.
        if command == "ticket-fields list":
            paginated = False
        body_defaults = {}
        if command in {"service tasks create", "service tasks update"}:
            body_defaults = {"type": "Service Task"}
        if command == "service groups create":
            body_defaults = {"group_type": "field_agent_group"}
        if command == "service make-technician":
            body_defaults = {"type": "field_agent"}
        operation = {
            "id": "freshdesk." + command.replace(" ", "."), "command": command.split(),
            "method": source["method"], "path": path, "path_parameters": params,
            "effect": effect, "paginated": paginated,
            "source_anchor": source["anchor"], "source_path": source["path"],
        }
        if paginated:
            operation["pagination"] = "cursor" if command.startswith("csat ") else "page"
        if defaults:
            operation["query_defaults"] = defaults
        if required_query:
            operation["required_query"] = required_query
        if body_defaults:
            operation["body_defaults"] = body_defaults
        if note:
            operation["note"] = note
        if effect == "write":
            operation["confirmation"] = True
        registry.append(operation)
    (BASE / "operations.json").write_text(json.dumps(registry, indent=2) + "\n")
    print(f"Wrote {len(registry)} reviewed operations")


if __name__ == "__main__":
    main()
