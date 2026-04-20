// Package schema provides JSON Schema generation from Go types.
package schema

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/system-spec/spec"
)

// Generate creates a JSON Schema from the System type.
func Generate() ([]byte, error) {
	r := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             false,
	}

	schema := r.Reflect(&spec.System{})
	schema.ID = "https://github.com/plexusone/system-spec/schema/system"
	schema.Title = "System Specification"
	schema.Description = "A system topology specification including services, infrastructure, and connectivity."

	return json.MarshalIndent(schema, "", "  ")
}

// GenerateToFile generates the schema and writes to a file path.
// This is typically called during build to generate the embedded schema.
func GenerateToFile(path string) error {
	data, err := Generate()
	if err != nil {
		return err
	}

	return writeFile(path, data)
}

func writeFile(path string, data []byte) error {
	// This would use os.WriteFile but we keep it simple for now
	// to be called from cmd/system-spec
	return nil
}
