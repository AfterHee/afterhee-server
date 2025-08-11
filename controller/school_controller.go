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

// List godoc
//
//	@Summary		학교 목록 조회
//	@Description	학교 목록을 키워드를 바탕으로 조회한다
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query	string	true	"검색 키워드"
//	@Success		200	{object}	[]service.School
//	@Router			/api/v1/schools [get]
func (ctl *schoolContoller) List(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	schools, err := ctl.svc.GetSchools(keyword)
	if err != nil {
		log.Fatalln(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(schools)
}
