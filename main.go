package main

import (
	"github.com/cloudquery/plugin-sdk/serve"
	"github.com/koltyakov/cq-source-sharepoint/plugin"
)

func main() {
	serve.Source(plugin.Plugin())
}
