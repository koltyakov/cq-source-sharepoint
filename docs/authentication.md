# SharePoint Source Plugin Authentication Options

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
