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
	var response GeminiMenuRecommendation
	json.Unmarshal([]byte(rawResponse), &response)

	return &response, nil
}

func CreatePrompt(values GeminiMenuRecommendationPromptValues) (*string, error) {
	tmpl := `Recommend one new menu item that falls under the following category and is not included in the exclusion menu, along with the reason. Please respond in Korean.

Category: {{ .Category }}
Exclusion Menu: {{range .SkipMenus}}{{.}}, {{end}}
`

	t := template.Must(template.New("MenuRecommendRequestPrompt").Parse(tmpl))

	var buf bytes.Buffer
	if err := t.Execute(&buf, values); err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}

type GeminiMenuRecommendationPromptValues struct {
	Category  string
	SkipMenus []string
}
