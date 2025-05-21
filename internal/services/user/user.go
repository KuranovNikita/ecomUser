package user

import (
	"context"
	"ecomUser/internal/domain/models"
	"ecomUser/internal/lib/jwtLib"
	"ecomUser/internal/storage/postgres"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
	secret       string
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		login string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
		log:          log,
	}
}

func (a *Auth) Login(ctx context.Context, login string, password string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", login),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			a.log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwtLib.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		a.log.Error("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}

func (a *Auth) SaveUser(ctx context.Context, email string, login string, password string) (int64, error) {
	const op = "Auth.SaveUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, login, passHash)

	if err != nil {
		log.Error("failed to save user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
