package pkg

import (
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	externalConfigName = "ExternalConfig"
	platformConfigName = "PlatformConfig"
	platformConfigRef  = "#/defs/" + platformConfigName
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

	addPlatformSchema(platformSchema, valuesSchema)
	return writeSchema(valuesSchema, outFile)
}

func addPlatformSchema(platformSchema, toSchema *jsonschema.Schema) {
	platformNode := &jsonschema.Schema{
		AdditionalProperties: nil,
		Description:          platformConfigName + " holds platform configuration",
		Properties:           platformSchema.Properties,
		Type:                 "object",
	}
	toSchema.Definitions[platformConfigName] = platformNode
	properties := jsonschema.NewProperties()
	properties.AddPairs(orderedmap.Pair[string, *jsonschema.Schema]{
		Key: "platform",
		Value: &jsonschema.Schema{
			Ref:         platformConfigRef,
			Description: "platform holds platform configuration",
			Type:        "object",
		},
	})
	externalConfigNode, ok := toSchema.Definitions[externalConfigName]
	if !ok {
		externalConfigNode = &jsonschema.Schema{
			AdditionalProperties: nil,
			Description:          externalConfigName + " holds external configuration",
			Properties:           properties,
		}
	} else {
		externalConfigNode.Properties = properties
	}
	toSchema.Definitions[externalConfigName] = externalConfigNode

	for defName, node := range platformSchema.Definitions {
		if _, exists := toSchema.Definitions[defName]; exists {
			panic("trying to overwrite definition " + defName + " this is unexpected")
		}
		toSchema.Definitions[defName] = node
	}
}

func writeSchema(schema *jsonschema.Schema, schemaFile string) error {
	prefix := ""
	schemaString, err := json.MarshalIndent(schema, prefix, "  ")
	if err != nil {
		return err
	}

	//err = os.MkdirAll(filepath.Dir(schemaFile), os.ModePerm)
	//if err != nil {
	//	return err
	//}
	//if _, err = os.Create(schemaFile); err != nil {
	//	return err
	//}

	err = os.WriteFile(schemaFile, schemaString, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
