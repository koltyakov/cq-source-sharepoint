# Plugin versioning

Versioning is managed via git tags. Then a new version tag is pushed to the repository, GoReleases pipeline automatically builds and publishes plugin to GitHub packages.

## Set version

```bash
git tag v1.7.2
```

### Push tag

```bash
git push --tags
```

### Delete tag

```bash
git tag -d v1.7.2-test
git push --delete origin v1.7.2-test
```
