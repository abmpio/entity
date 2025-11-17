package mongodb

import (
	"fmt"

	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
)

func init() {
	if app.IsServerMode() {
		fmt.Println("entity.mongodb starter init")
	}

	cli.ConfigureService(initMongodbConfigurator)
}
