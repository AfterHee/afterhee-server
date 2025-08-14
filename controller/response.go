package controller

import (
	"github.com/gofiber/fiber/v2"
)

// Resonse DTO

type CommonResponse struct {
	IsError bool        `json:"isError"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var errorMessageForStatusCode = map[int]string{
	fiber.StatusOK:                  "요청을 잘 수행했습니다",
	fiber.StatusBadRequest:          "잘못된 요청입니다",
	fiber.StatusInternalServerError: "오류가 발생했습니다",
}

func ResponseOf(isError bool, message string, data interface{}) *CommonResponse {
	return &CommonResponse{
		IsError: isError,
		Message: message,
		Data:    data,
	}
}

func ResponseOfCode(isError bool, statusCode int, data interface{}) *CommonResponse {
	return &CommonResponse{
		IsError: isError,
		Message: MessageOfCode(statusCode),
		Data:    data,
	}
}

func MessageOfCode(statusCode int) string {
	if message, ok := errorMessageForStatusCode[statusCode]; ok {
		return message
	}
	return "오류가 발생했습니다"
}
