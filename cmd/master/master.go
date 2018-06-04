//go:generate swagger generate spec
package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"test/test_swagger/master/model"
	"test/test_swagger/master/web"
)

type Flags struct {
	DBDataSource string
	DBDriver     string
}

func parseFlagsFromEnv() *Flags {
	flags := &Flags{
		DBDataSource: os.Getenv("MASTER_DB_DATASOURCE"),
		DBDriver:     os.Getenv("MASTER_DB_DRIVER"),
	}

	if flags.DBDataSource == "" {
		flags.DBDataSource = "test.db"
	}

	if flags.DBDriver == "" {
		flags.DBDriver = "sqlite3"
	}

	return flags
}

var components []web.Component = []web.Component{
	&web.SwaRouter{},
}

func main() {
	flag.Parse()
	flags := parseFlagsFromEnv()
	db, err := gorm.Open(flags.DBDriver, flags.DBDataSource)
	if err != nil {
		log.Fatal(err)
	}
	orm := &db

	orm.LogMode(true)
	model.InitSchema(orm)

	ms := map[string]web.MiddlewareFunc{}

	server := web.NewServer()

	for _, c := range components {
		c.Setup(orm, ms)
		server.RegisterRouter(c)
	}

	server.Run()
}
