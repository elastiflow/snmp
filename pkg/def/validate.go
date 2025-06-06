package def

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schemas/*.json
var schemas embed.FS

func validate[T any](definition T, schemaName string) error {
	schemaContent, err := schemas.ReadFile(schemaName)
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	schema, err := jsonschema.CompileString(schemaName, string(schemaContent))
	if err != nil {
		return fmt.Errorf("failed to add schema resource: %w", err)
	}

	data, err := json.Marshal(definition)
	if err != nil {
		return fmt.Errorf("failed to marshal definition: %w", err)
	}

	var v interface{}
	if unmarshalErr := json.Unmarshal(data, &v); unmarshalErr != nil {
		return fmt.Errorf("failed to unmarshal definition: %w", err)
	}

	return schema.Validate(v)
}
