package command

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/joelferrier/marsho/repository"
	"strings"
)

type FindCommand struct {
	Meta
	HelpText string
}

type findOpts struct {
	RepoUrl  string
	NoVerify bool
	KernVer  string
}

func (c *FindCommand) setHelp() {
	c.HelpText = `
Usage: marsho find [options] [kernel-version]
    Search repository for LiME kernel modules

    [options]
    -repo string      repository url
                      Default: https://threatresponse-lime-modules.s3.amazonaws.com/
    -gpg-no-verify    disable GPG Verification

    [kernel-version]  kernel module version eg. 4.4.10-22.54.amzn1.x86_64
                      Globs are supported eg. 4.4.10*amzn1.x86_64
`
}

func (c *FindCommand) Run(args []string) int {
	opts, err := findArgs(args)
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

	modules, err := repo.Find(opts.KernVer)
	if err != nil {
		log.Critical(err)
		return 0
	}

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	for _, mod := range modules {
		table.AddRow(fmt.Sprintf("kernel: %s", mod.Version), fmt.Sprintf("path: /modules/%s", mod.Name))
	}

	fmt.Println(table)
	fmt.Printf("\nMatched %d LiME modules for '%s' in %s\n", len(modules), opts.KernVer, repo.BaseUrl)
	return 1
}

func (c *FindCommand) Help() string {
	c.setHelp()
	return strings.TrimSpace(c.HelpText)
}

func (c *FindCommand) Synopsis() string {
	return "Search repository for LiME kernel modules"
}

func findArgs(args []string) (findOpts, error) {
	opts := findOpts{}

	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	repoUrl := findCmd.String("repo", "", "LiME Repository url")
	noVerify := flag.Bool("gpg-no-verify", false, "Disable GPG Verification")

	findCmd.Parse(args)
	log.Debug(fmt.Sprintf("parsed repoUrl: %s", *repoUrl))
	log.Debug(fmt.Sprintf("parsed noVerify: %t", *noVerify))

	var kernVer string
	if len(findCmd.Args()) != 1 {
		return opts, errors.New("find: missing kernel-version argument")
	} else {
		kernVer = findCmd.Args()[0]
	}

	opts.RepoUrl = *repoUrl
	opts.NoVerify = *noVerify
	opts.KernVer = kernVer

	return opts, nil
}
