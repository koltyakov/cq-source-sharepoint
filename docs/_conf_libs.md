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
