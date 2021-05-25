package generator

import (
	"k8s.io/gengo/types"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/kube-openapi/pkg/util/sets"
)

type NewPkgGenerator = func(pkg *types.Package, generatorArgs *args.GeneratorArgs) generator.Generator

type PkgGenerators []NewPkgGenerator

func (g PkgGenerators) Execute(generatorArgs *args.GeneratorArgs) error {
	return generatorArgs.Execute(g.NameSystems(), g.DefaultNameSystem(), g.Packages)
}

func (PkgGenerators) NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer("", nil),
	}
}

func (PkgGenerators) DefaultNameSystem() string {
	return "public"
}

func (g PkgGenerators) Packages(context *generator.Context, generatorArgs *args.GeneratorArgs) generator.Packages {
	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}

	for i := range inputs {
		pkg := context.Universe[i]

		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					for i := range g {
						generators = append(generators, g[i](pkg, generatorArgs))
					}
					return
				},
			})
	}

	return packages
}
