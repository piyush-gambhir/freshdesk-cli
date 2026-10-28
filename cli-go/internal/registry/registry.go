// Package registry is the reviewed inventory behind commands and API discovery.
package registry

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed operations.json
var raw []byte

type Operation struct {
	ID             string            `json:"id"`
	Command        []string          `json:"command"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	PathParameters []string          `json:"path_parameters"`
	Effect         string            `json:"effect"`
	Paginated      bool              `json:"paginated"`
	Pagination     string            `json:"pagination,omitempty"`
	Confirmation   bool              `json:"confirmation,omitempty"`
	QueryDefaults  map[string]string `json:"query_defaults,omitempty"`
	RequiredQuery  []string          `json:"required_query,omitempty"`
	BodyDefaults   map[string]any    `json:"body_defaults,omitempty"`
	SourceAnchor   string            `json:"source_anchor"`
	SourcePath     string            `json:"source_path"`
	Note           string            `json:"note,omitempty"`
}

var Operations = mustLoad()

func mustLoad() []Operation {
	var ops []Operation
	if err := json.Unmarshal(raw, &ops); err != nil {
		panic(err)
	}
	seen := map[string]bool{}
	for _, op := range ops {
		if seen[op.ID] || op.ID == "" || (op.Effect != "read" && op.Effect != "write") {
			panic(fmt.Sprintf("invalid registry operation %s", op.ID))
		}
		seen[op.ID] = true
	}
	return ops
}

func Hash() string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func Find(id string) (Operation, bool) {
	for _, op := range Operations {
		if op.ID == id || strings.Join(op.Command, " ") == id {
			return op, true
		}
	}
	return Operation{}, false
}
