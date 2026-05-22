package dto

import "time"

// Register Request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Name     string `json:"name" binding:"required,min=2,max=255"`
}

// Login Request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Refresh Request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Auth Response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
}

// Forgot Password Request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Reset Password Request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=100"`
}

// Change Password Request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=100"`
}

// User Response (безопасный)
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Whoami Response
type WhoamiResponse struct {
	Authenticated bool          `json:"authenticated"`
	User          *UserResponse `json:"user,omitempty"`
}
