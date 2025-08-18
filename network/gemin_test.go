package network_test

import (
	"0tak2/afterhee-server/network"
	"testing"
)

func TestProcessMealData(t *testing.T) {
	values := network.GeminiMenuRecommendationPromptValues{
		Category:  "중식",
		SkipMenus: []string{"김치찌개", "짜장면", "된장국"},
	}

	prompt, err := network.CreatePrompt(values)

	if err != nil {
		t.Error("예외가 발생했습니다... " + err.Error())
	}

	if *prompt != `Recommend one new menu item that falls under the following category and is not included in the exclusion menu, along with the reason. Please respond in Korean.

Category: 중식
Exclusion Menu: 김치찌개, 짜장면, 된장국, 
` {
		print(*prompt)
		t.Error("기댓값과 결과가 다릅니다.")
	}
}
