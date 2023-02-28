# Demo data seed

Creates dummy lists and feeds random data into them.

## Prerequisites

- SharePoint Online tenant
- SharePoint Online Addn-Only auth credentials
- Go 1.19+

## Configure connection

- Create `./config/private.json` file
- Populate with credentials and strategy name, e.g.:

```json
{
  "strategy": "addin",
  "siteUrl": "https://contoso.sharepoint.com/sites/cloudquery",
  "clientId": "e1990a0a-dcf7-4b71-8b96-2a53c7e323e0",
  "clientSecret": "1wlWvB7V...zG1AqSP8="
}
```

## Run provisioning

```bash
go run ./cmd/seed/...
```

It sould create lists and feed bunch of random data.

The process takes time due to the number of seeding items. Amend `./cmd/demo/main.go` to reduce or increase the number of items.

If you see progress drammatically slowing dow or stuck, reduce concurrency as you face SharePoint API throttling.
