package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"my-grpc/internal/lib/jwt"
	"my-grpc/internal/lib/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {

	const op = "service.auth.register"

	log := a.log.With(
		slog.String("op", op),
	)
	if email == "" {
		log.Error(ErrInvalidCredentials.Error())
		return 0, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	if password == "" {
		log.Error(ErrInvalidCredentials.Error())
		return 0, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil

}
func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	const op = "service.auth.login"

	log := a.log.With(
		slog.String("op", op),
	)

	if email == "" {
		log.Error(ErrInvalidCredentials.Error())
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	if password == "" {
		log.Error(ErrInvalidCredentials.Error())
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	if appId == 0 {
		log.Error(ErrInvalidCredentials.Error())
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	user, err := a.usrProvider.User(ctx, email)

	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error(err.Error())

		return "", fmt.Errorf("%s:%w", op, err)

	}

	app, err := a.appProvider.App(ctx, appId)

	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		log.Error(err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil

}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "service.auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	if userID == 0 {
		log.Error(ErrInvalidCredentials.Error())
		return false, fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)

	if err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return isAdmin, nil

}
