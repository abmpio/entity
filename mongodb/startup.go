package mongodb

import (
	"fmt"

	"github.com/abmpio/app/cli"
)

func init() {
	fmt.Println("entity.mongodb starter init")

	cli.ConfigureService(initMongodbConfigurator)
}
