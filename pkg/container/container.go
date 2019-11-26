package container

import (
	"io"
	"os"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/sirupsen/logrus"

	"github.com/spf13/afero"
)

// Container holds references to external systems that need to be captured or modified in tests
type Container struct {
	DebugLevel        value.DebugLevelValue
	DebugReportCaller bool
	FileSystem        *afero.Afero
	In                io.Reader
	Out               io.Writer
	Err               io.Writer
	Log               *logrus.Logger
}

// New instantiates the Container struct above
func New() *Container {
	return &Container{
		In:         os.Stdin,
		Out:        os.Stdout,
		Err:        os.Stderr,
		FileSystem: &afero.Afero{Fs: afero.NewOsFs()},
		Log:        logrus.New(),
	}
}
