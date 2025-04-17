//go:build ignore
// +build ignore

// This file defines the generate directive for producing mock implementations of the
// domain port interfaces. Run `go generate ./internal/domain/port` to regenerate mocks.
package port

//go:generate mockery --all --dir . --output ../../internal/mocks/domain/port --outpkg mocks