package network

import (
	"0tak2/afterhee-server/configuration"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const neisBaseURL = "https://open.neis.go.kr/hub"

type NEISMealRequest interface {
	// eduOfficeCode 시도교육청코드
	// schoolAdminCode 학교 행정표준코드
	// year 요청 연도
	// month 요청 월
	FetchMealPlan(eduOfficeCode string, schoolAdminCode string, year string, month string) (*MealServiceDietInfoResponse, error)
}

type neisMealRequest struct {
}

func NewNEISMealRequest() NEISMealRequest {
	return &neisMealRequest{}
}

func (r neisMealRequest) FetchMealPlan(eduOfficeCode string, schoolAdminCode string, year string, month string) (*MealServiceDietInfoResponse, error) {
	config := configuration.GetConfiguration()

	destination, _ := url.Parse(neisBaseURL + "/mealServiceDietInfo")

	params := url.Values{}
	params.Add("key", config.NEISAPIKey)
	params.Add("Type", "json")
	params.Add("ATPT_OFCDC_SC_CODE", eduOfficeCode)
	params.Add("SD_SCHUL_CODE", schoolAdminCode)
	params.Add("MLSV_YMD", year+month)

	destination.RawQuery = params.Encode()

	log.Println("Request -> " + destination.String())

	resp, err := http.Get(destination.String())

	if err != nil {
		log.Println("error requesting NEIS meal plane... " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status: %d\n", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("Response <-", string(bodyBytes))

	var result MealServiceDietInfoResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
