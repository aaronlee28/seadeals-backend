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

type intervalConfig struct {
	DeliveredOrderToDone     string
	WaitingForSellerToRefund string
	OnDeliveryToDelivered    string
}

type AppConfig struct {
	AppName                  string
	BaseURL                  string
	Port                     string
	ENV                      string
	JWTSecret                []byte
	JWTExpiredInMinuteTime   int64
	DBConfig                 dbConfig
	Interval                 intervalConfig
	IntervalCron             string
	DatabaseURL              string
	MailJetPublicKey         string
	MailJetSecretKey         string
	SeaLabsPayMerchantCode   string
	SeaLabsPayAPIKey         string
	SeaLabsPayTransactionURL string
	SeaLabsPayRefundURL      string
	NgrokURL                 string
	AWSMail                  string
	ShippingKey              string
	ShippingURL              string
	ShippingActionOnError    string
	RedirectPaymentBase      string
	TZ                       string
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
		Interval: intervalConfig{
			DeliveredOrderToDone:     getEnv("INTERVAL_DELIVERED_TO_DONE", "2 day"),
			WaitingForSellerToRefund: getEnv("INTERVAL_WAIT_FOR_SELL_TO_REFUND", "3 day"),
			OnDeliveryToDelivered:    getEnv("INTERVAL_ON_DELIVERY_TO_DELIVERED", "1 day"),
		},
		IntervalCron:             getEnv("INTERVAL_CRON", "@daily"),
		DatabaseURL:              getEnv("DATABASE_URL", ""),
		MailJetPublicKey:         getEnv("MAILJET_PUBLIC_KEY", ""),
		MailJetSecretKey:         getEnv("MAILJET_SECRET_KEY", ""),
		SeaLabsPayMerchantCode:   getEnv("SEA_LABS_PAY_MERCHANT_CODE", ""),
		SeaLabsPayAPIKey:         getEnv("SEA_LABS_PAY_API_KEY", ""),
		SeaLabsPayTransactionURL: getEnv("SEA_LABS_PAY_TRANSACTION_URL", ""),
		SeaLabsPayRefundURL:      getEnv("SEA_LABS_PAY_REFUND_URL", ""),
		NgrokURL:                 getEnv("NGROK_URL", ""),
		AWSMail:                  getEnv("AWS_MAIL", ""),
		ShippingKey:              getEnv("SHIPPING_API_KEY", ""),
		ShippingURL:              getEnv("SHIPPING_API_URL", ""),
		ShippingActionOnError:    getEnv("SHIPPING_ACTION_ON_ERROR", ""),
		RedirectPaymentBase:      getEnv("REDIRECT_PAY_BASE", "http://localhost:3000"),
		TZ:                       getEnv("TZ", "Asia/Jakarta"),
	}
}

func getEnv(key, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}
