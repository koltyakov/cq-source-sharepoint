# Configuration: Search queries

Search drived data source can be user only with user associated authentication strategies. E.g. it won't work with `addin` strategy.

```yaml
# sharepoint.yml
# ...
spec:
  # A map of search queries
  search:
    # Query name (whatever you want to name a resulted table)
    # Should be unique within other compound aliases
    documents:
      # Required, search query text
      # https://learn.microsoft.com/en-us/sharepoint/dev/general-development/sharepoint-search-rest-api-overview#querytext-parameter
      query_text: "*"
      # Optional, the managed properties to return in the search results
      # https://learn.microsoft.com/en-us/sharepoint/dev/general-development/sharepoint-search-rest-api-overview#selectproperties
      # By defining the list of properties, you also tell the plugin
      # to have correcponding columns in the table
      select_properties:
        - Size
        - Title
        - ContentTypeId
        - IsDocument
        - FileType
        - DocId
        - SPWebUrl
        - SiteId
        - WebId
        - ListId
      # Optional, whether duplicate items are removed from the results
      # https://learn.microsoft.com/en-us/sharepoint/dev/general-development/sharepoint-search-rest-api-overview#trimduplicates
      trim_duplicates: true
    profiles:
      query_text: "*",
      trim_duplicates: false
      # The result source ID to use for executing the search query.
      # https://learn.microsoft.com/en-us/sharepoint/dev/general-development/sharepoint-search-rest-api-overview#sourceid
      source_id: "b09a7990-05ea-4af9-81ef-edfab16c4e31"
```
