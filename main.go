package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
)

func main() {
	var (
		token, owner, repo, tag string
	)
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&owner, "owner", "", "")
	flag.StringVar(&repo, "repo", "", "")
	flag.StringVar(&tag, "tag", "", "")
	flag.Parse()

	ctx := context.Background()
	gh := github.NewTokenClient(ctx, token)
	opts := &github.PackageListOptions{
		PackageType: github.String("container"),
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	pkg, _, err := gh.Organizations.GetPackage(ctx, owner, "container", repo)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "get package: %s\n", err)
		os.Exit(1)
	}

	for {
		versions, resp, err := gh.Organizations.PackageGetAllVersions(ctx, owner, "container", pkg.GetName(), opts)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "list versions: %s\n", err)
			os.Exit(1)
		}

		for _, vsn := range versions {
			tags := vsn.Metadata.Container.Tags
			for _, t := range tags {
				if t == tag {
					fmt.Println(*(vsn.HTMLURL) + "?tag=" + tag)
					os.Exit(0)
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	_, _ = fmt.Fprintf(os.Stderr, "missing tag: %s\n", tag)
	os.Exit(1)
}
