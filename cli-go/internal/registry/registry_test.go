package registry

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestEveryDocumentedOperationIsMapped(t *testing.T) {
	data, err := os.ReadFile("testdata/freshdesk-api.json")
	if err != nil {
		t.Fatal(err)
	}
	var snapshot struct {
		Operations []struct{ Anchor, Method, Path string }
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Operations) != len(Operations) {
		t.Fatalf("snapshot %d operations, registry %d", len(snapshot.Operations), len(Operations))
	}
	used := map[string]bool{}
	for i, source := range snapshot.Operations {
		op := Operations[i]
		if op.SourceAnchor != source.Anchor || op.Method != source.Method || op.SourcePath != source.Path {
			t.Fatalf("unmapped source operation: %+v", source)
		}
		path := strings.Join(op.Command, " ")
		if used[path] {
			t.Fatalf("duplicate command: %s", path)
		}
		used[path] = true
		if !strings.HasPrefix(op.Path, "/api/v2/") && op.Path != "/scenario_automations.json" {
			t.Fatalf("unscoped path %s", op.Path)
		}
		if len(regexp.MustCompile(`\{[^}]+\}`).FindAllString(op.Path, -1)) != len(op.PathParameters) {
			t.Fatalf("bad placeholders for %s", op.ID)
		}
		if op.Method != "GET" && op.Effect != "write" {
			t.Fatalf("mutation is classified as read: %s", op.ID)
		}
	}
}

func TestCloneIsWriteDespiteGET(t *testing.T) {
	op, ok := Find("ticket-forms clone")
	if !ok || op.Path != "/api/v2/ticket-forms/{id}/clone" || op.Effect != "write" {
		t.Fatalf("unsafe clone: %+v", op)
	}
}
