package data

import (
	"comment/internal/conf"
	"comment/pkg/log"
	"gorm.io/gorm"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCommentRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	log.Info(nil, "init Data.", "conf.Data", c)

	if err := c.ValidateAll(); err != nil {
		log.Fatal(nil, "validate config error.", "err", err)
		return nil, nil, err
	}
	log.Info(nil, "validate conf.Data successful.")

	cleanup := func() {
		log.Info(nil, "closing the data resources")
	}

	data := &Data{}
	if c.Database.Driver == "mysql" || c.Database.Driver == "" {
		// 使用 mysql
		db, err := NewDB(c.Database)
		if err != nil {
			log.Fatal(nil, "new db error.", "err", err)
			return nil, nil, err
		}
		data.db = db
		log.Info(nil, "new db successful.")
	} else {
		log.Fatal(nil, "database driver error.", "driver", c.Database.Driver)
	}

	return data, cleanup, nil
}
