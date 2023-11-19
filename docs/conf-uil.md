# Configuration: User Information List

Quite often you'd need getting User Information List for Author and Editor fields joining. This is a special case, and we have a dedicated configuration for it.

```yaml
# sharepoint.yml
# ...
spec:
  lists:
    _catalogs/users: # UIL list URI, source of People Picker lookup
      select:
        - Title
        - FirstName
        - LastName
        - JobTitle
        - Department
        - EMail
        - Deleted
      filter: "UserName ne null"
      alias: "user"
```
