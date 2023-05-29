package main

import (
	"github.com/cloudquery/plugin-sdk/v3/serve"
	"github.com/koltyakov/cq-source-sharepoint/plugin"
)

func main() {
	serve.Source(plugin.Plugin())
}
