package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg"
)

func (f *Filme) PrintVersion() {
	fmt.Fprintf(f.Out, "filme %s (%s) (%s)\n", pkg.Version, pkg.Commit, pkg.BuildTime)
}
