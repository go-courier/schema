module github.com/go-courier/schema

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-courier/codegen v1.1.2
	github.com/go-courier/ptr v1.0.1
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	k8s.io/gengo v0.0.0-20210203185629-de9496dff47b
	k8s.io/kube-openapi v0.0.0-20210524163139-412c2b45c7d3
)

// https://github.com/kubernetes/gengo/pull/203
replace k8s.io/gengo => github.com/morlay/gengo v0.0.0-20210527040048-db12fc2590e7
