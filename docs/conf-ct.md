# Configuration: Content Types based rollup

Content Types based rollup allows to fetch data from multiple lists or document libraries based on the Content Type configuration.

All items based on the parent content type are fetched from all lists and subwebs below the context site URL.

```yaml
# sharepoint.yml
# ...
spec:
  # A map of Content Types with the rollup configurations
  content_types:
    # Base Content Type name or ID (e.g. "0x0101009D1CB255D" must be in quotes)
    Task:
      # REST `$select` OData modificator, fields entity properties array
      select:
        - Title
        - AssignedTo/Title
      # REST `$expand` OData modificator, fields entity properties array
      expand:
        - AssignedTo
      # Optional, an alias for the table name
      # the name of the alias is prefixed with `rollup_`
      alias: "task"
```
