package core

import (
	"github.com/syncfuture/go/sconfig"
	"github.com/syncfuture/host"
	"github.com/syncfuture/host/sfasthttp"
)

var (
	Config sconfig.IConfigProvider
	Host   host.IWebHost
)

func init() {
	Config = sconfig.NewJsonConfigProvider()
	Host = sfasthttp.NewFHWebHost(Config)
}
