package cli

import "github.com/leaanthony/clir"

type Cli struct {
	*clir.Cli
}

func (cli *Cli) Init() {
	cli.RegisterSearchAction()
	cli.RegisterGetAction()
}
