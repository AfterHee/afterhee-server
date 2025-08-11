package service

import "0tak2/afterhee-server/repository"

// Domain Entity
type School struct {
	Id                int64  `json:"id"`
	AdminStandardCode string `json:"adminStandardCode"`
	SchoolName        string `json:"schoolName"`
	RoadAddress       string `json:"roadAddress"`
	RoadDetailAddress string `json:"roadDetailAddress"`
}

// Service
type SchoolService interface {
	GetSchools(keyword string) ([]School, error)
}

type schoolService struct {
	repo repository.SchoolRepository
}

func NewSchoolService(repo repository.SchoolRepository) SchoolService {
	return &schoolService{repo: repo}
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
			AdminStandardCode: nilToEmptyString(schoolModel.AdminStandardCode),
			SchoolName:        nilToEmptyString(schoolModel.SchoolName),
			RoadAddress:       nilToEmptyString(schoolModel.RoadAddress),
			RoadDetailAddress: nilToEmptyString(schoolModel.RoadDetailAddress),
		})
	}

	return schools, err
}

func nilToEmptyString(nilOrString *string) string {
	if nilOrString != nil {
		return *nilOrString
	}

	return ""
}
