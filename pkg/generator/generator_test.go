package generator

import (
	"k8s.io/gengo/args"
	"testing"

	enumerationgenerators "github.com/go-courier/schema/pkg/enumeration/generators"
	jsonschemagenerators "github.com/go-courier/schema/pkg/jsonschema/generators"
)

func TestPkgGenerator(t *testing.T) {
	arguments := args.Default()

	arguments.InputDirs = []string{
		"github.com/go-courier/schema/testdata/a",
		"github.com/go-courier/schema/testdata/b",
	}

	arguments.OutputFileBaseName = "zz_generated"

	g := PkgGenerators{
		enumerationgenerators.NewEnumGen,
		jsonschemagenerators.NewJsonSchemaGen,
	}

	if err := g.Execute(arguments); err != nil {
		t.Fatalf("Error: %v", err)
	}
}
