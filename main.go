package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v56/github"
)

func main() {
	outp := os.Getenv("GITHUB_OUTPUT")
	if outp == "" {
		_, _ = fmt.Fprintf(os.Stderr, "unset: $GITHUB_OUTPUT\n")
		os.Exit(1)
	}

	var (
		token, owner, repo, tag string
	)
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&owner, "owner", "", "")
	flag.StringVar(&repo, "repo", "", "")
	flag.StringVar(&tag, "tag", "", "")
	flag.Parse()

	if _, after, found := strings.Cut(tag, ":"); found {
		tag = after
	}

	ctx := context.Background()
	gh := github.NewClient(nil).WithAuthToken(token)
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
					f, err := os.OpenFile(outp, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
						os.Exit(1)
					}
					if _, err := f.WriteString(
						fmt.Sprintf("url=%s?tag=%s\n", *(vsn.HTMLURL), tag),
					); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
						os.Exit(1)
					}
					_ = f.Close()
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
