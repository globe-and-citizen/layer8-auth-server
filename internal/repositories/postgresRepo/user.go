package postgresRepo

import (
	"fmt"
	gormModels2 "globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

func (r *PostgresRepository) UpdateUser(updates gormModels2.User) error {
	tx := r.db.Begin()
	user := gormModels2.User{}
	if err := tx.Where("username = ?", updates.Username).First(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not find user: %e", err)
	}

	err := tx.Model(&user).Updates(map[string]interface{}{
		"public_key": updates.PublicKey,
		"stored_key": updates.ScramStoredKey,
		"server_key": updates.ScramServerKey,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update user: %e", err)
	}

	userMetadata := gormModels2.UserMetadata{
		ID:              user.ID,
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

func (r *PostgresRepository) GetUserByID(userId uint) (gormModels2.User, error) {
	var user gormModels2.User
	e := r.db.Where("id = ?", userId).First(&user).Error

	if e != nil {
		return gormModels2.User{}, e
	}

	return user, e
}

func (r *PostgresRepository) GetUserByUsername(username string) (gormModels2.User, error) {
	var user gormModels2.User

	err := r.db.Model(&gormModels2.User{}).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return gormModels2.User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) GetUserProfile(userID uint) (gormModels2.User, gormModels2.UserMetadata, error) {
	var user gormModels2.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return gormModels2.User{}, gormModels2.UserMetadata{}, err
	}
	var userMetadata gormModels2.UserMetadata
	if err := r.db.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return gormModels2.User{}, gormModels2.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *PostgresRepository) PrecheckUserRegister(user gormModels2.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create a new user: %v", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return r.db.Model(&gormModels2.User{}).
		Where("username=?", username).
		Updates(map[string]interface{}{"stored_key": storedKey, "server_key": serverKey}).Error
}
