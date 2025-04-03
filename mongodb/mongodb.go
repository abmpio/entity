package mongodb

import (
	"fmt"
	"time"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/configurationx"
	"github.com/abmpio/configurationx/options/mongodb"
	"github.com/abmpio/mongodbr"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongodbConfigurator(cliApp cli.CliApplication) {
	initMongodb()
}

// 初始化mongodb
func initMongodb() {
	mongodbOptions := configurationx.GetInstance().Mongodb

	if len(mongodbOptions.MongodbList) <= 0 {
		err := fmt.Errorf("无法初始化mongodb,没有配置好Mongodb的字符串连接")
		log.Logger.Error(err.Error())
		panic(err)
	}
	for eachKey, eachOption := range mongodbOptions.MongodbList {
		var client *mongo.Client
		var err error
		opts := make([]func(*options.ClientOptions), 0)
		if eachOption.EnableCommandMonitor != nil && *eachOption.EnableCommandMonitor {
			// enable command monitor
			opts = append(opts, mongodbr.EnableMongodbMonitor())
		}

		if eachKey == mongodb.AliasName_Default {
			if mongodbr.DefaultClient() == nil {
				client, err = mongodbr.SetupDefaultClient(eachOption.Uri, opts...)
				if err != nil {
					log.Logger.Error(fmt.Sprintf("key:%s,err:%s",
						eachKey,
						err.Error()))
					panic(err)
				}
			}
		} else {
			client, err = mongodbr.RegistClient(eachKey, eachOption.Uri, opts...)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("key:%s,err:%s",
					eachKey,
					err.Error()))
				panic(err)
			}
		}
		//测试ping
		for {
			err = mongodbr.Ping(client)
			if err == nil {
				break
			}
			log.Logger.Warn(fmt.Sprintf("key:%s,err:%s",
				eachKey,
				err.Error()))
			log.Logger.Warn("2s后重新测试...")
			time.Sleep(2 * time.Second)
		}
	}
}
