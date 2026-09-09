import http from 'k6/http';
import { check } from 'k6';

// Staircase: warm-up, then seven distinct RPM plateaus from 100k up to
// 400k, each held long enough (30s, several Prometheus scrape intervals)
// to show as a flat step on the Grafana "Request rate" graph rather than
// blending into a smooth ramp. Targets the gateway (nginx on :8000), not
// the app directly, so numbers reflect the whole system (gateway
// rate-limit + Redis cache-aside), not a bare Go handler.
//
// RPM -> RPS (rounded): 100k/150k/200k/250k/300k/350k/400k RPM ==
// 1667/2500/3334/4167/5000/5834/6667 RPS.
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';
const WARMUP_RPS = 300; // ~18k RPM -- just enough to fill Redis/DB caches and connection pools, not a load test in itself
const STEPS_RPS = [1667, 2500, 3334, 4167, 5000, 5834, 6667]; // 100k..400k RPM
const RAMP = '8s'; // between-step ramp, short relative to the 30s hold so each plateau still reads as flat
const HOLD = '30s';

function stepStages() {
  const stages = [];
  for (const rps of STEPS_RPS) {
    stages.push({ target: rps, duration: RAMP });
    stages.push({ target: rps, duration: HOLD });
  }
  return stages;
}

export const options = {
  scenarios: {
    highload: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 1000,
      maxVUs: 3000,
      stages: [
        { target: WARMUP_RPS, duration: '5s' }, // ramp into warm-up
        { target: WARMUP_RPS, duration: '15s' }, // 20s total warm-up: fills catalog cache, primes DB/Redis connection pools
        ...stepStages(),
        { target: 0, duration: '10s' },
      ],
    },
  },
  thresholds: {
    // Per-endpoint thresholds via k6's tag-scoped metrics -- these don't
    // fail the run (no abortOnFail), they're just asserted so the summary
    // clearly flags which endpoint (if any) missed the mark.
    'http_req_failed{endpoint:catalog_categories}': ['rate<0.01'],
    'http_req_failed{endpoint:catalog_services}': ['rate<0.01'],
    'http_req_failed{endpoint:health}': ['rate<0.01'],
    'http_req_failed{endpoint:master_profile}': ['rate<0.01'],
    'http_req_failed{endpoint:my_requests}': ['rate<0.01'],
    'http_req_duration{endpoint:catalog_categories}': ['p(95)<200'],
  },
};

// Runs once before the load starts, result is shared read-only across all
// VUs -- one login, one token reused everywhere, not one login per VU/iter.
export function setup() {
  const res = http.post(
    `${BASE_URL}/api/auth/login`,
    JSON.stringify({ username: 'master1', password: 'Demo12345' }),
    { headers: { 'Content-Type': 'application/json' } },
  );
  if (res.status !== 200) {
    throw new Error(`setup login failed: ${res.status} ${res.body}`);
  }
  return { token: res.json('accessToken') };
}

// Weighted mix: catalog reads (cached, the actual 100k-RPM target) dominate,
// health is light background traffic, authorized reads are a smaller slice
// -- confirms the gateway+cache path under load without hammering
// Postgres-backed write paths, which were deliberately left out.
export default function (data) {
  const roll = Math.random();
  const authHeaders = { headers: { Authorization: `Bearer ${data.token}` } };

  if (roll < 0.55) {
    const res = http.get(`${BASE_URL}/api/catalog/categories`, {
      tags: { endpoint: 'catalog_categories' },
    });
    check(res, { 'categories 200': (r) => r.status === 200 });
  } else if (roll < 0.75) {
    const res = http.get(`${BASE_URL}/api/catalog/services`, {
      tags: { endpoint: 'catalog_services' },
    });
    check(res, { 'services 200': (r) => r.status === 200 });
  } else if (roll < 0.9) {
    const res = http.get(`${BASE_URL}/health/live`, {
      tags: { endpoint: 'health' },
    });
    check(res, { 'health 200': (r) => r.status === 200 });
  } else if (roll < 0.95) {
    const res = http.get(`${BASE_URL}/api/master/profile`, {
      ...authHeaders,
      tags: { endpoint: 'master_profile' },
    });
    check(res, { 'master profile 200': (r) => r.status === 200 });
  } else {
    const res = http.get(`${BASE_URL}/api/requests`, {
      ...authHeaders,
      tags: { endpoint: 'my_requests' },
    });
    check(res, { 'my requests 200': (r) => r.status === 200 });
  }
}
