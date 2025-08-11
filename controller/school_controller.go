package controller

import (
	"0tak2/afterhee-server/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

type SchoolContoller interface {
	List(c *fiber.Ctx) error
}

type schoolContoller struct {
	svc service.SchoolService
}

func NewSchoolController(svc service.SchoolService) SchoolContoller {
	return &schoolContoller{svc: svc}
}

func (ctl *schoolContoller) List(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	schools, err := ctl.svc.GetSchools(keyword)
	if err != nil {
		log.Fatalln(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(schools)
}
