package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Twilio   Twilio
	Cookie   Cookie
	Store    Store
	Session  Session
	Metrics  Metrics
	Logger   Logger
	AWS      AWS
	Jaeger   Jaeger
	Firebase Firebase
	OAuth    OAuth
	Nats     NatsConfig
	Stripe   Stripe
	Paypal   Paypal
	Adyen    Adyen
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	EmailFrom         string
	SMTPHost          string
	SMTPPort          string
	SMTPUser          string
	SMTPPassword      string
	WebhookURL        string
	PublicKeyPath     string
	Environment       string
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// NATS config
type NatsConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Twilio config
type Twilio struct {
	AccountSID string
	AuthToken  string
	From       string
}

// Firebase config struct
type Firebase struct {
	Type                string `mapstructure:"type" json:"type"`
	ProjectID           string `mapstructure:"project_id" json:"project_id"`
	PrivateKeyID        string `mapstructure:"private_key_id" json:"private_key_id"`
	PrivateKey          string `mapstructure:"private_key" json:"private_key"`
	ClientEmail         string `mapstructure:"client_email" json:"client_email"`
	ClientID            string `mapstructure:"client_id" json:"client_id"`
	AuthURI             string `mapstructure:"auth_uri" json:"auth_uri"`
	TokenURI            string `mapstructure:"token_uri" json:"token_uri"`
	AuthProviderCertURL string `mapstructure:"auth_provider_x509_cert_url" json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `mapstructure:"client_x509_cert_url" json:"client_x509_cert_url"`
}

// OAuth config struct
type OAuth struct {
	Google      OAuthProvider
	Facebook    OAuthProvider
	Apple       OAuthProvider
	StateString string
}

type OAuthProvider struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	Enabled      bool   `mapstructure:"enabled"`
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Store config
type Store struct {
	ImagesFolder string
}

// AWS S3
type AWS struct {
	Endpoint       string
	MinioAccessKey string
	MinioSecretKey string
	UseSSL         bool
	MinioEndpoint  string
}

// AWS S3
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// PSP
type Stripe struct {
	Secret string
	Key    string
}

type Paypal struct {
	Secret string
	Key    string
}

type Adyen struct {
	Secret string
	Key    string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".") 

	// Enable reading environment variables
	v.SetEnvPrefix("POSTGRES") // Prefix env vars with POSTGRES_
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	_ = v.BindEnv("postgres.postgresqlhost")
	_ = v.BindEnv("postgres.postgresqlport")
	_ = v.BindEnv("postgres.postgresqluser")
	_ = v.BindEnv("postgres.postgresqlpassword")
	_ = v.BindEnv("postgres.postgresqldbname")
	_ = v.BindEnv("postgres.postgresqlsslmode")

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
