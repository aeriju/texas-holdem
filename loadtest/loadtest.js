/**
 * Holdem Load Test - k6
 * Tests REST endpoints under load.
 *
 * Run: k6 run loadtest/loadtest.js
 * With backend: BASE_URL=http://localhost:8080 k6 run loadtest/loadtest.js
 * Options: k6 run --vus 50 --duration 30s loadtest/loadtest.js
 */
import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 50 },
    { duration: '10s', target: 100 },
    { duration: '20s', target: 50 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<600'],
    checks: ['rate>0.99'],
  },
};

export default function () {
  const roll = Math.random();

  if (roll < 0.34) {
    const payload = JSON.stringify({
      hole: ['HA', 'HK'],
      community: ['C2', 'D3', 'S4', 'H5', 'C6'],
    });
    const res = http.post(`${BASE_URL}/api/v1/best-hand`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
      'best-hand status 200': (r) => r.status === 200,
    });
  } else if (roll < 0.67) {
    const payload = JSON.stringify({
      hand1: {
        hole: ['HA', 'HK'],
        community: ['C2', 'D3', 'S4', 'H5', 'C6'],
      },
      hand2: {
        hole: ['S9', 'S8'],
        community: ['C2', 'D3', 'S4', 'H5', 'C6'],
      },
    });
    const res = http.post(`${BASE_URL}/api/v1/heads-up`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
      'heads-up status 200': (r) => r.status === 200,
    });
  } else {
    const payload = JSON.stringify({
      hole: ['HA', 'HK'],
      community: ['C2', 'D3', 'S4'],
      players: 6,
      simulations: 2000,
    });
    const res = http.post(`${BASE_URL}/api/v1/odds`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
      'odds status 200': (r) => r.status === 200,
    });
  }

  sleep(0.3);
}
