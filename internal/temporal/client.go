package temporal

import (
	"sync"

	"github.com/shayansm2/askhn/internal/config"
	"go.temporal.io/sdk/client"
)

var once sync.Once
var c client.Client

func GetClient() client.Client {
	once.Do(func() {
		var err error
		c, err = client.Dial(client.Options{
			HostPort: config.Load().TemporalHost,
		})
		if err != nil {
			panic("Unable to create Temporal client: " + err.Error())
		}
	})
	return c
}
