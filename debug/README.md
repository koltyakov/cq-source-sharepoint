# Local functional debug

## Datasource

Copy `.env.example` to `.env` and fill in the values.

The example values stands for SharePoint Online Addn-Only auth scenario. See more details in [SharePoint Online Addn-Only auth](https://go.spflow.com/auth/strategies/addin).

## Run sync

### SQLite destination

In current directory run:

```bash
make sync-sqlite
```

Check console output for errors. And data in the local database.

### Postgresql destination

- make sure `destination` is aligned in `sharepoint.yml`
- run `docker-compose up -d` to start postgresql

```bash
make sync-postgresql
```
