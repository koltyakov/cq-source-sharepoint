# Configuration: Lists

SharePoint lists is the main artifact for customizable data storage in SharePoint.

Lists fetching configuration follows same naming as SharePoint REST API.

```yaml
# sharepoint.yml
# ...
spec:
  # A map of URIs with the list configurations
  # If no lists are provided, nothing will be fetched
  lists:
    # List or Document library URI - a relative path without a site URL
    # Can be checker in the browser URL (exclude site URL and view page path)
    Lists/ListEntityName:
      # REST `$select` OData modificator, fields entity properties array
      # Wildcard selectors `*` are intentionally not supported
      # If not provided, only default fields will be fetched (ID, Created, AuthorId, Modified, EditorId)
      select:
        - Title
        - Author/Title
        # Fields mapping via `->` arrow alias, when a specific field name is considered
        - EditorId -> editor
      # REST `$expand` OData modificator, fields entity properties array
      # When expanding an entity use selection of a nested entity property(s)
      # Optional, and in most of the cases we recommend to avoid it and
      # prefer to map nested entities to the separate tables
      expand:
        - Author
      # REST `$filter` OData modificator, a filter string
      # Don't use filters for large entities which potentially
      # can return more than 5000 in a view
      # such filtering will throttle no matter top limit is set
      filter: "Active eq true"
      # Optional, an alias for the table name
      # Don't map different lists to the same table - such scenario is not supported
      alias: "my_table"
    Lists/AnotherList:
      select:
        - Title
```
