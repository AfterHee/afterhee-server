# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

- 식단 조회 캐싱 시, TTL을 기존 3일에서 1일로 변경

### Deprecated
                                                  
### Removed

### Fixed

- NEIS에서 반환하는 식단 정보가 빈 슬라이스여도 Redis에 저장

### Security

## [0.0.7] - 2026-02-07

### Fixed

- NEIS에서 반환하는 식단 정보가 완전히 없으면 오류를 반환하던 문제

## [0.0.6] - 2026-02-07

- 식단 조회 캐싱

## [0.0.5] - 2025-09-11

### Added

- HealthCheck 엔드포인트 추가

### Fixed

- Fiber 기본 에러를 전역 에러 핸들러가 처리하지 못하는 문제 수정

## [0.0.4] - 2025-09-07

### Fixed

- "해산물" 추천이 안되던 문제 수정 (#26)

## [0.0.3] - 2025-08-26

### Fixed

- SuggestController 입력 값 검증 (#22)

## [0.0.2] - 2025-08-23

### Changed

- 카테고리 픽스 (#16)
- 내부 프롬프팅 개선

## [0.0.1] - 2025-08-18

### Added

- 기본 엔드포인트 구현

[unreleased]: https://github.com/AfterHee/afterhee-server/compare/v0.0.6...HEAD
[0.0.6]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.6
[0.0.5]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.5
[0.0.4]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.4
[0.0.3]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.3
[0.0.2]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.2
[0.0.1]: https://github.com/AfterHee/afterhee-server/releases/tag/v0.0.1
