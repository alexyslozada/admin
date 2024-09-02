package db

import "gitlab.com/EDteam/workshop-ai-2024/admin/domain"

func RunMigration() {
	db, err := NewGorm()
	if err != nil {
		panic(err)
	}

	// Users table
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		panic(err)
	}

	// Insert first user
	user := domain.User{
		Email:    "admin@ed.team",
		Password: "$2a$12$W.oQN6MlDwpkQt.v1hmgmeMDAHwhrbFJVituXFgWoIt5tSmj8eGEG",
	}
	err = db.Create(&user).Error
	if err != nil {
		panic(err)
	}
}
