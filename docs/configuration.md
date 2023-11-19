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

```yaml
# sharepoint.yml
# ...
spec:
  auth:
    strategy: "azurecert"
    creds:
      siteUrl: "https://contoso.sharepoint.com/sites/cloudquery"
      tenantId: "e1990a0a-dcf7-4b71-8b96-2a53c7e323e0"
      clientId: "2a53c7e323e0-e1990a0a-dcf7-4b71-8b96"
      certPath: "/path/to/cert.pfx"
      certPass: "certpass"
```

`creds` options are unique for different auth strategies. See more details in [Auth strategies](https://go.spflow.com/auth/strategies).

We recomment Azure AD (`azurecert`) or Add-In (`addin`) auth for production scenarios for SharePoint Online. Yet, other auth strategies are also available, e.g. `saml`, `device`. Some of the APIs could require using user contextual auth, for instance, Search API can't work without a user context.

SharePoint On-Premise auth is also supported, based on your farm configuration you can use: `ntlm`, `adfs` to name a few.

Need to hands on quickly without configuring Azure Apps or Addins or asking an admin to turn on app password? Try On-Demand auth:

```yaml
# sharepoint.yml
# ...
spec:
  auth:
    strategy: "ondemand"
    creds:
      siteUrl: "https://contoso.sharepoint.com/sites/cloudquery"
```

### Config: Lists

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

### Config: Document libraries

Document listariries are the same as lists in SharePoint, but with a few differences. And it's common to expand File entity to get file metadata.

Also, a document library URI usually doesn't contain `Lists/` prefix.

```yaml
# sharepoint.yml
# ...
spec:
  lists:
    Shared Documents:
      select:
        - FileLeafRef
        - FileRef
        - FileDirRef
        - File/Length
      expand:
        - File
      alias: "document"
```

### Config: User Information List

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

### Config: Content Types

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

### Config: Managed Metadata

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

### Config: Search

Search-drived datasource can be user only with user associated authentication strategies. E.g. it won't work with `addin` strategy.

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

### Config: User Profiles

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

### Interactive Schema Builder

The plugin ships with configuration utility `spctl`.

![](https://github.com/koltyakov/cq-source-sharepoint/blob/main/assets/spctl.gif?raw=true)

It can be downloaded from [releases](https://github.com/koltyakov/cq-source-sharepoint/releases): `spctl_[OS]_[ARCH].zip`.

On a macOS System Settings / Security allowance is needed for it to run.
