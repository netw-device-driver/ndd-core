package clicmd

import (
	"github.com/netw-device-driver/ndd-runtime/pkg/parser"
	"github.com/spf13/afero"
)

type BuildChild struct {
	name   string
	linter parser.Linter
	fs     afero.Fs
}

type PushChild struct {
	tag string
	fs  afero.Fs
}

// pushProviderCmd pushes a Provider.
type PushProviderCmd struct {
	Tag string `arg:"" help:"Tag of the package to be pushed. Must be a valid OCI image tag."`
}
