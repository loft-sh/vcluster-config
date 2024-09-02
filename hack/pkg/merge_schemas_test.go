package pkg

import (
	"os"
	"reflect"
	"testing"
)

func TestRunMergeSchemas(t *testing.T) {
	cases := []struct {
		valuesSchema, version, outFile string
		expected                       string
	}{
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.14",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.15",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.16",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.17",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.18",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.19",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.20",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-alpha.21",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.1",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.2",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.3",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.4",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.5",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.6",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.0.0-beta.7",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.1.0-alpha.0",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.1.0-alpha.1",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.1.0-alpha.2",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.1.0-alpha.3",
		},
		{
			valuesSchema: "testdata/values.schema.json",
			version:      "v4.1.0-alpha.4",
		},
	}

	for _, tc := range cases {
		t.Run(tc.version, func(t *testing.T) {
			outFile, err := os.CreateTemp("", tc.version+"_vcluster.schema.json")
			assertNoError(t, err)
			expectedOutFileName := "testdata/" + tc.version + "_vcluster.schema.json"
			platformConfigSchemaFile := "testdata/" + tc.version + "_platform.schema.json"
			t.Logf("created merged file: %s\n", outFile.Name())
			err = RunMergeSchemas(tc.valuesSchema, platformConfigSchemaFile, outFile.Name())
			assertNoError(t, err)
			expected, err := os.ReadFile(expectedOutFileName)
			assertNoError(t, err)
			got, err := os.ReadFile(outFile.Name())
			assertNoError(t, err)
			if !reflect.DeepEqual(expected, got) {
				t.Fatalf("expected merged schema as %s got %s\n", expectedOutFileName, outFile.Name())
			}
		})
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
