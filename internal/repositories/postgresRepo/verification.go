package postgresRepo

import (
	"database/sql"
	"globe-and-citizen/layer8/auth-server/internal/models"
)

func (r *PostgresRepository) SaveProofOfEmailVerification(
	userId uint, verificationCode string, emailProof []byte, zkKeyPairId uint,
) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Model(
		&models.User{},
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
		&models.EmailVerificationData{},
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(
		&models.UserMetadata{},
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

func (r *PostgresRepository) SaveEmailVerificationData(data models.EmailVerificationData) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(
		models.EmailVerificationData{UserId: data.UserId},
	).Assign(data).FirstOrCreate(&models.EmailVerificationData{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetEmailVerificationData(userId uint) (models.EmailVerificationData, error) {
	var data models.EmailVerificationData
	e := r.db.Where("user_id = ?", userId).First(&data).Error
	if e != nil {
		return models.EmailVerificationData{}, e
	}

	return data, nil
}

func (r *PostgresRepository) SavePhoneNumberVerificationData(
	data models.PhoneNumberVerificationData,
) error {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelReadCommitted})

	err := tx.Where(
		models.PhoneNumberVerificationData{UserId: data.UserId},
	).Assign(data).FirstOrCreate(&models.PhoneNumberVerificationData{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepository) GetPhoneNumberVerificationData(userID uint) (models.PhoneNumberVerificationData, error) {
	var data models.PhoneNumberVerificationData
	err := r.db.Where("user_id = ?", userID).First(&data).Error
	if err != nil {
		return models.PhoneNumberVerificationData{}, err
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
		&models.User{},
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
		&models.PhoneNumberVerificationData{},
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(
		&models.UserMetadata{},
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
		&models.User{},
	).Where(
		"id = ?", userID,
	).Update("telegram_session_id_hash", sessionID).Error
}
