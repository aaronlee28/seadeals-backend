package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var Testing = "testing"

type dbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type AppConfig struct {
	AppName                  string
	BaseURL                  string
	Port                     string
	ENV                      string
	JWTSecret                []byte
	JWTExpiredInMinuteTime   int64
	DBConfig                 dbConfig
	DatabaseURL              string
	MailJetPublicKey         string
	MailJetSecretKey         string
	SeaLabsPayMerchantCode   string
	SeaLabsPayAPIKey         string
	SeaLabsPayTransactionURL string
	NgrokURL                 string
	AWSMail                  string
}

var Config = AppConfig{}

func Reset() {
	Config = AppConfig{
		AppName:                getEnv("APP_NAME", "Sea Deals"),
		BaseURL:                getEnv("BASE_URL", "localhost"),
		Port:                   getEnv("PORT", "8080"),
		ENV:                    getEnv("ENV", Testing),
		JWTSecret:              []byte(getEnv("JWT_SECRET", "")),
		JWTExpiredInMinuteTime: 15,
		DBConfig: dbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "seadeals_db"),
			Port:     getEnv("DB_PORT", "5432"),
		},
		DatabaseURL:              getEnv("DATABASE_URL", ""),
		MailJetPublicKey:         getEnv("MAILJET_PUBLIC_KEY", ""),
		MailJetSecretKey:         getEnv("MAILJET_SECRET_KEY", ""),
		SeaLabsPayMerchantCode:   getEnv("SEA_LABS_PAY_MERCHANT_CODE", ""),
		SeaLabsPayAPIKey:         getEnv("SEA_LABS_PAY_API_KEY", ""),
		SeaLabsPayTransactionURL: getEnv("SEA_LABS_PAY_TRANSACTION_URL", ""),
		NgrokURL:                 getEnv("NGROK_URL", ""),
		AWSMail:                  getEnv("AWS_MAIL", ""),
	}
}

func getEnv(key, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}
