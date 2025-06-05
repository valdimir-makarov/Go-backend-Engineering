import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  vus: 15,         // 15 concurrent users
  duration: '30s', // Test duration of 30 seconds
  rps: 150,        // Target of 150 requests/second
  thresholds: {
    http_req_duration: ['p(90)<500'], // 90% of requests below 500ms
  },
};

export default function () {
  http.post('http://localhost:8080/execute', JSON.stringify({
    language: 'python',
    code: 'print("Hello, World!")',
    method: 'docker'
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  sleep(0.2); // Small delay to keep it steady and realistic
}
