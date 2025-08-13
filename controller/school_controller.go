package controller

import (
	"0tak2/afterhee-server/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

type SchoolContoller interface {
	List(c *fiber.Ctx) error
	ListMeals(c *fiber.Ctx) error
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
//	@Success 200 {object} controller.CommonResponse{data=[]service.School}
//	@Router			/api/v1/schools [get]
func (ctl *schoolContoller) List(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	schools, err := ctl.svc.GetSchools(keyword)
	if err != nil {
		log.Println(err)
		return ErrorOf(fiber.StatusInternalServerError, "50000", MessageOfCode(fiber.StatusInternalServerError))
	}
	return UTF8Json(c, ResponseOfCode(false, fiber.StatusOK, schools))
}

// List godoc
//
//	@Summary		식단표 조회
//	@Description	특정 학교의 월별 식단표를 제공한다
//	@Accept			json
//	@Produce		json
//	@Param			eduOfficeCode	query	string	true	"시도교육청코드"
//	@Param			schoolCode	query	string	true	"학교 행정표준코드"
//	@Param			year	query	string	true	"요청하려는 연도 (YYYY)"
//	@Param			month	query	string	true	"요청하려는 달 (MM)"
//	@Success 200 {object} controller.CommonResponse{data=[]service.Meal}
//	@Router			/api/v1/schools/meals [get]
func (ctl *schoolContoller) ListMeals(c *fiber.Ctx) error {
	eduOfficeCode := c.Query("eduOfficeCode")
	schoolCode := c.Query("schoolCode")
	year := c.Query(("year"))
	month := c.Query(("month"))

	meals, err := ctl.svc.GetMealPlans(eduOfficeCode, schoolCode, year, month)
	if err != nil {
		log.Println(err)
		return ErrorOf(fiber.StatusInternalServerError, "50000", MessageOfCode(fiber.StatusInternalServerError))
	}
	return UTF8Json(c, ResponseOfCode(false, fiber.StatusOK, meals))
}
