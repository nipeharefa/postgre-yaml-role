package yaml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
	"gopkg.in/yaml.v3"
)

// ParseFile parses a YAML file and returns the corresponding resource
func ParseFile(filename string) ([]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

// ParseDirectory parses all YAML files in a directory
func ParseDirectory(dir string) ([]interface{}, error) {
	var resources []interface{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is YAML
		ext := filepath.Ext(path)
		if ext == ".yaml" || ext == ".yml" {
			fileResources, err := ParseFile(path)
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", path, err)
			}
			resources = append(resources, fileResources...)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resources, nil
}

// Parse parses YAML from an io.Reader
func Parse(reader io.Reader) ([]interface{}, error) {
	decoder := yaml.NewDecoder(reader)
	var resources []interface{}

	for {
		var resource map[string]interface{}
		err := decoder.Decode(&resource)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Determine the kind of resource
		kind, ok := resource["kind"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'kind' field in YAML")
		}

		// Convert to the appropriate type
		switch kind {
		case "Role":
			var role types.Role
			if err := yaml.Unmarshal(resourceToYAML(resource), &role); err != nil {
				return nil, err
			}
			resources = append(resources, role)
		case "User":
			var user types.User
			if err := yaml.Unmarshal(resourceToYAML(resource), &user); err != nil {
				return nil, err
			}
			resources = append(resources, user)
		case "Grant":
			var grant types.Grant
			if err := yaml.Unmarshal(resourceToYAML(resource), &grant); err != nil {
				return nil, err
			}
			resources = append(resources, grant)
		default:
			return nil, fmt.Errorf("unknown resource kind: %s", kind)
		}
	}

	return resources, nil
}

// resourceToYAML converts a map to YAML bytes
func resourceToYAML(resource map[string]interface{}) []byte {
	bytes, _ := yaml.Marshal(resource)
	return bytes
}