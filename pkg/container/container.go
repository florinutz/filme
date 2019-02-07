package container

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

// Container holds references to external systems that need to be captured or modified in tests
type Container struct {
	FileSystem *afero.Afero
	Out        io.Writer
	Err        io.Writer
}

// New instantiates the Container struct above
func New() *Container {
	return &Container{
		FileSystem: &afero.Afero{Fs: afero.NewOsFs()},
		Out:        os.Stdout,
		Err:        os.Stderr,
	}
}
