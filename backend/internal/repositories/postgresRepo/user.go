package postgresRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
)

func (r *PostgresRepository) AddUser(newUser gormModels.User) error {
	tx := r.db.Begin()
	if err := tx.Where("username = ?", newUser.Username).First(&newUser).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not find user: %e", err)
	}

	err := tx.Model(&r).Updates(map[string]interface{}{
		"public_key": newUser.PublicKey,
		"stored_key": newUser.StoredKey,
		"server_key": newUser.ServerKey,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update user: %e", err)
	}

	userMetadata := gormModels.UserMetadata{
		ID:              newUser.ID,
		IsEmailVerified: false,
		DisplayName:     "",
		Color:           "",
	}

	err = tx.Create(&userMetadata).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create user metadata entry: %e", err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) FindUserByID(userId uint) (gormModels.User, error) {
	var user gormModels.User
	e := r.db.Where("id = ?", userId).First(&user).Error

	if e != nil {
		return gormModels.User{}, e
	}

	return user, e
}

func (r *PostgresRepository) GetUserByUsername(username string) (gormModels.User, error) {
	var user gormModels.User

	err := r.db.Model(&gormModels.User{}).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return gormModels.User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) GetUserProfile(userID uint) (gormModels.User, gormModels.UserMetadata, error) {
	var user gormModels.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return gormModels.User{}, gormModels.UserMetadata{}, err
	}
	var userMetadata gormModels.UserMetadata
	if err := r.db.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return gormModels.User{}, gormModels.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *PostgresRepository) PrecheckUserRegister(user gormModels.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create a new user: %v", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return r.db.Model(&gormModels.User{}).
		Where("username=?", username).
		Updates(map[string]interface{}{"stored_key": storedKey, "server_key": serverKey}).Error
}
