package service

import (
	"0tak2/afterhee-server/network"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 포인터 값 할당을 위한 헬퍼 함수
func intPtr(i int) *int {
	return &i
}

// --- Mocks ---

type mockNEISRequest struct {
	mock.Mock
}

func (m *mockNEISRequest) FetchMealPlan(eduOfficeCode string, schoolAdminCode string, fromDateString string, toDateString string) (*network.MealServiceDietInfoResponse, error) {
	args := m.Called(eduOfficeCode, schoolAdminCode, fromDateString, toDateString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*network.MealServiceDietInfoResponse), args.Error(1)
}

type extendedMockCache struct {
	mock.Mock
}

func (m *extendedMockCache) GetValue(ctx context.Context, key string) (*string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *extendedMockCache) SetValue(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// --- Test Suite ---

func TestSchoolService_GetMealPlans(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 10, 1, 23, 59, 59, 0, time.UTC)

	sidoCode := "B10"
	schoolCode := "7010537"

	t.Run("캐시에 데이터가 존재할 경우 네트워킹 없이 즉시 반환한다", func(t *testing.T) {
		mockCache := new(extendedMockCache)
		mockNEIS := new(mockNEISRequest)
		service := NewSchoolService(nil, mockCache, mockNEIS)

		cachedJSON := `[{"schoolName": "테스트고등학교", "mlsvYmd": "20231001", "dishNames": ["밥", "국"]}]`
		mockCache.On("GetValue", ctx, mock.Anything).Return(&cachedJSON, nil)

		result, err := service.GetMealPlans(ctx, sidoCode, schoolCode, from, to)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "테스트고등학교", result[0].SchoolName)
		mockNEIS.AssertNotCalled(t, "FetchMealPlan")
	})

	t.Run("캐시가 없고 네트워크 성공 시, 데이터를 파싱하고 캐시에 저장한다", func(t *testing.T) {
		mockCache := new(extendedMockCache)
		mockNEIS := new(mockNEISRequest)
		service := NewSchoolService(nil, mockCache, mockNEIS)

		mockCache.On("GetValue", ctx, mock.Anything).Return(nil, nil)

		// 수정된 부분: intPtr 사용 및 필드 구조 매칭
		mockResponse := &network.MealServiceDietInfoResponse{
			MealServiceDietInfo: []network.MealServiceDietInfo{
				{
					Head: []network.HeadInfo{
						{ListTotalCount: intPtr(1)},
					},
				},
				{
					Row: []network.MealRow{
						{
							SchulNm: "네트워크고교",
							MlsvYmd: "20231001",
							DdishNm: "돈까스<br/>김치",
						},
					},
				},
			},
		}
		mockNEIS.On("FetchMealPlan", sidoCode, schoolCode, "20231001", "20231001").Return(mockResponse, nil)
		mockCache.On("SetValue", ctx, mock.Anything, mock.Anything, 3*24*time.Hour).Return(nil)

		result, err := service.GetMealPlans(ctx, sidoCode, schoolCode, from, to)

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, "네트워크고교", result[0].SchoolName)
		mockCache.AssertCalled(t, "SetValue", ctx, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("급식 정보가 없는 경우 에러를 반환한다", func(t *testing.T) {
		mockCache := new(extendedMockCache)
		mockNEIS := new(mockNEISRequest)
		service := NewSchoolService(nil, mockCache, mockNEIS)

		mockCache.On("GetValue", ctx, mock.Anything).Return(nil, nil)

		emptyResponse := &network.MealServiceDietInfoResponse{
			MealServiceDietInfo: []network.MealServiceDietInfo{
				{
					Head: []network.HeadInfo{
						{ListTotalCount: intPtr(0)},
					},
				},
			},
		}
		mockNEIS.On("FetchMealPlan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(emptyResponse, nil)

		result, err := service.GetMealPlans(ctx, sidoCode, schoolCode, from, to)

		assert.Error(t, err)
		assert.Equal(t, "no meal rows", err.Error())
		assert.Nil(t, result)
	})
}
