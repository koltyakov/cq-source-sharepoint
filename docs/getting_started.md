---
hub-title: Getting Started
---

### Install CloudQuery

Follow [quickstart instructions](https://www.cloudquery.io/docs/quickstart/).

### Source sample data

Provision and seed some sample data. [See more](https://github.com/koltyakov/cq-source-sharepoint/blob/main/cmd/seed/README.md). Which satisfy the schema below.

### Setup authentication

```bash
# .env or env vars export
# See more details in https://go.spflow.com/auth/strategies
SP_SITE_URL=https://contoso.sharepoint.com/sites/site
```

or use "ondeman" auth.

### Source configuration

```yaml
# sharepoint.yml
kind: source
spec:
  name: sharepoint
  registry: cloudquery
  path: "koltyakov/sharepoint"
  version: "VERSION_SOURCE_SHAREPOINT"
  destinations: ["sqlite"]
  tables: ["*"]
  spec:
    auth:
      strategy: "ondemand"
      creds:
        siteUrl: ${SP_SITE_URL}
        # align creds with the used strategy
    lists:
      _catalogs/users:
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

### Destination configuration

For the sake of simplicity, we'll use SQLite destination.

```yaml
# sqlite.yml
kind: destination
spec:
  name: sqlite
  path: cloudquery/sqlite
  version: "VERSION_DESTINATION_SQLITE"
  spec:
    connection_string: ./sp.db
```

### Run CloudQuery

```bash
# With auth environment variables exported
cloudquery sync sharepoint.yml sqlite.yml
```

You should see the following output:

```bash
Loading spec(s) from sharepoint_reg.yml, sqlite.yml
Downloading https://github.com/koltyakov/...sharepoint_darwin_arm64.zip
Downloading 100% |█████████████████████████████████████| (5.2/5.2 MB, 10 MB/s)
Migration completed successfully.
Starting sync for: sharepoint (v2.1.0) -> [sqlite (v2.4.15)]
Sync completed successfully. Resources: 37478, Errors: 0, Panics: 0, Time: 21s
```

Check for destination database data.
