package config

import (
	"log"
	"os"
)

type Config struct {
	Env         string
	HttpAddress string
	LoginDB     string
	PasswordDB  string
	HostDB      string
	PortDB      string
	NameDB      string
	JWTSecret   string
}

func MustLoad() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is not set")
	}

	httpAddress := os.Getenv("HTTP_ADDRESS")
	if env == "" {
		log.Fatal("HTTP_ADDRESS is not set")
	}

	loginDB := os.Getenv("LOGIN_DB")
	if loginDB == "" {
		log.Fatal("LOGIN_DB is not set")
	}

	passwordDB := os.Getenv("PASSWORD_DB")
	if passwordDB == "" {
		log.Fatal("PASSWORD_DB is not set")
	}

	hostDB := os.Getenv("HOST_DB")
	if hostDB == "" {
		log.Fatal("HOST_DB is not set")
	}

	portDB := os.Getenv("PORT_DB")
	if portDB == "" {
		log.Fatal("PORT_DB is not set")
	}

	nameDB := os.Getenv("NAME_DB")
	if nameDB == "" {
		log.Fatal("NAME_DB is not set")
	}

	jwtSercet := os.Getenv("JWT_SECRET")
	if nameDB == "" {
		log.Fatal("NAME_DB is not set")
	}

	return &Config{
		Env:         env,
		HttpAddress: httpAddress,
		LoginDB:     loginDB,
		PasswordDB:  passwordDB,
		HostDB:      hostDB,
		PortDB:      portDB,
		NameDB:      nameDB,
		JWTSecret:   jwtSercet,
	}
}
