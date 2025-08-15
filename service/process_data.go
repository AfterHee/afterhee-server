package service

import (
	"html"
	"regexp"
	"strings"
)

func NilToEmptyString(nilOrString *string) string {
	if nilOrString != nil {
		return *nilOrString
	}

	return ""
}

func ExtractMenus(ddish string) []string {
	var brSplitter = regexp.MustCompile(`(?i)<br\s*/?>`)
	var trailingAllergen = regexp.MustCompile(`\s*\(\d{1,2}(?:\.\d{1,2})*\)\s*$`)

	if ddish == "" {
		return nil
	}

	lines := brSplitter.Split(ddish, -1)
	out := make([]string, 0, len(lines))

	for _, raw := range lines {
		// HTML 엔티티 복원 및 공백 정리
		s := html.UnescapeString(strings.TrimSpace(raw))
		s = strings.ReplaceAll(s, "\u00A0", " ") // nbsp → space
		s = strings.Join(strings.Fields(s), " ") // 다중 공백 압축

		// 맨 끝의 알레르기 번호 괄호를 반복적으로 제거
		for {
			before := s
			s = trailingAllergen.ReplaceAllString(s, "")
			if s == before {
				break
			}
		}

		if s != "" {
			out = append(out, strings.TrimSuffix(s, "."))
		}
	}
	return out
}
