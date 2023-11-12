package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"inf_test/model"
	"log"
)

const DBName = "inf_test.db"

func init() {
	db := openDB()
	db.AutoMigrate(&model.User{})
}

func AddUser(user model.User) error {
	db := openDB()
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func SelectUser(id, password string) (model.User, error) {
	var user model.User
	db := openDB()
	if err := db.
		Where("user_id = ? AND password = ?", id, password).
		First(&user).
		Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func UpdateUser(newUser model.User) error {
	var user model.User
	db := openDB()
	if err := db.
		Where("user_id = ? AND password = ?", newUser.UserID, newUser.Password).
		First(&user).Error; err != nil {
		return err
	}
	if newUser.Nickname != "" {
		user.Nickname = newUser.Nickname
	}
	if newUser.Comment != "" {
		user.Comment = newUser.Comment
	}
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(id, password string) error {
	db := openDB()
	if err := db.
		Where("user_id = ? AND password = ?", id, password).
		Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

// データベース開く
func openDB() (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(DBName), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
