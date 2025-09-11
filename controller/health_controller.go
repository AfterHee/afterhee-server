package controller

import (
	"github.com/gofiber/fiber/v2"
)

// Controller Types
type HealthController interface {
	Check(c *fiber.Ctx) error
}

type healthController struct {
}

func NewHealthController() HealthController {
	return &healthController{}
}

// Suggest godoc
//
//	@Summary		헬스 체크
//	@Description	항상 HTTP 200 코드를 반환한다. 서버 상태를 체크하는 용도로 사용한다.
//	@Success 200	{string} string
//	@Router			/api/v1/healthcheck [get]
func (ctl *healthController) Check(c *fiber.Ctx) error {
	return c.Status(200).SendString("Okay")
}
