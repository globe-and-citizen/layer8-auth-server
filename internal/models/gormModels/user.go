package gormModels

type User struct {
	ID       uint   `gorm:"column:id; primaryKey; unique; autoIncrement; not null"`
	Username string `gorm:"column:username; unique; not null"`

	EmailZkSalt           string `gorm:"column:email_salt"`
	EmailVerificationCode string `gorm:"column:email_verification_code"`
	EmailZkProof          []byte `gorm:"column:email_zk_proof"`
	EmailZkKeyPairID      uint   `gorm:"column:email_zk_id"`

	PhoneZkSalt           string `gorm:"column:phone_salt"`
	PhoneVerificationCode string `gorm:"column:phone_verification_code"`
	PhoneZkProof          []byte `gorm:"column:phone_zk_proof"`
	PhoneZkKeyPairID      uint   `gorm:"column:phone_zk_id"`

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
