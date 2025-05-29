package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	GRPCPort    int
	GRPCTimeout time.Duration
}

func MustLoad() *Config {

	godotenv.Load(".env")

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
		log.Fatal("JWT_SECRET is not set")
	}

	grpcPortStr := os.Getenv("GRPC_PORT")
	if grpcPortStr == "" {
		log.Fatal("GRPC_PORT is not set")
	}

	grpcPort, err := strconv.Atoi(grpcPortStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_PORT value: %v", err)
	}

	grpcTimeoutStr := os.Getenv("GRPC_TIMEOUT")
	if grpcTimeoutStr == "" {
		log.Fatal("GRPC_TIMEOUT is not set")
	}

	grpcTimeout, err := time.ParseDuration(grpcTimeoutStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_TIMEOUT value: %v", err)
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
		GRPCPort:    grpcPort,
		GRPCTimeout: grpcTimeout,
	}
}
