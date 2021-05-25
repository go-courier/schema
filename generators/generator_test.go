package generators

import (
	"context"
	"testing"

	"github.com/go-courier/gengo/pkg/gengo"

	// ugly hack for testdata
	_ "github.com/gorilla/handlers"
)

func TestGenerator(t *testing.T) {
	t.Skip()

	c, err := gengo.NewContext(&gengo.GeneratorArgs{
		Entrypoint: []string{
			"github.com/go-courier/schema/testdata",
		},
		OutputFileBaseName: "zz_generated",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Execute(context.Background(), gengo.GetRegisteredGenerators()...); err != nil {
		t.Fatal(err)
	}
}
