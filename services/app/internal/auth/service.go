package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
)

type AuthService struct {
	usersRepository interfaces.IUsersRepository
	config          *config.Config
}

func NewService(usersRepository interfaces.IUsersRepository, config *config.Config) *AuthService {
	return &AuthService{
		usersRepository: usersRepository,
		config:          config,
	}
}

func (s *AuthService) Register(dto dtos.RegisterUserDto) (dtos.RegisterUserResponseDto, error) {
	_, found, err := s.usersRepository.FindByName(dto.Name)

	if err != nil {
		return dtos.RegisterUserResponseDto{}, err
	}

	if found {
		return dtos.RegisterUserResponseDto{}, ErrUserAlreadyExists
	}

	if dto.Password != dto.ConfirmedPassword {
		return dtos.RegisterUserResponseDto{}, ErrPasswordsDoNotMatch
	}

	passwordHash, err := argon2id.CreateHash(dto.Password, argon2id.DefaultParams)

	if err != nil {
		return dtos.RegisterUserResponseDto{}, ErrHashGeneration
	}

	user := entities.User{
		Id:           uuid.NewString(),
		Name:         dto.Name,
		PasswordHash: passwordHash,
	}

	s.usersRepository.Save(user)

	return dtos.RegisterUserResponseDto{
		Id: user.Id,
	}, nil
}

func (s *AuthService) Login(dto dtos.LoginUserDto) (dtos.LoginUserResponseDto, error) {
	user, found, err := s.usersRepository.FindByName(dto.Name)

	if err != nil {
		return dtos.LoginUserResponseDto{}, err
	}

	if !found {
		return dtos.LoginUserResponseDto{}, ErrUserNotFound
	}

	math, err := argon2id.ComparePasswordAndHash(dto.Password, user.PasswordHash)

	if err != nil {
		return dtos.LoginUserResponseDto{}, err
	}

	if !math {
		return dtos.LoginUserResponseDto{}, ErrInvalidCredentials
	}

	bytes, err := os.ReadFile(s.config.JwtPrivateKeyPath)

	if err != nil {
		return dtos.LoginUserResponseDto{}, fmt.Errorf("failed to read jwt private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)

	if err != nil {
		return dtos.LoginUserResponseDto{}, fmt.Errorf("failed to parse jwt private key: %w", err)
	}

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return dtos.LoginUserResponseDto{}, err
	}

	return dtos.LoginUserResponseDto{
		AccessToken: tokenString,
	}, nil
}
