package gormModels

type User struct {
	ID       uint   `gorm:"primaryKey; unique; autoIncrement; not null"`
	Username string `gorm:"column:username; unique; not null"`

	EmailVerificationCode string `gorm:"column:verification_code"`
	EmailZkProof          []byte `gorm:"column:email_proof"`
	EmailZkKeyPairId      uint   `gorm:"column:zk_key_pair_id"`

	PhoneNumberVerificationCode string `gorm:"column:phone_number_verification_code"`
	PhoneNumberZkProof          []byte `gorm:"column:phone_number_zk_proof"`
	PhoneNumberZkPairID         uint   `gorm:"column:phone_number_zk_pair_id"`

	PublicKey []byte `gorm:"column:public_key"`

	ScramSalt           string `gorm:"column:salt"`
	ScramIterationCount int    `gorm:"column:iteration_count"` // fixme? this is system configuration -> doesn't need to be stored per user? - if we change the configuration, all users will need to reset their passwords
	ScramServerKey      string `gorm:"column:server_key"`
	ScramStoredKey      string `gorm:"column:stored_key"`

	TelegramSessionIDHash []byte `gorm:"column:telegram_session_id_hash"`
}

func (User) TableName() string {
	return "users"
}
