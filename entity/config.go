package entity

type PostgresConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
}

type ServerConfig struct {
	Host      string `env:"HOST"`
	Port      int    `env:"PORT"`
	JWTSecret string `env:"JWT_SECRET"`
}
