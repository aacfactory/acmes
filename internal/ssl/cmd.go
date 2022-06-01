package ssl

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

var Command = &cli.Command{
	Name:        "ca",
	Usage:       "ca -c {common name} -e {expire days} -o {out dir}",
	Description: "generate self signed ca",
	ArgsUsage:   "",
	Category:    "",
	Action: func(c *cli.Context) error {
		cn := strings.TrimSpace(c.String("cn"))
		out := strings.TrimSpace(c.String("out"))
		expires := c.Int("expires")
		err := generate(cn, expires, out)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("acmes: ca was generated succeed, see %s", out))
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Required: true,
			Name:     "cn",
			Value:    "",
			Usage:    "common name for ca",
			Aliases:  []string{"c"},
		},
		&cli.IntFlag{
			Required: true,
			Name:     "expires",
			Value:    0,
			Usage:    "expire days for ca",
			Aliases:  []string{"e"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "out",
			Value:    "",
			Usage:    "out dir for ca",
			Aliases:  []string{"o"},
		},
	},
	HelpName:           "",
	CustomHelpTemplate: "",
}
