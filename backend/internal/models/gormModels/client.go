package gormModels

type Client struct {
	ID                  string `gorm:"column:id; not null" json:"id"`
	Secret              string `gorm:"column:secret" json:"secret"`
	Name                string `gorm:"column:name" json:"name"`
	RedirectURI         string `gorm:"column:redirect_uri" json:"redirect_uri"`
	BackendURI          string `gorm:"column:backend_uri" json:"backend_uri"`
	Username            string `gorm:"column:username; unique; not null" json:"username"`
	ScramSalt           string `gorm:"column:salt; not null" json:"salt"`
	ScramIterationCount int    `gorm:"column:iteration_count" json:"iteration_count"`
	ScramServerKey      string `gorm:"column:server_key" json:"server_key"`
	ScramStoredKey      string `gorm:"column:stored_key" json:"stored_key"`
	NTorX509Certificate string `gorm:"column:x509_certificate_bytes" json:"x509_certificate_bytes"`
}

func (Client) TableName() string {
	return "clients"
}
