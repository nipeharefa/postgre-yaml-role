package yaml

import (
	"strings"
	"testing"

	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
)

func TestParse(t *testing.T) {
	yamlData := `apiVersion: pgctl/v1
kind: Role
metadata:
  name: test-role
spec:
  login: true
  inherit: true
  createrole: false
  createdb: false
  replication: false
  bypassrls: false
  connectionLimit: 5
  comment: "A test role"
`
	resources, err := Parse(strings.NewReader(yamlData))
	if err != nil {
		t.Fatalf("Error parsing YAML: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(resources))
	}

	role, ok := resources[0].(types.Role)
	if !ok {
		t.Fatalf("Expected Role resource, got %T", resources[0])
	}

	if role.Metadata.Name != "test-role" {
		t.Errorf("Expected role name 'test-role', got '%s'", role.Metadata.Name)
	}

	if !role.Spec.Login {
		t.Error("Expected Login to be true")
	}
}