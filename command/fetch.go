package command

import (
	"errors"
	"flag"
	"fmt"
	"github.com/joelferrier/marsho/repository"
	"strings"
)

type FetchCommand struct {
	Meta
	HelpText string
}

type fetchOpts struct {
	RepoUrl  string
	NoVerify bool
	KernVer  string
}

func (c *FetchCommand) setHelp() {
	c.HelpText = `
Usage: marsho fetch [options] [kernel-version]
    Fetch LiME kernel module

    [options]
    -repo string   repository url
                   Default: https://threatresponse-lime-modules.s3.amazonaws.com/
    -gpg-no-verify disable GPG Verification

    [kernel-version]
    kernel module version, eg. 4.4.10-22.54.amzn1.x86_64
`
}

func (c *FetchCommand) Run(args []string) int {
	opts, err := fetchArgs(args)
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
	repo.Get(opts.KernVer)

	return 1
}

func (c *FetchCommand) Help() string {
	c.setHelp()
	return strings.TrimSpace(c.HelpText)
}

func (c *FetchCommand) Synopsis() string {
	return "Fetch LiME kernel module"
}

func fetchArgs(args []string) (fetchOpts, error) {
	opts := fetchOpts{}

	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	repoUrl := fetchCmd.String("repo", "", "LiME Repository url")
	noVerify := flag.Bool("gpg-no-verify", false, "Disable GPG Verification")

	fetchCmd.Parse(args)
	log.Debug(fmt.Sprintf("parsed repoUrl: %s", *repoUrl))
	log.Debug(fmt.Sprintf("parsed noVerify: %t", *noVerify))

	var kernVer string
	if len(fetchCmd.Args()) != 1 {
		return opts, errors.New("fetch: missing kernel-version argument")
	} else {
		kernVer = fetchCmd.Args()[0]
	}
	log.Debug(fmt.Sprintf("parsed kernVer: %s", kernVer))

	opts.RepoUrl = *repoUrl
	opts.NoVerify = *noVerify
	opts.KernVer = kernVer

	return opts, nil
}
