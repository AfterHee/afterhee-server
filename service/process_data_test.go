package service_test

import (
	"0tak2/afterhee-server/service"
	"testing"
)

func TestProcessMealData(t *testing.T) {
	rawData := "추가밥. <br/>탄탄멘. (1.2.4.5.6.10.13.15.16)<br/>오꼬노미야끼. (1.5.6.9.10.13.17.18)<br/>고구마파이. (1.5.6)<br/>우유. (2)<br/>락교. (13)"
	result := service.ExtractMenus(rawData)

	if !equalSlice(
		result,
		[]string{
			"추가밥",
			"탄탄멘",
			"오꼬노미야끼",
			"고구마파이",
			"우유",
			"락교",
		},
	) {
		t.Error("기댓값과 결과가 다릅니다.")
	}
}

func TestNilToEmptyString(t *testing.T) {
	str1 := "Nil이 아닌 스트링"
	str1ptr := &str1
	result1 := service.NilToEmptyString(str1ptr)
	expected1 := "Nil이 아닌 스트링"

	if result1 != expected1 {
		t.Error("기댓값과 결과가 다릅니다. expected: " + expected1 + " result: " + result1)
	}

	result2 := service.NilToEmptyString(nil)
	expected2 := ""

	if result2 != expected2 {
		t.Error("기댓값과 결과가 다릅니다. expected: " + expected2 + " result: " + result2)
	}
}

func equalSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
