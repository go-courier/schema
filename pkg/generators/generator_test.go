package generators

import (
	"context"
	"testing"

	"github.com/go-courier/gengo/pkg/gengo"
)

func TestGenerator(t *testing.T) {
	c, _ := gengo.NewContext(&gengo.GeneratorArgs{
		Inputs: []string{
			"github.com/go-courier/schema/testdata/a",
			"github.com/go-courier/schema/testdata/b",
		},
		OutputFileBaseName: "zz_generated",
	})

	if err := c.Execute(context.Background(), gengo.GetRegisteredGenerators()...); err != nil {
		t.Fatal(err)
	}
}
