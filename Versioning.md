# Plugin versioning

Versioning is managed via git tags. Then a new version tag is pushed to the repository, GoReleases pipeline automatically builds and publishes plugin to GitHub packages.

## Set version

```bash
git tag v2.0.0
```

### Push tag

```bash
git push --tags
```

### Delete tag

```bash
git tag -d v2.0.0-test
git push --delete origin v2.0.0-test
```

## Publish

Make sure you logged in to CloudQuery Hub and switched to the right team:

```bash
cloudquery login
cloudquery switch koltyakov
```

Package the plugin:

```bash
go run main.go package --docs-dir docs -m @CHANGELOG.md $(git describe --tags --abbrev=0) .
```

Publish the plugin (draft):

```bash
cloudquery plugin publish
```

Publish the plugin (release):

```bash
cloudquery plugin publish --finalize
```
