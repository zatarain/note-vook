package models

import "gorm.io/gorm"

type DataAccessInterface interface {
	Create(interface{}) *gorm.DB
	First(interface{}, ...interface{}) *gorm.DB
	Find(interface{}, ...interface{}) *gorm.DB
	AutoMigrate(...interface{}) error
}

type MockedDataAccessInterface interface {
	gorm.DB
	DataAccessInterface
}
