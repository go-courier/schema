package jsonschema_test

import (
	"testing"

	"github.com/go-courier/schema/pkg/jsonschema"
	"github.com/go-courier/schema/pkg/testutil"

	"github.com/go-courier/schema/testdata/a"
)

func TestFromGoValue(t *testing.T) {
	p := a.Struct{}

	testutil.PrintJSON(jsonschema.FromGoValue(p))
}
