package main

import (
	"github.com/cloudquery/plugin-sdk/v2/serve"
	"github.com/koltyakov/cq-source-sharepoint/plugin"
)

func main() {
	serve.Source(plugin.Plugin())
}
