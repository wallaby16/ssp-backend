package billing

import (
	"time"

	"log"

	"github.com/coreos/etcd/client"
)

var cfg client.Config
var Api client.KeysAPI

func init() {
	cfg = client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,

		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	Api = client.NewKeysAPI(c)
}

