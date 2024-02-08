package cli

import "github.com/urfave/cli/v2"

type Command[T any] interface {
	action(cCtx *cli.Context) error
	extract(cCtx *cli.Context) (T, error)
	flags() []cli.Flag
	load(cCtx *cli.Context, data T)
}
