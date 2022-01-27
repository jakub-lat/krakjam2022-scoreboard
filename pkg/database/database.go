package database

import "gorm.io/gorm"

type DB struct {
	*gorm.DB
}

func NewDB(db *gorm.DB) (*DB, error) {
	d := &DB{db}
	if err := d.AutoMigrate(&Player{}, &GameRun{}, &GameRunLevel{}); err != nil {
		return nil, err
	}
	return d, nil
}
