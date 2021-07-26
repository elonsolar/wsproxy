package example

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/elonsolar/wsproxy/client"
	"testing"
)

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "cfg", "./client.toml", "")
}

// go test -args -cfg ../client.toml
func TestProxy(t *testing.T) {
	flag.Parse()

	var pxyCfg client.Config
	if _, err := toml.DecodeFile(cfg, &pxyCfg); err != nil {
		panic(err)
	}

	client.Proxy(&pxyCfg)
	select {}
}
