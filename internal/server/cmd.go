package server

import (
	"github.com/urfave/cli/v2"
	"strings"
)

var Command = &cli.Command{
	Name:        "serve",
	Usage:       "serve --port 443 --ca {ca_path} --cakey {ca_key_path} --level info --email {email} --store {file:///some_dir_path} --provider {provider}",
	Description: "run acmes http server",
	ArgsUsage:   "",
	Category:    "",
	Action: func(c *cli.Context) error {
		return serve(options{
			port:     c.Int("port"),
			ca:       strings.TrimSpace(c.String("ca")),
			key:      strings.TrimSpace(c.String("cakey")),
			level:    strings.TrimSpace(c.String("level")),
			store:    strings.TrimSpace(c.String("store")),
			email:    strings.TrimSpace(c.String("email")),
			provider: strings.TrimSpace(c.String("provider")),
		})
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Value:   80,
			Usage:   "port for http server",
			EnvVars: []string{"ACMES_PORT"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "ca",
			Value:    "",
			Usage:    "ca file for http server",
			EnvVars:  []string{"ACMES_CA"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "cakey",
			Value:    "",
			Usage:    "ca key file for http server",
			EnvVars:  []string{"ACMES_CAKEY"},
		},
		&cli.StringFlag{
			Name:    "level",
			Value:   "info",
			Usage:   "level for logger",
			EnvVars: []string{"ACMES_LOG_LEVEL"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "store",
			Value:    "",
			Usage:    "store for certs",
			EnvVars:  []string{"ACMES_STORE"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "email",
			Value:    "",
			Usage:    "user email for acme",
			EnvVars:  []string{"ACMES_EMAIL"},
		},
		&cli.StringFlag{
			Required: true,
			Name:     "provider",
			Value:    "",
			Usage:    "dns provider for acme",
			EnvVars:  []string{"ACMES_DNS_PROVIDER"},
		},
	},
}
