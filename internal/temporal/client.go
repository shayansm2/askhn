package temporal

import (
	"sync"

	"go.temporal.io/sdk/client"
)

var once sync.Once
var c client.Client

func GetClient() client.Client {
	once.Do(func() {
		var err error
		c, err = client.Dial(client.Options{})
		if err != nil {
			panic("Unable to create Temporal client: " + err.Error())
		}
	})
	return c
}
