package service

import (
	"0tak2/afterhee-server/network"
)

// Domain Entity
type MenuRecommendCommand struct {
	Category  string
	SkipMenus []string
}

type MenuRecommendation struct {
	Menu   string `json:"recommendation"`
	Reason string `json:"reason"`
}

// Service
type GeminiService interface {
	GetMenuRecommendation(command MenuRecommendCommand) (*MenuRecommendation, error)
	GetMenuRecommendationWithCustomPrompt(prompt string) (*MenuRecommendation, error)
}

type geminiService struct {
	gemini network.GeminiRequest
}

func NewGeminiService(gemini network.GeminiRequest) GeminiService {
	return &geminiService{
		gemini: gemini,
	}
}

func (s geminiService) GetMenuRecommendation(command MenuRecommendCommand) (*MenuRecommendation, error) {
	result, err := s.gemini.RecommendMealMenu(command.SkipMenus, command.Category)

	if err != nil {
		return nil, err
	}

	return &MenuRecommendation{
		Menu:   result.Menu,
		Reason: result.Reason,
	}, nil
}

func (s geminiService) GetMenuRecommendationWithCustomPrompt(prompt string) (*MenuRecommendation, error) {
	result, err := s.gemini.PredictManually(prompt)

	if err != nil {
		return nil, err
	}

	return &MenuRecommendation{
		Menu:   result.Menu,
		Reason: result.Reason,
	}, nil
}
