import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const uuid = 'a98216bb-9b3d-4d57-8da3-dda89ac03b51';  // Replace with a valid user UUID
  const response = http.get(`http://localhost:8088/users/${uuid}/balances`);
  check(response, { 'status was 200': (r) => r.status === 200 });
}
