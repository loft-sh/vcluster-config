package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"
)

const (
	platformConfigName = "PlatformConfig"
)

func RunMergeSchemas(valuesSchemaFile, platformConfigSchemaFile, outFile string) error {
	platformBytes, err := os.ReadFile(platformConfigSchemaFile)
	if err != nil {
		return err
	}
	platformSchema := &jsonschema.Schema{}
	err = json.Unmarshal(platformBytes, platformSchema)
	if err != nil {
		return err
	}
	valuesBytes, err := os.ReadFile(valuesSchemaFile)
	if err != nil {
		return err
	}
	valuesSchema := &jsonschema.Schema{}
	err = json.Unmarshal(valuesBytes, valuesSchema)
	if err != nil {
		return err
	}
	if err := addPlatformSchema(platformSchema, valuesSchema); err != nil {
		return err
	}
	return writeSchema(valuesSchema, outFile)
}

func addPlatformSchema(platformSchema, toSchema *jsonschema.Schema) error {
	for defName, node := range platformSchema.Definitions {
		if _, exists := toSchema.Definitions[defName]; exists {
			panic("trying to overwrite definition " + defName + " this is unexpected")
		}
		toSchema.Definitions[defName] = node
	}

	for defName, def := range toSchema.Definitions {
		if defName == platformConfigName {
			for pair := platformSchema.Properties.Oldest(); pair != nil; pair = pair.Next() {
				pair := pair
				def.Properties.AddPairs(*pair)
			}
		}
	}
	return nil
}

func writeSchema(schema *jsonschema.Schema, schemaFile string) error {
	prefix := ""
	schemaString, err := json.MarshalIndent(schema, prefix, "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(schemaFile), os.ModePerm)
	if err != nil {
		return err
	}
	if _, err = os.Create(schemaFile); err != nil {
		return err
	}
	return os.WriteFile(schemaFile, schemaString, os.ModePerm)
}
