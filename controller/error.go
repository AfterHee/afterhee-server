package controller

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Error Handler

func GlobalErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := MessageOfCode(code)

	// Retrieve the custom status code if it's a *fiber.Error
	var e ControllerError
	if errors.As(err, &e) {
		code = e.StatusCode
		message = "[" + e.DetailCode + "] " + e.Message
	}

	ctx.Status(code)
	return UTF8Json(ctx, ResponseOf(
		true,
		message,
		nil,
	))
}

// Error

type ControllerError struct {
	StatusCode int
	DetailCode string
	Message    string
}

func (e ControllerError) Error() string {
	return strconv.Itoa(e.StatusCode) + " (" + e.DetailCode + ") " + e.Message
}

func ErrorOf(httpStatusCode int, detailCode string, message string) error {
	return ControllerError{
		StatusCode: httpStatusCode,
		DetailCode: detailCode,
		Message:    message,
	}
}
