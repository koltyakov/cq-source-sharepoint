# Debug

## Debug gRPC server

### Run gRPC server (CLI)

```bash
go run main.go serve
```

## Run gRPC server (VSCode)

As an alternative to the previous step, you can run the gRPC server in debug mode from VSCode. Use `Debug Launch` configuration. You can use breakpoints and step through the code.

### Run sync command

Make sure source configuration in `debug/debug.yml` corresponds your SharePoint environment.

Run the sync:

```bash
cloudquery sync debug/debug.yml
```

## Debug with plugin build

See more: [debug/README.md](debug/README.md)
