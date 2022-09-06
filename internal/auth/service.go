package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/turbak/joom-calendar/internal/creating"
	"net/http"
	"strings"
	"time"
)

type Storage interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user creating.User) (int, error)
}

type Authorizer interface {
	GetAccessToken(code string) (string, error)
	GetUser(accessToken string) (*User, error)
}

type Service struct {
	storage Storage

	authorizer Authorizer
	jwtKey     []byte
}

func NewService(storage Storage, authorizer Authorizer, jwtKey []byte) *Service {
	return &Service{
		storage:    storage,
		authorizer: authorizer,
		jwtKey:     jwtKey,
	}
}

func (s *Service) AuthenticateGithub(ctx context.Context, code string) (string, error) {
	token, err := s.authorizer.GetAccessToken(code)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %v", err)
	}
	user, err := s.authorizer.GetUser(token)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %v", err)
	}

	foundUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			user.ID, err = s.storage.CreateUser(ctx, creating.User{
				Name:  user.Name,
				Email: user.Email,
			})
		} else {
			return "", fmt.Errorf("failed to get user: %v", err)
		}
	}

	if foundUser != nil {
		user.ID = foundUser.ID
	}

	return s.generateToken(*user)
}

func (s *Service) generateToken(user User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return tokenString, nil
}

func (s *Service) validateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.jwtKey, nil
		},
	)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return errors.New("couldn't parse claims")
	}

	if err = claims.Valid(); err != nil {
		return fmt.Errorf("claims not valid: %v", err)
	}

	return nil
}

func (s *Service) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "token is empty", http.StatusUnauthorized)
				return
			}

			if strings.HasPrefix(token, "Bearer ") {
				token = strings.TrimPrefix(token, "Bearer ")
			}

			if err := s.validateToken(token); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
