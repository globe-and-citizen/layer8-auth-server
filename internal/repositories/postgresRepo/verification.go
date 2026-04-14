package postgresRepo

import (
	"context"
	"database/sql"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

func (r *PostgresRepository) SaveProofOfEmailVerification(
	ctx context.Context, userId uint, salt string, verificationCode string, emailProof []byte, zkKeyPairId uint,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(&gormModels.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"email_salt":              salt,
			"email_verification_code": verificationCode,
			"email_zk_proof":          emailProof,
			"email_zk_id":             zkKeyPairId,
		}).Error

	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	err = tx.Where("user_id = ?", userId).
		Delete(&gormModels.EmailVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	err = tx.Model(&gormModels.UserMetadata{}).
		Where("id = ?", userId).
		Update("is_email_verified", true).
		Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveEmailVerificationData(
	ctx context.Context, data gormModels.EmailVerificationData,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(gormModels.EmailVerificationData{UserId: data.UserId}).
		Assign(data).
		FirstOrCreate(&gormModels.EmailVerificationData{}).
		Error

	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetEmailVerificationData(ctx context.Context, userId uint) (*gormModels.EmailVerificationData, error) {
	var data gormModels.EmailVerificationData
	e := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&data).Error
	if e != nil {
		return nil, utils.StackError(e)
	}

	return &data, nil
}

func (r *PostgresRepository) SavePhoneVerificationData(ctx context.Context, data gormModels.PhoneVerificationData) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(gormModels.PhoneVerificationData{UserId: data.UserId}).
		Assign(data).
		FirstOrCreate(&gormModels.PhoneVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetPhoneVerificationData(
	ctx context.Context,
	userID uint,
) (*gormModels.PhoneVerificationData, error) {
	var data gormModels.PhoneVerificationData
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return nil, utils.StackError(err)
	}

	return &data, nil
}

func (r *PostgresRepository) SaveProofOfPhoneVerification(
	ctx context.Context,
	userID uint,
	salt string,
	verificationCode string,
	zkProof []byte,
	zkPairID uint,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(&gormModels.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"phone_salt":              salt,
			"phone_verification_code": verificationCode,
			"phone_zk_proof":          zkProof,
			"phone_zk_id":             zkPairID,
		}).Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	err = tx.Where("user_id = ?", userID).
		Delete(&gormModels.PhoneVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	err = tx.Model(&gormModels.UserMetadata{}).
		Where("id = ?", userID).
		Update("is_phone_number_verified", true).
		Error
	if err != nil {
		tx.Rollback()
		return utils.StackError(err)
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveTelegramSessionIDHash(ctx context.Context, userID uint, sessionID []byte) error {
	return utils.StackError(r.db.WithContext(ctx).Model(&gormModels.User{}).
		Where("id = ?", userID).
		Update("telegram_session_id_hash", sessionID).
		Error)
}
