package controller

import (
	"0tak2/afterhee-server/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Request Body DTO
type SuggestMenuRequest struct {
	Category  string
	SkipMenus []string
}

type SuggestTestRequest struct {
	Prompt string
}

// Controller Types
type SuggestContoller interface {
	Suggest(c *fiber.Ctx) error
	SuggestTest(c *fiber.Ctx) error
}

type suggestContoller struct {
	svc service.GeminiService
}

func NewSuggestController(svc service.GeminiService) SuggestContoller {
	return &suggestContoller{svc: svc}
}

// Suggest godoc
//
//	@Summary		메뉴 제안
//	@Description	메뉴를 제안한다. category는 다음만 허용한다: "한식", "일식", "중식", "양식", "아시안", "세계음식", "찜", "국물", "볶음", "밥", "면", "빵", "해물", "고기", "야채", "빠른 식사", "디저트",
//	@Accept			json
//	@Produce		json
//	@Param request body SuggestMenuRequest true "request suggest menu body"
//	@Success 200 {object} controller.CommonResponse{data=controller.SuggestMenuRequest}
//	@Router			/api/v1/suggest [post]
func (ctl *suggestContoller) Suggest(c *fiber.Ctx) error {
	var request SuggestMenuRequest

	if err := c.BodyParser(&request); err != nil {
		log.Println("failed to parse request body... " + err.Error())
		return err
	}

	// Validate Request
	if !validateCategory(request.Category) {
		log.Println("not pre-defined category... input=" + request.Category)
		return ErrorOf(fiber.StatusBadRequest, "40000", MessageOfCode(fiber.StatusBadRequest))
	}

	if !validateSkipMenus(request.SkipMenus) {
		log.Println("so long. max length is 30...")
		return ErrorOf(fiber.StatusBadRequest, "40000", MessageOfCode(fiber.StatusBadRequest))
	}

	menu, err := ctl.svc.GetMenuRecommendation(service.MenuRecommendCommand{
		Category:  request.Category,
		SkipMenus: request.SkipMenus,
	})

	if err != nil {
		log.Println(err)
		return ErrorOf(fiber.StatusInternalServerError, "50000", MessageOfCode(fiber.StatusInternalServerError))
	}

	return UTF8Json(c, ResponseOfCode(false, fiber.StatusOK, menu))
}

// SuggestTest godoc
//
//	@Summary		매뉴얼 프롬프트 테스팅
//	@Description	프롬프트로 LLM에 요청을 보낸다.
//	@Accept			json
//	@Produce		json
//	@Param request body SuggestTestRequest true "request suggest menu body"
//	@Success 200 {object} controller.CommonResponse{data=controller.SuggestMenuRequest}
//	@Router			/api/v1/test/suggest [post]
func (ctl *suggestContoller) SuggestTest(c *fiber.Ctx) error {
	var request SuggestTestRequest

	if err := c.BodyParser(&request); err != nil {
		log.Println("failed to parse request body... " + err.Error())
		return err
	}

	// Validate Request
	if request.Prompt == "" {
		log.Println("prompt is empty...")
		return ErrorOf(fiber.StatusBadRequest, "40000", MessageOfCode(fiber.StatusBadRequest))
	}

	menu, err := ctl.svc.GetMenuRecommendationWithCustomPrompt(request.Prompt)

	if err != nil {
		log.Println(err)
		return ErrorOf(fiber.StatusInternalServerError, "50000", MessageOfCode(fiber.StatusInternalServerError))
	}

	return UTF8Json(c, ResponseOfCode(false, fiber.StatusOK, menu))
}

// 카테고리가 정의된 카테고리인지 확인한다.
func validateCategory(rawCategory string) bool {
	allowCategories := []string{
		"한식", "일식", "중식", "양식", "아시안",
		"세계음식", "찜", "국물", "볶음", "밥",
		"면", "빵", "해물", "고기", "야채",
		"빠른 식사", "디저트",
	}

	for _, allow := range allowCategories {
		if allow == rawCategory {
			return true
		}
	}

	return false
}

// 순회하면서 각 아이템의 글자 수를 확인한다. 30자 초과이면 검증에 실패한다.
func validateSkipMenus(rawSkipMenus []string) bool {
	maxLength := 30

	for _, inputMenu := range rawSkipMenus {
		if len([]rune(inputMenu)) > maxLength {
			return false
		}
	}

	return true
}
