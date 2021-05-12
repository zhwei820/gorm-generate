package main

import (
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"time"
)

func connect(dsn string) (*gorm.DB, error) {
	var err error
	connection, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		log.Println(dsn)
		log.Println(err)
		log.Fatal("database configuration load error.")
	}

	if err != nil {
		return nil, errors.New("db connection error:" + err.Error())
	}
	//connection.LogMode(true)
	sqldb, _ := connection.DB()

	sqldb.SetConnMaxLifetime(time.Duration(300) * time.Second)
	sqldb.SetMaxOpenConns(200)
	sqldb.SetMaxIdleConns(50)
	return connection.Unscoped(), nil
}
