package pkg

import (
	"github.com/xeipuuv/gojsonschema"
	"os"
	"reflect"
	"testing"
)

func TestRunMergeSchemas(t *testing.T) {
	cases := []struct {
		valuesSchema   string
		platformSchema string
		expected       string
	}{
		{
			valuesSchema:   "testdata/values.schema.json",
			platformSchema: "testdata/platform.schema.json",
			expected:       "testdata/vcluster.schema.json",
		},
	}

	for _, tc := range cases {
		t.Run("latest", func(t *testing.T) {
			outFile, err := os.CreateTemp("", "got_vcluster.schema.json")
			assertNoError(t, err)
			t.Logf("created merged file: %s\n", outFile.Name())
			err = RunMergeSchemas(tc.valuesSchema, tc.platformSchema, outFile.Name())
			assertNoError(t, err)
			expected, err := os.ReadFile(tc.expected)
			assertNoError(t, err)
			got, err := os.ReadFile(outFile.Name())
			assertNoError(t, err)
			if !reflect.DeepEqual(expected, got) {
				t.Fatalf("expected merged schema as %s got %s\n", tc.expected, outFile.Name())
			}
			cwd, err := os.Getwd()
			assertNoError(t, err)
			schemaLoader := gojsonschema.NewReferenceLoader("file://" + cwd + "/testdata/vcluster.schema.json")
			exampleConfig := `
{
	"external": {
		"platform": {
			"autoSleep": {
				"afterInactivity": 100
			},
      
    		"apiKey": {
				"secretName": "foo",
      			"namespace": "bar"
			}
		}
	}
}
`

			example := gojsonschema.NewStringLoader(exampleConfig)
			_, err = gojsonschema.Validate(schemaLoader, example)
			assertNoError(t, err)

		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
