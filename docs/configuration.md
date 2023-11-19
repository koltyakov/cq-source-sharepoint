# SharePoint Source Plugin Configuration Reference

A single source `yml` configuration assumes fetching data from a single SharePoint site. If you need to fetch data from multiple sites, you can create multiple source configurations. Alternatevely, search based data fetching can be used for rollup scenarios grabbing data from as many sites as needed.

```yaml
# sharepoint.yml
# ...
spec:
  # auth:
  # lists:
  # content_types:
  # mmd:
  # search:
  # profiles:
```

See more details about configuration options in the following sections:

- [auth](authentication)
- [lists](config-lists)
- [content_types](conf-ct)
- [mmd](conf-mmd)
- [search](conf-search)
- [profiles](conf-ups)
