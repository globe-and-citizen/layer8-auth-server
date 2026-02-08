package postgresRepo

import (
	"context"
	"database/sql"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

func (r *PostgresRepository) SaveProofOfEmailVerification(
	ctx context.Context, userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(&gormModels.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"verification_code": verificationCode,
			"email_proof":       emailProof,
			"zk_key_pair_id":    zkKeyPairId,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("user_id = ?", userId).
		Delete(&gormModels.EmailVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&gormModels.UserMetadata{}).
		Where("id = ?", userId).
		Update("is_email_verified", true).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveEmailVerificationData(ctx context.Context, data gormModels.EmailVerificationData) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(gormModels.EmailVerificationData{UserId: data.UserId}).
		Assign(data).
		FirstOrCreate(&gormModels.EmailVerificationData{}).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetEmailVerificationData(ctx context.Context, userId uint) (gormModels.EmailVerificationData, error) {
	var data gormModels.EmailVerificationData
	e := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&data).Error
	if e != nil {
		return gormModels.EmailVerificationData{}, e
	}

	return data, nil
}

func (r *PostgresRepository) SavePhoneNumberVerificationData(ctx context.Context, data gormModels.PhoneNumberVerificationData) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(gormModels.PhoneNumberVerificationData{UserId: data.UserId}).
		Assign(data).
		FirstOrCreate(&gormModels.PhoneNumberVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetPhoneNumberVerificationData(
	ctx context.Context,
	userID uint,
) (gormModels.PhoneNumberVerificationData, error) {
	var data gormModels.PhoneNumberVerificationData
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return gormModels.PhoneNumberVerificationData{}, err
	}

	return data, nil
}

func (r *PostgresRepository) SaveProofOfPhoneNumberVerification(
	ctx context.Context,
	userID uint,
	phoneNumberVerificationCode string,
	phoneNumberZkProof []byte,
	phoneNumberZkPairID uint,
) error {
	tx := r.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(&gormModels.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"phone_number_verification_code": phoneNumberVerificationCode,
			"phone_number_zk_proof":          phoneNumberZkProof,
			"phone_number_zk_pair_id":        phoneNumberZkPairID,
		}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("user_id = ?", userID).
		Delete(&gormModels.PhoneNumberVerificationData{}).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&gormModels.UserMetadata{}).
		Where("id = ?", userID).
		Update("is_phone_number_verified", true).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveTelegramSessionIDHash(ctx context.Context, userID uint, sessionID []byte) error {
	return r.db.WithContext(ctx).Model(&gormModels.User{}).
		Where("id = ?", userID).
		Update("telegram_session_id_hash", sessionID).
		Error
}
