# cinema-booking-system-backend
[Vue 3 Frontend]
   |  REST API / WebSocket
   v
[Go Echo Backend]
   |------> [MongoDB]   (movies, showtimes, bookings, users, audit_logs)
   |------> [Redis]     (seat locks, pub/sub)
   |------> [Cleanup Job]
   |------> [MQ Consumer]

### Functional
- [ ] login ได้
- [ ] movies/showtimes/seats ใช้งานได้
- [ ] lock seat ได้
- [ ] confirm booking ได้
- [ ] release booking ได้
- [ ] timeout cleanup ได้
- [ ] realtime refresh ได้
- [ ] admin bookings ดูได้
- [ ] admin audit logs ดูได้

### Non-functional
- [ ] ไม่มี double booking ใน flow หลัก
- [ ] role แยก USER / ADMIN
- [ ] env ไม่ hardcode
- [ ] repo structure

## Run with Docker

1. Create environment files:
   - `backend/.env`
   - `frontend/.env`

2. Place Firebase service account file at:
   - `backend/firebase-key.json`

3. Start all services:

```bash
docker compose up --build