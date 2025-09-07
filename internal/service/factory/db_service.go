package factory

import (
	"sync"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/database"
	"gorm.io/gorm"
)

var (
	databaseServiceInstanceLock = &sync.Mutex{}
	databaseServiceInstance     *gorm.DB
)

func NewDBService(cfg *config.Config) (*gorm.DB, error) {
	if databaseServiceInstance == nil {
		databaseServiceInstanceLock.Lock()
		defer databaseServiceInstanceLock.Unlock()
		db, err := database.Connect(cfg)
		if err != nil {
			return nil, err
		}
		databaseServiceInstance = db
	}
	return databaseServiceInstance, nil
}
