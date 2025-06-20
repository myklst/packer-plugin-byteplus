package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/volcengine/packer-plugin-volcengine/builder/ecs"
	datasource "github.com/volcengine/packer-plugin-volcengine/datasource/image"
	"github.com/volcengine/packer-plugin-volcengine/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("images", new(datasource.Datasource))
	pps.RegisterBuilder("ecs", new(ecs.Builder))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
