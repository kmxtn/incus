package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	cli "github.com/lxc/incus/v6/internal/cmd"
	"github.com/lxc/incus/v6/internal/i18n"
	"github.com/lxc/incus/v6/internal/version"
)

type cmdVersion struct {
	global *cmdGlobal
}

// Command returns a cobra.Command for use with (*cobra.Command).AddCommand.
func (c *cmdVersion) Command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = usage("version", i18n.G("[<remote>:]"))
	cmd.Short = i18n.G("Show local and remote versions")
	cmd.Long = cli.FormatSection(i18n.G("Description"), i18n.G(
		`Show local and remote versions`))

	cmd.RunE = c.Run

	return cmd
}

// Run runs the actual command logic.
func (c *cmdVersion) Run(cmd *cobra.Command, args []string) error {
	// Quick checks.
	exit, err := c.global.checkArgs(cmd, args, 0, 1)
	if exit {
		return err
	}

	fmt.Printf(i18n.G("Client version: %s\n"), version.Version)

	// Remote version
	remote := ""
	if len(args) == 1 {
		remote = args[0]
		if !strings.HasSuffix(remote, ":") {
			remote = remote + ":"
		}
	}

	ver := i18n.G("unreachable")
	resources, err := c.global.parseServers(remote)
	if err == nil {
		resource := resources[0]
		info, _, err := resource.server.GetServer()
		if err == nil {
			ver = info.Environment.ServerVersion
		}
	}

	fmt.Printf(i18n.G("Server version: %s\n"), ver)

	return nil
}
