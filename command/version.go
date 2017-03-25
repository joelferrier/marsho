package command

import "fmt"

type VersionCommand struct {
	Meta

	BuildTime string
	Revision  string
	Version   string
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(args []string) int {
	fmt.Printf("marsho v%s\n", c.Version)
	fmt.Printf("git commit hash: %s\n", c.Revision)
	fmt.Printf("build time: %s\n", c.BuildTime)
	return 1
}

func (c *VersionCommand) Synopsis() string {
	return "Print marsho version"
}
