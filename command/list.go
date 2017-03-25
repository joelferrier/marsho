package command

import (
	"flag"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/joelferrier/marsho/repository"
	"strings"
)

type ListCommand struct {
	Meta
	HelpText string
}

type listOpts struct {
	RepoUrl  string
	NoVerify bool
}

func (c *ListCommand) setHelp() {
	c.HelpText = `
Usage: marsho list [options]
    List availible LiME kernel modules

    [options]
    -repo string   repository url
                   Default: https://threatresponse-lime-modules.s3.amazonaws.com/
    -gpg-no-verify disable GPG Verification
`
}

func (c *ListCommand) Run(args []string) int {
	opts, err := listArgs(args)
	if err != nil {
		fmt.Printf("%s\n\n%s\n", err, c.Help())
		return 0
	}
	repo := repository.DefaultRepository()
	if opts.RepoUrl != "" {
		// ensure repo url has a trailing slash
		if opts.RepoUrl[len(opts.RepoUrl)-1:] != "/" {
			opts.RepoUrl = opts.RepoUrl + "/"
		}
		repo.BaseUrl = opts.RepoUrl
	}

	//manifest, err := repository.List(conf)
	manifest, err := repo.List()
	if err != nil {
		log.Critical(err)
		return 0
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	for _, mod := range manifest.Modules {
		table.AddRow(fmt.Sprintf("kernel: %s", mod.Version), fmt.Sprintf("path: /modules/%s", mod.Name))
	}

	fmt.Println(table)
	fmt.Printf("\nFound %d LiME modules in %s\n", len(manifest.Modules), repo.BaseUrl)
	return 1
}

func (c *ListCommand) Help() string {
	c.setHelp()
	return strings.TrimSpace(c.HelpText)
}

func (c *ListCommand) Synopsis() string {
	return "List availible LiME kernel modules"
}

func listArgs(args []string) (listOpts, error) {
	opts := listOpts{}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	repoUrl := listCmd.String("repo", "", "LiME Repository url")
	noVerify := flag.Bool("gpg-no-verify", false, "Disable GPG Verification")

	listCmd.Parse(args)
	log.Debug(fmt.Sprintf("parsed repoUrl: %s", *repoUrl))
	log.Debug(fmt.Sprintf("parsed noVerify: %t", *noVerify))

	opts.RepoUrl = *repoUrl
	opts.NoVerify = *noVerify

	return opts, nil
}
