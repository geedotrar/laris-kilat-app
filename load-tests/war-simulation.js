import http from 'k6/http';
import { check } from 'k6';

const TARGET_MODE = __ENV.MODE || 'naive';

export const options = {
    scenarios: {
        flash_sale_rush: {
            executor: 'per-vu-iterations',
            vus: 500,        // 500 virtual users
            iterations: 1,   // Each user presses the buy button once 
            maxDuration: '10s',
        },
    },
};

export default function () {
    const url = `http://host.docker.internal:8081/api/v1/checkout/${TARGET_MODE}`;
    const payload = JSON.stringify(
        {
            user_id: `user-${__VU}`,
            item_id: 'item-sneaker-01',
        }
    );

    const res = http.post(url, payload, {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        '200 OK (Success)': (r) => r.status === 200,
        '409 Conflict (Out of Stock / Duplicate)': (r) => r.status === 409,
    });
}