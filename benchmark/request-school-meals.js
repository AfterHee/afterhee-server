import { check } from 'k6';
import http from 'k6/http';

/**
 * 테스트 옵션
 */
const testerOptions = {
    testServerIpAndPort: 'localhost:8080',
    eduOfficeCode: 'B10',
    schoolCode: '7010266',
    startDate: '2026-01-01',
    endDate: '2026-01-31',
}

export const options = {
  hosts: { 'test.afterhee.server': testerOptions.testServerIpAndPort },
  stages: [
    { duration: '10s', target: 10 },
    { duration: '50s', target: 10 },
  ],
  userAgent: 'AfterHeeStessTest/1.0',
};

export default function () {
  const { eduOfficeCode, schoolCode, startDate, endDate } = testerOptions
  const res = http.get(`http://test.afterhee.server/api/v1/schools/meals?eduOfficeCode=${eduOfficeCode}&schoolCode=${schoolCode}&from=${startDate}&to=${endDate}`);
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}
