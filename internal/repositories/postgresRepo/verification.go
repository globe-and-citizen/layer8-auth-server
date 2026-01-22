package postgresRepo

import (
	"database/sql"
	gormModels2 "globe-and-citizen/layer8/auth-server/internal/models/gormModels"
)

func (r *PostgresRepository) SaveProofOfEmailVerification(
	userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(
		&gormModels2.User{},
	).Where(
		"id = ?", userId,
	).Updates(map[string]interface{}{
		"verification_code": verificationCode,
		"email_proof":       emailProof,
		"zk_key_pair_id":    zkKeyPairId,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(
		"user_id = ?", userId,
	).Delete(
		&gormModels2.EmailVerificationData{},
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(
		&gormModels2.UserMetadata{},
	).Where(
		"id = ?", userId,
	).Update("is_email_verified", true).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveEmailVerificationData(data gormModels2.EmailVerificationData) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(
		gormModels2.EmailVerificationData{UserId: data.UserId},
	).Assign(data).FirstOrCreate(&gormModels2.EmailVerificationData{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetEmailVerificationData(userId uint) (gormModels2.EmailVerificationData, error) {
	var data gormModels2.EmailVerificationData
	e := r.db.Where("user_id = ?", userId).First(&data).Error
	if e != nil {
		return gormModels2.EmailVerificationData{}, e
	}

	return data, nil
}

func (r *PostgresRepository) SavePhoneNumberVerificationData(
	data gormModels2.PhoneNumberVerificationData,
) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(
		gormModels2.PhoneNumberVerificationData{UserId: data.UserId},
	).Assign(data).FirstOrCreate(&gormModels2.PhoneNumberVerificationData{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetPhoneNumberVerificationData(userID uint) (gormModels2.PhoneNumberVerificationData, error) {
	var data gormModels2.PhoneNumberVerificationData
	err := r.db.Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return gormModels2.PhoneNumberVerificationData{}, err
	}

	return data, nil
}

func (r *PostgresRepository) SaveProofOfPhoneNumberVerification(
	userID uint,
	phoneNumberVerificationCode string,
	phoneNumberZkProof []byte,
	phoneNumberZkPairID uint,
) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(
		&gormModels2.User{},
	).Where(
		"id = ?", userID,
	).Updates(map[string]interface{}{
		"phone_number_verification_code": phoneNumberVerificationCode,
		"phone_number_zk_proof":          phoneNumberZkProof,
		"phone_number_zk_pair_id":        phoneNumberZkPairID,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where(
		"user_id = ?", userID,
	).Delete(
		&gormModels2.PhoneNumberVerificationData{},
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(
		&gormModels2.UserMetadata{},
	).Where(
		"id = ?", userID,
	).Update("is_phone_number_verified", true).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) SaveTelegramSessionIDHash(userID uint, sessionID []byte) error {
	return r.db.Model(
		&gormModels2.User{},
	).Where(
		"id = ?", userID,
	).Update("telegram_session_id_hash", sessionID).Error
}
