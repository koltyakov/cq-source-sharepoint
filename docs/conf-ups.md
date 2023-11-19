# Configuration: User Profiles

User Profiles are fetched via Search API, so the search should be configured in the farm.

Search drived data source can be user only with user associated authentication strategies. E.g. it won't work with `addin` strategy.

```yaml
# sharepoint.yml
# ...
spec:
  # Include `profiles` property to fetch user profiles
  # Object structure for extensibility (adding custom properties)
  profiles:
    enabled: true
    # Optional, an alias for the table name
    alias: "profile"
```
