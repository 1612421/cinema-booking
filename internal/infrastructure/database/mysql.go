package database

import (
	"github.com/1612421/cinema-booking/pkg/go-kit/database/mysql"
	"gorm.io/gorm"

	"github.com/1612421/cinema-booking/config"
)

func NewMySQLDB(cfg *config.Config) (*gorm.DB, func(), error) {
	db, f, err := mysql.ConnectMySQL(cfg.MySQL, cfg.Service.Name)
	if err != nil {
		return nil, nil, err
	}

	sdb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	err = sdb.Ping()
	if err != nil {
		return nil, nil, err
	}

	return db, f, err
}
