package filme

import "github.com/florinutz/filme/pkg/config"

// Filme is the logic behind the cli app.
type Filme struct {
	*config.Config
}

func New() *Filme {
	return &Filme{
		Config: config.New(),
	}
}
