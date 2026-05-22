package repository

import (
	"time"

	"paving-tiles-api/internal/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// User operations
func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByYandexID ищет пользователя по ID от Яндекс
func (r *AuthRepository) FindUserByYandexID(yandexID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("yandex_id = ? AND deleted_at IS NULL", yandexID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByVkID(vkID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("vk_id = ? AND deleted_at IS NULL", vkID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *AuthRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}

// Token operations
func (r *AuthRepository) SaveToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) FindTokenByHash(tokenHash string) (*models.Token, error) {
	var token models.Token
	err := r.db.Where("token_hash = ? AND is_revoked = false AND expires_at > ?",
		tokenHash, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) RevokeToken(tokenID uint) error {
	now := time.Now()
	return r.db.Model(&models.Token{}).Where("id = ?", tokenID).Updates(map[string]interface{}{
		"is_revoked": true,
		"revoked_at": now,
	}).Error
}

func (r *AuthRepository) RevokeAllUserTokens(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.Token{}).Where("user_id = ? AND is_revoked = false", userID).Updates(map[string]interface{}{
		"is_revoked": true,
		"revoked_at": now,
	}).Error
}

func (r *AuthRepository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Token{}).Error
}
