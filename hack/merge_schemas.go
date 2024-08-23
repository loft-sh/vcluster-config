package main

import (
	"os"

	"encoding/json"
	"path/filepath"

	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const (
	PlatformConfigSchema = "platform.schema.json"
	ValuesSchema         = "values.schema.json"
	OutFile              = "vcluster.schema.json"
	externalConfigName   = "ExternalConfig"
	platformConfigName   = "PlatformConfig"
	platformConfigRef    = "#/defs/" + platformConfigName
)

func main() {
	platformBytes, err := os.ReadFile(PlatformConfigSchema)
	panicOnErr(err)
	platformSchema := &jsonschema.Schema{}
	err = json.Unmarshal(platformBytes, platformSchema)
	panicOnErr(err)

	valuesBytes, err := os.ReadFile(ValuesSchema)
	panicOnErr(err)
	valuesSchema := &jsonschema.Schema{}
	err = json.Unmarshal(valuesBytes, valuesSchema)
	panicOnErr(err)

	addPlatformSchema(platformSchema, valuesSchema)
	panicOnErr(writeSchema(valuesSchema, OutFile))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err.Error())
	}
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

	err = os.MkdirAll(filepath.Dir(schemaFile), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(schemaFile, schemaString, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
