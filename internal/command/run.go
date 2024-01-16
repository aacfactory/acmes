package command

import (
	"github.com/aacfactory/acmes/internal/server"
	"github.com/aacfactory/acmes/internal/ssl"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	copyright = `Copyright 2021 Wang Min Xiang

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

	usage = `see COMMANDS`

	description = `get more providers at https://go-acme.github.io/lego/dns/`
)

func Run() error {
	app := &cli.App{
		Name:        "acmes",
		HelpName:    "help",
		Usage:       "acmes tool",
		UsageText:   usage,
		ArgsUsage:   "",
		Version:     "v1.1.0",
		Description: description,
		Commands: []*cli.Command{
			ssl.Command,
			server.Command,
		},
		Authors: []*cli.Author{
			{
				Name:  "Wang Min Xiang",
				Email: "wangminxiang@aacfactory.co",
			},
		},
		Copyright: copyright,
	}
	return app.Run(os.Args)
}
