# lookup-package-url
GitHub Action that provides a container package's web page URL as an output to other actions.

For example, to directly link to the package page from a GitHub Deployment.

```yaml
      - id: lookup-url
        name: Lookup Container Package URL
        uses: eigenbot-app/lookup-package-url@main
        with:
          owner: my-github-org
          repo: my-github-repo
          tag: the container's tag from your docker build step, e.g., ${{ steps.meta.outputs.tags }}
          token: token to access the GitHub API, e.g., ${{ secrets.GITHUB_TOKEN }}
```

Access the output with: `${{ steps.lookup-url.outputs.url }}`.
