package main

import (
	"github.com/loft-sh/vcluster-config/hack/pkg"
)

const (
	PlatformConfigSchema = "platform.schema.json"
	ValuesSchema         = "values.schema.json"
	OutFile              = "vcluster.schema.json"
)

func main() {
	err := pkg.RunMergeSchemas(ValuesSchema, PlatformConfigSchema, OutFile)
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
