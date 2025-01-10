package mongodb

import (
	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/cli"
)

func init() {
	log.Logger.Info("entity.mongodb starter init")

	cli.ConfigureService(initMongodbConfigurator)
}
