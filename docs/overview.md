---
name: SharePoint
stage: GA
title: SharePoint Source Plugin
description: CloudQuery SharePoint Source Plugin documentation
---

# SharePoint Source Plugin

The SharePoint Source plugin allows you to fetch data from SharePoint and load it into any supported CloudQuery destination (e.g. PostgreSQL, BigQuery, Snowflake, and [more](https://hub.cloudquery.io/plugins/destination)).

## Features

- Lists and Document Libraries data fetching
- Content Types based rollup
- User Information List data fetching
- Search Query data source
- User Profiles data source
- Managed Metadata data source
- SharePoint Online support
- SharePoint On-Premise support

![demo](https://github.com/koltyakov/cq-source-sharepoint/blob/main/assets/demo.gif?raw=true)

## Schema

```yaml
kind: source
spec:
  name: sharepoint
  registry: cloudquery
  path: koltyakov/sharepoint
  version: "VERSION_SOURCE_SHAREPOINT"
  destinations: ["postgresql"]
  tables: ["*"]
  spec:
    # Spec is mandatory
    # This plugin follows idealogy of explicit configuration
```
