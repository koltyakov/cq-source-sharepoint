A single source `yml` configuration assumes fetching data from a single SharePoint site. If you need to fetch data from multiple sites, you can create multiple source configurations. Alternatevely, search based data fetching can be used for rollup scenarios grabbing data from as many sites as needed.

```yaml
# sharepoint.yml
# ...
spec:
  auth:
    # ...
  lists:
    # ...
  content_types:
    # ...
  mmd:
    # ...
  search:
    # ...
  profiles:
    # ...
```

### Config: Authentication

:conf_auth

### Config: Lists

:conf_lists

### Config: Document libraries

:conf_libs

### Config: User Information List

:conf_uil

### Config: Content Types

:conf_cts

### Config: Managed Metadata

:conf_mmd

### Config: Search

:conf_search

### Config: User Profiles

:conf_ups

### Interactive Schema Builder

:conf_spctl
