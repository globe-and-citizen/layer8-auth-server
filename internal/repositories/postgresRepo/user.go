package postgresRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/dto/tmp"
	"globe-and-citizen/layer8/auth-server/internal/models"
)

func (r *PostgresRepository) AddUser(req tmp.UserRegisterDTO) error {
	newUser := models.User{}

	tx := r.db.Begin()
	if err := tx.Where("username = ?", req.Username).First(&newUser).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not find user: %e", err)
	}

	err := tx.Model(&r).Updates(map[string]interface{}{
		"public_key": req.PublicKey,
		"stored_key": req.StoredKey,
		"server_key": req.ServerKey,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update user: %e", err)
	}

	userMetadata := models.UserMetadata{
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

func (r *PostgresRepository) FindUserByID(userId uint) (models.User, error) {
	var user models.User
	e := r.db.Where("id = ?", userId).First(&user).Error

	if e != nil {
		return models.User{}, e
	}

	return user, e
}

func (r *PostgresRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	err := r.db.Model(&models.User{}).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) GetUserProfile(userID uint) (models.User, models.UserMetadata, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, models.UserMetadata{}, err
	}
	var userMetadata models.UserMetadata
	if err := r.db.Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return models.User{}, models.UserMetadata{}, err
	}
	return user, userMetadata, nil
}

func (r *PostgresRepository) PrecheckUserRegister(req tmp.RegisterUserPrecheckDTO, salt string, iterCount int) error {
	user := models.User{
		Username:       req.Username,
		Salt:           salt,
		IterationCount: iterCount,
		PublicKey:      []byte{},
	}

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create a new user: %v", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateUserPassword(username string, storedKey string, serverKey string) error {
	return r.db.Model(&models.User{}).
		Where("username=?", username).
		Updates(map[string]interface{}{"stored_key": storedKey, "server_key": serverKey}).Error
}
