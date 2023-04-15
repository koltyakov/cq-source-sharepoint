# Demo data seed

Creates dummy lists and feeds random data into them.

## Prerequisites

- SharePoint Online tenant
- SharePoint Online Addn-Only auth credentials
- Go 1.19+

## Run provisioning

On a Mac/Linux machine, run the following command:

```bash
SP_SITE_URL="https://contoso.sharepoint.com/sites/site" go run ./cmd/seed/...
```

On a Windows machine, run the following command:

```powershell
$env:SP_SITE_URL="https://contoso.sharepoint.com/sites/site"
go run ./cmd/seed/...
```

> If you are using On-Premise SharePoint with NTLM authentication, modify `./cmd/demo/main.go` to use `ntlm` auth provider.

It should create lists and feed bunch of random data.

The process takes time due to the number of seeding items. Amend `./cmd/demo/main.go` to reduce or increase the number of items.

If you see progress drammatically slowing dow or stuck, reduce concurrency as you face SharePoint API throttling.
