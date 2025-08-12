package service

import (
	"0tak2/afterhee-server/network"
	"0tak2/afterhee-server/repository"
	"errors"
	"log"
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
	SidoEduOfficeCode string  `json:"sidoEduOfficeCode"`
	SidoEduOfficeName string  `json:"sidoEduOfficeName"`
	SchoolAdminCode   string  `json:"schoolAdminCode"`
	SchoolName        string  `json:"schoolName"`
	MmealScCode       string  `json:"mmealScCode"`
	MmealScNm         string  `json:"mmealScNm"`
	MlsvYmd           string  `json:"mlsvYmd"`
	MlsvFgr           float32 `json:"mlsvFgr"`
	DishName          string  `json:"dishName"`
	OrplcInfo         string  `json:"orplcInfo"`
	CalInfo           string  `json:"calInfo"`
	NtrInfo           string  `json:"ntrInfo"`
	MlsvFromYmd       string  `json:"mlsvFromYmd"`
	MlsvToYmd         string  `json:"mlsvToYmd"`
	LoadDtm           string  `json:"loadDtm"`
}

// Service
type SchoolService interface {
	GetSchools(keyword string) ([]School, error)
	GetMealPlans(sidoEduOfficeCode string, schoolStandardCode string, year string, month string) ([]Meal, error)
}

type schoolService struct {
	repo repository.SchoolRepository
	neis network.NEISMealRequest
}

func NewSchoolService(repo repository.SchoolRepository, neis network.NEISMealRequest) SchoolService {
	return &schoolService{
		repo: repo,
		neis: neis,
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
			SidoEduOfficeCode: nilToEmptyString(schoolModel.SidoEduOfficeCode),
			SidoEduOfficeName: nilToEmptyString(schoolModel.SidoEduOfficeName),
			AdminStandardCode: nilToEmptyString(schoolModel.AdminStandardCode),
			SchoolName:        nilToEmptyString(schoolModel.SchoolName),
			RoadAddress:       nilToEmptyString(schoolModel.RoadAddress),
			RoadDetailAddress: nilToEmptyString(schoolModel.RoadDetailAddress),
		})
	}

	return schools, err
}

func (s schoolService) GetMealPlans(sidoEduOfficeCode string, schoolStandardCode string, year string, month string) ([]Meal, error) {
	result, err := s.neis.FetchMealPlan(sidoEduOfficeCode, schoolStandardCode, year, month)
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
			DishName:          row.DdishNm,
			OrplcInfo:         row.OrplcInfo,
			CalInfo:           row.CalInfo,
			NtrInfo:           row.NtrInfo,
			MlsvFromYmd:       row.MlsvFromYmd,
			MlsvToYmd:         row.MlsvToYmd,
			LoadDtm:           row.LoadDtm,
		})
	}
	return meals, err
}

func nilToEmptyString(nilOrString *string) string {
	if nilOrString != nil {
		return *nilOrString
	}

	return ""
}
