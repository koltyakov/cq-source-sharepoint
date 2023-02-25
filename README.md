# cq-source-sharepoint

[CloudQuery](https://github.com/cloudquery/cloudquery) SharePoint Source community plugin.

## Schema

```yaml
kind: source
spec:
  name: "sharepoint"
  registry: "github"
  path: "koltyakov/sharepoint"
  version: "v1.0.0" # provide the latest stable version
  destinations: ["postgresql"] # provide the list of used destinations
  spec:
    # Spec is mandatory
    # This plugin follows idealogy of explicit configuration
    # we can change this in future based on the feedback
```

### Authentication options

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

We always recomment Azure AD (`azurecert`) or Add-In (`addin`) auth for production scenarios for SharePoint Online. Yet, other auth strategies are available for testing and development purposes, e.g. `saml`, `device`.

SharePoint On-Premise auth is also supported, based on your farm configuration you can use: `ntlm`, `adfs` to name a few.

### Entities configuration

So far, the plugin supports lists and document libraries data fetching. Base on feedback and use cases, we have a strategy for extending the plugin to support other SharePoint API entities, e.g. Managed Metadata, UPS, etc.

A single source `yml` configuration assumes fetching data from a single SharePoint site. If you need to fetch data from multiple sites, you can create multiple source configurations.


```yaml
# sharepoint.yml
# ...
  spec:
    # A map of URIs to the list configuration
    # If no lists are provided, nothing will be fetched
    lists:
      # List or Document library URI - a relative path without a site URL
      # Can be checker in the browser URL (exclude site URL and view page path)
      Lists/ListEntityName:
        # REST's `$select` OData modificator, fields entity properties array
        # Wildcard selectors `*` are intentionally not supported
        # If not provided, only default fields will be fetched (ID, Created, AuthorId, Modified, EditorId)
        select:
          - Title
          - Author/Title
        # REST's `$expand` OData modificator, fields entity properties array
        # When expanding an entity use selection of a nested entity property(s)
        # Optional, and in most of the cases we recommend to avoid it and 
        # prefer to map nested entities to the separate tables
        expand:
          - Author
        # Optional, an alias for the table name
        # Don't map different lists to the same table - such scenariou is not supported
        alias: "my_table"
      Lists/AnotherList:
        select:
          - Title
```

#### User Information List

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
          - IsSiteAdmin
          - Deleted
        alias: "user"
```

#### Document libraries

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

## Schema e2e sample

```bash
# .env or env vars export
# See more details in https://go.spflow.com/auth/strategies
SP_AUTH_STRATEGY=addin
SP_SITE_URL=https://contoso.sharepoint.com/sites/site
SP_CLIENT_ID=97e6ed51-777c-42da-8f07-b035a5ac057b
SP_CLIENT_SECRET="1wlWvB...AqSP8="
```

```yaml
# sharepoint.yml
kind: source
spec:
  name: "sharepoint"
  registry: "github"
  path: "koltyakov/sharepoint"
  version: "v1.0.0"
  destinations: ["postgresql"]
  spec:
    auth:
      strategy: "${SP_AUTH_STRATEGY}"
      creds:
        siteUrl: ${SP_SITE_URL}
        clientId: ${SP_CLIENT_ID}
        clientSecret: ${SP_CLIENT_SECRET}
    lists:
      _catalogs/users:
        select:
          - Title
        alias: "user"
      Shared Documents:
        select:
          - FileLeafRef
          - FileRef
          - FileDirRef
          - Author/Title
          - File/Length
        expand:
          - Author
          - File
        alias: "document"
      Lists/Managers:
        select:
          - Title
        alias: "manager"
      Lists/Customers:
        select:
          - Title
          - RoutingNumber
          - Region
          - Revenue
          - ManagerId
        alias: "customer"
      Lists/Orders:
        select:
          - Title
          - CustomerId
          - OrderNumber
          - OrderDate
          - Total
        alias: "order"
```

Powered by [Gosip](https://github.com/koltyakov/gosip).