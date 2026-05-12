package database

import (
	"github.com/Tuxi4k/timesnap/internal/config"
	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if cfg.Database.Migrate {
		err = db.AutoMigrate(&deadline.Deadline{})
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
