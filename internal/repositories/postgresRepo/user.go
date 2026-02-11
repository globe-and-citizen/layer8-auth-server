package postgresRepo

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (r *PostgresRepository) UpdateUser(ctx context.Context, updates gormModels.User) error {
	tx := r.db.WithContext(ctx).Begin()
	user := gormModels.User{}
	if err := tx.Where("username = ?", updates.Username).First(&user).Error; err != nil {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("could not find user: %w", err))
	}

	err := tx.Model(&user).Updates(map[string]interface{}{
		"public_key": updates.PublicKey,
		"stored_key": updates.ScramStoredKey,
		"server_key": updates.ScramServerKey,
	}).Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("could not update user: %w", err))
	}

	userMetadata := gormModels.UserMetadata{
		ID:              user.ID,
		IsEmailVerified: false,
		DisplayName:     "",
		Color:           "",
	}

	err = tx.Create(&userMetadata).Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(fmt.Errorf("could not create user metadata entry: %e", err))
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, userId uint) (gormModels.User, error) {
	var user gormModels.User
	e := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if e != nil {
		return gormModels.User{}, utils.StackError(e)
	}

	return user, e
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (gormModels.User, error) {
	var user gormModels.User

	err := r.db.WithContext(ctx).Model(&gormModels.User{}).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		return gormModels.User{}, utils.StackError(err)
	}

	return user, nil
}

func (r *PostgresRepository) GetUserProfile(ctx context.Context, userID uint) (gormModels.User, gormModels.UserMetadata, error) {
	var user gormModels.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return gormModels.User{}, gormModels.UserMetadata{}, utils.StackError(err)
	}
	var userMetadata gormModels.UserMetadata
	if err := r.db.WithContext(ctx).Where("id = ?", userID).Find(&userMetadata).Error; err != nil {
		return gormModels.User{}, gormModels.UserMetadata{}, utils.StackError(err)
	}
	return user, userMetadata, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user gormModels.User) error {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return utils.StackError(fmt.Errorf("failed to create a new user: %v", err))
	}

	return nil
}

func (r *PostgresRepository) UpdateUserPassword(ctx context.Context, username string, storedKey string, serverKey string) error {
	return utils.StackError(r.db.WithContext(ctx).Model(&gormModels.User{}).
		Where("username=?", username).
		Updates(map[string]interface{}{"stored_key": storedKey, "server_key": serverKey}).Error)
}
