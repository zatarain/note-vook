package models

import "gorm.io/gorm"

type DataAccessInterface interface {
	AutoMigrate(...interface{}) error
	Create(interface{}) *gorm.DB
	Delete(interface{}, ...interface{}) *gorm.DB
	First(interface{}, ...interface{}) *gorm.DB
	Find(interface{}, ...interface{}) *gorm.DB
	Model(interface{}) *gorm.DB
	Updates(interface{}) *gorm.DB
	Where(interface{}, ...interface{}) *gorm.DB
}

type MockedDataAccessInterface interface {
	gorm.DB
	DataAccessInterface
}
