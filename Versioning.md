# Plugin versioning

Versioning is managed via git tags. Then a new version tag is pushed to the repository, GoReleases pipeline automatically builds and publishes plugin to GitHub packages.

## Set version

```bash
git tag v2.0.0.
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
