package service

import (
	"0tak2/afterhee-server/network"
	"0tak2/afterhee-server/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

// Domain Entity
type School struct {
	Id                int64  `json:"id"`
	SidoEduOfficeCode string `json:"sidoEduOfficeCode"`
	SidoEduOfficeName string `json:"sidoEduOfficeName"`
	AdminStandardCode string `json:"adminStandardCode"`
	SchoolName        string `json:"schoolName"`
	RoadAddress       string `json:"roadAddress"`
	RoadDetailAddress string `json:"roadDetailAddress"`
}

type Meal struct {
	SidoEduOfficeCode string   `json:"sidoEduOfficeCode"`
	SidoEduOfficeName string   `json:"sidoEduOfficeName"`
	SchoolAdminCode   string   `json:"schoolAdminCode"`
	SchoolName        string   `json:"schoolName"`
	MmealScCode       string   `json:"mmealScCode"`
	MmealScNm         string   `json:"mmealScNm"`
	MlsvYmd           string   `json:"mlsvYmd"`
	MlsvFgr           float32  `json:"mlsvFgr"`
	DishName          []string `json:"dishNames"`
	OrplcInfo         string   `json:"orplcInfo"`
	CalInfo           string   `json:"calInfo"`
	NtrInfo           string   `json:"ntrInfo"`
	MlsvFromYmd       string   `json:"mlsvFromYmd"`
	MlsvToYmd         string   `json:"mlsvToYmd"`
	LoadDtm           string   `json:"loadDtm"`
}

// Service
type SchoolService interface {
	GetSchools(keyword string) ([]School, error)
	GetMealPlans(ctx context.Context, sidoEduOfficeCode string, schoolStandardCode string, from time.Time, to time.Time) ([]Meal, error)
}

type schoolService struct {
	repo  repository.SchoolRepository
	cache repository.CacheRepository
	neis  network.NEISMealRequest
}

func NewSchoolService(repo repository.SchoolRepository, cache repository.CacheRepository, neis network.NEISMealRequest) SchoolService {
	return &schoolService{
		repo:  repo,
		cache: cache,
		neis:  neis,
	}
}

func (s schoolService) GetSchools(keyword string) ([]School, error) {
	result, err := s.repo.FetchSchools(keyword)

	if err != nil {
		return nil, err
	}

	var schools []School
	for _, schoolModel := range result {
		schools = append(schools, School{
			Id:                schoolModel.Id,
			SidoEduOfficeCode: NilToEmptyString(schoolModel.SidoEduOfficeCode),
			SidoEduOfficeName: NilToEmptyString(schoolModel.SidoEduOfficeName),
			AdminStandardCode: NilToEmptyString(schoolModel.AdminStandardCode),
			SchoolName:        NilToEmptyString(schoolModel.SchoolName),
			RoadAddress:       NilToEmptyString(schoolModel.RoadAddress),
			RoadDetailAddress: NilToEmptyString(schoolModel.RoadDetailAddress),
		})
	}

	sort.Slice(schools, func(i, j int) bool {
		return schools[i].Id < schools[j].Id
	})

	return schools, nil
}

func (s schoolService) GetMealPlans(ctx context.Context, sidoEduOfficeCode string, schoolStandardCode string, from time.Time, to time.Time) ([]Meal, error) {
	// Fetch cached data
	cacheKey := fmt.Sprintf("%s_%s_%d_%d", sidoEduOfficeCode, schoolStandardCode, from.Unix(), to.Unix())
	cachedValue, cacheFetchErr := s.cache.GetValue(ctx, cacheKey)

	if cacheFetchErr == nil && cachedValue != nil {
		var cachedData []Meal
		err := json.Unmarshal([]byte(*cachedValue), &cachedData)
		if err == nil {
			return cachedData, nil
		}
		log.Printf("failed to unmarshal cached data: %v", err)
	}

	// Fetch server data
	result, err := s.neis.FetchMealPlan(sidoEduOfficeCode, schoolStandardCode, timeToString(from), timeToString((to)))
	if err != nil {
		return nil, err
	}

	if len(result.MealServiceDietInfo) < 2 {
		return nil, errors.New("no meal rows") // TODO: Custom Error
	}

	rows := result.MealServiceDietInfo[1].Row

	var meals []Meal
	for _, row := range rows {
		log.Println(row)
		meals = append(meals, Meal{
			SidoEduOfficeCode: row.AtptOfcdcScCode,
			SidoEduOfficeName: row.AtptOfcdcScNm,
			SchoolAdminCode:   row.SdSchulCode,
			SchoolName:        row.SchulNm,
			MmealScCode:       row.MmealScCode,
			MmealScNm:         row.MmealScNm,
			MlsvYmd:           row.MlsvYmd,
			MlsvFgr:           row.MlsvFgr,
			DishName:          ExtractMenus(row.DdishNm),
			OrplcInfo:         row.OrplcInfo,
			CalInfo:           row.CalInfo,
			NtrInfo:           row.NtrInfo,
			MlsvFromYmd:       row.MlsvFromYmd,
			MlsvToYmd:         row.MlsvToYmd,
			LoadDtm:           row.LoadDtm,
		})
	}

	sort.Slice(meals, func(i, j int) bool {
		return meals[i].MlsvYmd < meals[j].MlsvYmd
	})

	jsonBytes, marshalErr := json.Marshal(meals)
	if marshalErr != nil {
		log.Printf("failed to marshal data that will be planned to cache. skip cache it. error=%s", marshalErr.Error())
	} else {
		ttl := 3 * 24 * time.Hour
		s.cache.SetValue(ctx, cacheKey, string(jsonBytes), ttl)
	}

	return meals, nil
}

func timeToString(time time.Time) string {
	return time.Format("20060102")
}
