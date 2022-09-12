package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/VadimGossip/crudFinManager/pkg/util"
	audit "github.com/VadimGossip/grpcAuditLog/pkg/domain"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/domain"
	"github.com/golang-jwt/jwt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type TokensRepository interface {
	Create(ctx context.Context, token domain.Token) error
	Get(ctx context.Context, token string) (domain.Token, error)
}

type UserAuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type Users struct {
	userRepo    UsersRepository
	tokenRepo   TokensRepository
	auditClient UserAuditClient
	hasher      PasswordHasher
	hmacSecret  []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewUsers(userRepo UsersRepository, tokenRepo TokensRepository, auditClient UserAuditClient, hasher PasswordHasher, secret []byte, accessTTL, refreshTTL time.Duration) *Users {
	return &Users{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		auditClient: auditClient,
		hasher:      hasher,
		hmacSecret:  secret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Surname:      inp.Surname,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}
	if s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	user, err = s.userRepo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_REGISTER,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.ID),
		AuthorID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignUp",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}

		return "", "", err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  int64(user.ID),
		AuthorID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignIn",
		}).Error("failed to send log request:", err)
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *Users) ParseToken(token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}
func (s *Users) GetRefreshTokenTTL() time.Duration {
	return s.refreshTTL
}

func (s *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := s.tokenRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return s.generateTokens(ctx, token.UserID)
}

func (s *Users) generateTokens(ctx context.Context, userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(userId),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.accessTTL).Unix(),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.NewRandString(32)
	if err != nil {
		return "", "", err
	}

	if err := s.tokenRepo.Create(ctx, domain.Token{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
