package mongodb

import (
	"github.com/abmpio/app/cli"
)

func init() {
	cli.ConfigureService(initMongodbConfigurator)
}
