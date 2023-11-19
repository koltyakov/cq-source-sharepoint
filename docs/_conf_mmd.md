To configure managed metadata fetching, you need to provide a term set ID (GUID) and an optional alias for the table name.

```yaml
# sharepoint.yml
# ...
spec:
  # A map of MMD term sets IDs (GUIDs)
  mmd:
    # Term set ID
    8ed8c9ea-7052-4c1d-a4d7-b9c10bffea6f:
      # Optional, an alias for the table name
      # the name of the alias is prefixed with `mmd_`
      alias: "department"
```
