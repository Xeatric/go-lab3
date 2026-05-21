package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ErrorHandler middleware для обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			statusCode, errorResponse := handleError(err)
			c.JSON(statusCode, errorResponse)
			c.Abort()
		}
	}
}

// Logger middleware для логирования запросов
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Логируем входящий запрос
		log.Printf("[%s] %s %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())

		c.Next()

		// Логируем время выполнения
		log.Printf("Request processed in %v", time.Since(start))
	}
}

func handleError(err error) (int, ErrorResponse) {
	// Проверка на ошибки валидации
	if strings.Contains(err.Error(), "binding") ||
		strings.Contains(err.Error(), "validation") {
		return http.StatusBadRequest, ErrorResponse{
			Error:   "Validation Error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	// Проверка на отсутствие записи
	if err == gorm.ErrRecordNotFound ||
		err.Error() == "tile not found" {
		return http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "The requested tile was not found",
			Code:    http.StatusNotFound,
		}
	}

	// Другие ошибки
	return http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal Server Error",
		Message: "Something went wrong",
		Code:    http.StatusInternalServerError,
	}
}
