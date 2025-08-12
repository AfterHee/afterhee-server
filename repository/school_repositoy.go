package repository

import (
	"database/sql"
	"time"

	_ "github.com/marcboeker/go-duckdb/v2"
)

// Data Model
type SchoolModel struct {
	Id                        int64
	SidoEduOfficeCode         *string
	SidoEduOfficeName         *string
	AdminStandardCode         *string
	SchoolName                *string
	SchoolNameEng             *string
	SchoolTypeName            *string
	SidoName                  *string
	JurisdictionOrgName       *string
	EstablishmentName         *string
	RoadZipCode               *string
	RoadAddress               *string
	RoadDetailAddress         *string
	PhoneNumber               *string
	HomepageURL               *string
	GenderMixTypeName         *string
	FaxNumber                 *string
	HighSchoolTypeName        *string
	IndustrySpecialClassYN    *string
	HighSchoolGeneralSpecName *string
	SpecialPurposeHSCategory  *string
	EntranceTermType          *string
	DayNightType              *string
	EstablishedAt             *time.Time
	FoundedAt                 *time.Time
	ModifiedAt                *time.Time
}

// DAO
type SchoolRepository interface {
	FetchSchools(keyword string) ([]SchoolModel, error)
}

type schoolRepository struct {
	db *sql.DB
}

func NewSchoolRepository(db *sql.DB) SchoolRepository {
	return &schoolRepository{db: db}
}

func (r schoolRepository) FetchSchools(keyword string) ([]SchoolModel, error) {
	rows, err := r.db.Query(`
		SELECT Id, SidoEduOfficeCode, SidoEduOfficeName, SchoolName, AdminStandardCode, RoadAddress, RoadDetailAddress
		FROM schools
		WHERE SchoolName LIKE ?`,
		"%"+keyword+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []SchoolModel
	for rows.Next() {
		var s SchoolModel
		if err := rows.Scan(&s.Id, &s.SidoEduOfficeCode, &s.SidoEduOfficeName, &s.SchoolName, &s.AdminStandardCode, &s.RoadAddress, &s.RoadDetailAddress); err != nil {
			return nil, err
		}
		schools = append(schools, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schools, nil
}
