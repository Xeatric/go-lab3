package service

import (
	"errors"

	"math/rand"
	"time"

	"paving-tiles-api/internal/auth/repository"
	"paving-tiles-api/internal/auth/utils"
	"paving-tiles-api/internal/config"
	"paving-tiles-api/internal/dto"
	"paving-tiles-api/internal/models"
)

type AuthService struct {
	repo   *repository.AuthRepository
	config *config.Config
}

// NewAuthService - конструктор
func NewAuthService(repo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:   repo,
		config: cfg,
	}
}

// Register - регистрация нового пользователя
func (s *AuthService) Register(req *dto.RegisterRequest, userAgent, ip string) (*dto.AuthResponse, error) {
	// Проверка существования пользователя
	existingUser, _ := s.repo.FindUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хеширование пароля
	passwordHash, salt, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Создание пользователя
	user := &models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Salt:         salt,
		Name:         req.Name,
		Role:         "user",
		IsActive:     true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Генерация токенов
	return s.generateTokens(user, userAgent, ip)
}

// Login - вход пользователя
func (s *AuthService) Login(req *dto.LoginRequest, userAgent, ip string) (*dto.AuthResponse, error) {
	// Поиск пользователя
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверка пароля
	if !utils.CheckPassword(req.Password, user.PasswordHash, user.Salt) {
		return nil, errors.New("invalid email or password")
	}

	// Обновление времени последнего входа
	s.repo.UpdateLastLogin(user.ID)

	// Генерация токенов
	return s.generateTokens(user, userAgent, ip)
}

// Refresh - обновление токенов
func (s *AuthService) Refresh(refreshToken, userAgent, ip string) (*dto.AuthResponse, error) {
	// Хеширование refresh token
	tokenHash := utils.HashToken(refreshToken)

	// Поиск токена в БД
	storedToken, err := s.repo.FindTokenByHash(tokenHash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Поиск пользователя
	user, err := s.repo.FindUserByID(storedToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Отзыв старого токена
	s.repo.RevokeToken(storedToken.ID)

	// Генерация новых токенов
	return s.generateTokens(user, userAgent, ip)
}

// Logout - выход из текущей сессии
func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := utils.HashToken(refreshToken)
	storedToken, err := s.repo.FindTokenByHash(tokenHash)
	if err != nil {
		return errors.New("token not found")
	}

	return s.repo.RevokeToken(storedToken.ID)
}

// LogoutAll - выход из всех сессий
func (s *AuthService) LogoutAll(userID uint) error {
	return s.repo.RevokeAllUserTokens(userID)
}

// Whoami - получение информации о текущем пользователе
func (s *AuthService) Whoami(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

// FindOrCreateUser - поиск или создание пользователя через OAuth
func (s *AuthService) FindOrCreateUser(email, name, provider, providerID string) (*models.User, error) {
	var user *models.User
	var err error

	// Сначала ищем по ID провайдера
	if provider == "yandex" {
		user, err = s.repo.FindUserByYandexID(providerID)
	}

	// Если не нашли, ищем по email
	if user == nil {
		user, err = s.repo.FindUserByEmail(email)
	}

	// Если пользователь существует и у него нет ID провайдера, обновляем его
	if user != nil && user.YandexID == "" && provider == "yandex" {
		user.YandexID = providerID
		err = s.repo.UpdateUser(user)
		if err != nil {
			return nil, err
		}
	}

	// Если пользователь не найден, создаем нового
	if user == nil {
		// Генерируем случайный пароль для OAuth пользователя
		randomPassword := generateRandomString(32)
		passwordHash, salt, err := utils.HashPassword(randomPassword)
		if err != nil {
			return nil, err
		}

		user = &models.User{
			Email:        email,
			PasswordHash: passwordHash,
			Salt:         salt,
			Name:         name,
			YandexID:     providerID,
			Role:         "user",
			IsActive:     true,
		}
		err = s.repo.CreateUser(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GenerateTokensForUser - генерация токенов для существующего пользователя
func (s *AuthService) GenerateTokensForUser(user *models.User, userAgent, ip string) (*dto.AuthResponse, error) {
	return s.generateTokens(user, userAgent, ip)
}

// generateTokens - генерация пары токенов
func (s *AuthService) generateTokens(user *models.User, userAgent, ip string) (*dto.AuthResponse, error) {
	// Генерация Access Token (JWT)
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		s.config.JWTAccessSecret,
		s.config.JWTAccessExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Генерация Refresh Token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Хеширование и сохранение Refresh Token в БД
	tokenHash := utils.HashToken(refreshToken)
	token := &models.Token{
		UserID:    user.ID,
		TokenHash: tokenHash,
		TokenType: "refresh",
		ExpiresAt: time.Now().Add(s.config.JWTRefreshExpiration),
		IsRevoked: false,
		UserAgent: userAgent,
		IPAddress: ip,
	}

	if err := s.repo.SaveToken(token); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateAccessToken - валидация токена
func (s *AuthService) ValidateAccessToken(tokenString string) (*utils.JWTClaims, error) {
	return utils.ValidateAccessToken(tokenString, s.config.JWTAccessSecret)
}

// generateRandomString - генерация случайной строки
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
