package main

import "time"

type envConfig struct {
	DBPort             int           `mapstruct:"DB_PORT"`
	DBName             string        `mapstruct:"DB_NAME"`
	DBPassword         string        `mapstruct:"DB_PASSWORD"`
	HTTPServerPort     int           `mapstruct:"HTTP_SERVER_PORT"`
	HTTPServerHost     string        `mapstruct:"HTTP_SERVER_HOST"`
	RunMode            string        `mapstruct:"RUN_MODE"`
	MySQLLog           string        `mapstruct:"MYSQL_LOG"`
	AccessLog          string        `mapstruct:"ACCESS_LOG"`
	ViewPath           string        `mapstruct:"VIEW_PATH"`
	SmtpAuthAddress    string        `mapstruct:"SMTP_AUTH_ADDRESS"`
	SmtpServerAddress  string        `mapstruct:"SMTP_SERVER_ADDRESS"`
	MailSenderName     string        `mapstruct:"MAIL_SENDER_NAME"`
	MailSenderAddress  string        `mapstruct:"MAIL_SENDER_ADDRESS"`
	MailSenderPassword string        `mapstruct:"MAIL_SENDER_PASSWORD"`
	TokenSecretKey     string        `mapstruct:"TOKEN_SECRET_KEY"`
	MaxActivateTime    time.Duration `mapstruct:"MAX_ACTIVATE_TIME"`
	ShortDefaultExpire time.Duration `mapstruct:"SHORT_DEFAULT_EXPIRE"`
	RedisAddress       string        `mapstruct:"REDIS_ADDRESS"`
	SessionLifeTime    time.Duration `mapstruct:"SESSION_LIFETIME"`
}
