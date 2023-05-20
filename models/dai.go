package models

import "gorm.io/gorm"

type DataAccessInterface interface {
	Create(interface{}) *gorm.DB
	First(interface{}, ...interface{}) *gorm.DB
}

type MockedDataAccessInterface interface {
	gorm.DB
	DataAccessInterface
}
