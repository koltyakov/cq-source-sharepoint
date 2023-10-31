package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

var (
	Name    = "sharepoint"
	Kind    = "source"
	Team    = "koltyakov"
	Version = "development"
)

func NewPlugin() *plugin.Plugin {
	return plugin.NewPlugin(Name, Version, NewClient, plugin.WithKind(Kind), plugin.WithTeam(Team))
}
