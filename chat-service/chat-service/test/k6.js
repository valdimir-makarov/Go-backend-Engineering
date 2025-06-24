import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  vus: 50,       // 50 virtual users
  duration: '30s', // for 30 seconds
};

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fakePayload.signature"; // Replace with actual if needed
const groupId = "11111111-1111-1111-1111-111111111111";

export default function () {
  const url = `ws://yourdomain.com/wsgrpmsg?user_id=${__VU}&token=${token}`;

  const res = ws.connect(url, null, function (socket) {
    socket.on('open', function () {
      console.log(`VU ${__VU} connected`);

      const message = {
        sender_id: __VU,        // Virtual user ID
        group_id: groupId,
        content: `Hello from VU ${__VU}`,
      };

      // Send message as JSON string
      socket.send(JSON.stringify(message));

      // Wait for a response
      socket.on('message', function (msg) {
        console.log(`VU ${__VU} received: ${msg}`);
      });

      sleep(1); // Let the socket stay open briefly

      socket.close();
    });

    socket.on('close', function () {
      console.log(`VU ${__VU} disconnected`);
    });

    socket.on('error', function (e) {
      console.log(`VU ${__VU} error: ${e.error()}`);
    });
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}
