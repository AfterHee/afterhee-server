package network

import (
	"0tak2/afterhee-server/configuration"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"text/template"

	"google.golang.org/genai"
)

type GeminiRequest interface {
	RecommendMealMenu(skipMenus []string, category string) (*GeminiMenuRecommendation, error)
	PredictManually(prompt string) (*GeminiMenuRecommendation, error)
}

type geminiRequest struct {
}

func NewGeminiRequest() GeminiRequest {
	return &geminiRequest{}
}

func (r geminiRequest) RecommendMealMenu(skipMenus []string, category string) (*GeminiMenuRecommendation, error) {
	prompt, err := CreatePrompt(GeminiMenuRecommendationPromptValues{
		Category:  category,
		SkipMenus: skipMenus,
	})

	if err != nil {
		log.Println("failed to create prompt using template... " + err.Error())
		return nil, err
	}

	return r.requestToGemini(*prompt)
}

func (r geminiRequest) PredictManually(prompt string) (*GeminiMenuRecommendation, error) {
	return r.requestToGemini(prompt)
}

func (r geminiRequest) requestToGemini(prompt string) (*GeminiMenuRecommendation, error) {
	serverConfig := configuration.GetConfiguration()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: serverConfig.GeminiAPIKey,
	})
	if err != nil {
		log.Println("failed to create a gemini client..." + err.Error())
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"recommendation": {Type: genai.TypeString},
				"reason":         {Type: genai.TypeString},
			},
			PropertyOrdering: []string{"recommendation", "reason"},
		},
	}

	log.Println("will send prediction request to gemini with: " + prompt)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		log.Println("failed to request to gemini..." + err.Error())
		return nil, err
	}

	var rawResponse string = result.Text()
	log.Println("received response from gemini... response: " + rawResponse)

	var response GeminiMenuRecommendation
	json.Unmarshal([]byte(rawResponse), &response)

	return &response, nil
}

func preprocessLLMInputs(values GeminiMenuRecommendationPromptValues) GeminiMenuRecommendationPromptValues {
	switch values.Category {
	case "아시안":
		values.Category = "인도 음식, 태국 음식, 베트남 음식, 네팔 음식, 몽골 음식"
	case "세계음식":
		values.Category = "터키 등 중동 음식, 멕세코 등 남미 음식, 샤실리크 등 러시아 음식, 그 외 아시안 음식을 제외한 세계 음식"
	case "볶음":
		values.Category = "볶음 요리"
	case "밥":
		values.Category = "덮밥, 비빔밥, 볶음밥 등 각종 밥 요리"
	case "면":
		values.Category = "면 요리"
	case "야채":
		values.Category = "샐러드, 포케 등 야채가 풍부하게 들어간 음식"
	case "빠른 식사":
		values.Category = "분식, 햄버거, 샌드위치 등 빠르게 주문해서 먹을 수 있는 음식. 컵밥은 제외."
	}

	return values
}

func CreatePrompt(values GeminiMenuRecommendationPromptValues) (*string, error) {
	processedValues := preprocessLLMInputs(values)

	tmpl := `Recommend one new menu item that falls under the following category and is not included in the exclusion menu, along with the reason. Please respond in Korean.

Category: {{ .Category }}
Exclusion Menu: {{range .SkipMenus}}{{.}}, {{end}}
`

	t := template.Must(template.New("MenuRecommendRequestPrompt").Parse(tmpl))

	var buf bytes.Buffer
	if err := t.Execute(&buf, processedValues); err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}

type GeminiMenuRecommendationPromptValues struct {
	Category  string
	SkipMenus []string
}
