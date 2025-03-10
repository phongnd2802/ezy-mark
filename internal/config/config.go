package config

import "github.com/spf13/viper"

type Config struct {
	Mode           string `mapstructure:"MODE"`
	DBSource       string `mapstructure:"DB_SOURCE"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	SessionSecret  string `mapstructure:"SESSION_SECRET"`
	SenderName     string `mapstructure:"SENDER_NAME"`
	SenderEmail    string `mapstructure:"SENDER_EMAIL"`
	SenderPassword string `mapstructure:"SENDER_PASSWORD"`
	SecretKeyJWT   string `mapstructure:"SECRET_KEY_JWT"`
	MinIOEndPoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOUseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
	MinIOHost      string `mapstructure:"MINIO_HOST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
