package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const SeatLockTTL = 5 * time.Minute

type SeatLockValue struct {
	BookingID  string `json:"booking_id"`
	UserID     string `json:"user_id"`
	ShowtimeID string `json:"showtime_id"`
	SeatID     string `json:"seat_id"`
	ExpiresAt  string `json:"expired_at"`
}

type SeatLockService struct {
	redisClient *redis.Client
}

func NewSeatLockService(redisClient *redis.Client) *SeatLockService {
	return &SeatLockService{
		redisClient: redisClient,
	}
}

func (s *SeatLockService) BuildKey(showtimeID, seatID string) string {
	return fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeID, seatID)
}

// แบบเก่า ไม่ใช้แล้ว
func (s *SeatLockService) LockSeat(
	ctx context.Context,
	showtimeID string,
	seatID string,
	value SeatLockValue,
) (bool, error) {
	key := s.BuildKey(showtimeID, seatID)

	raw, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	ok, err := s.redisClient.SetNX(ctx, key, raw, SeatLockTTL).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

//แบบใหม่
// func (s *SeatLockService) LockSeat(
// 	ctx context.Context,
// 	showtimeID string,
// 	seatID string,
// 	value SeatLockValue,
// ) (error) {
// 	key := s.BuildKey(showtimeID, seatID)

// 	raw, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}

// 	ok := s.redisClient.Set(ctx, key, raw, SeatLockTTL)
// 	if ok.Err() == redis.Nil {
// 		// Key was not set because it already exists
// 		fmt.Println("Key already exists, not set")
// 	} else if ok.Err() != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *SeatLockService) GetSeatLock(
	ctx context.Context,
	showtimeID string,
	seatID string,
) (*SeatLockValue, error) {
	key := s.BuildKey(showtimeID, seatID)

	result, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var value SeatLockValue
	if err := json.Unmarshal([]byte(result), &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func (s *SeatLockService) ReleaseSeatLock(
	ctx context.Context,
	showtimeID string,
	seatID string,
) error {
	key := s.BuildKey(showtimeID, seatID)
	return s.redisClient.Del(ctx, key).Err()
}

func (s *SeatLockService) GetSeatLockTTL(
	ctx context.Context,
	showtimeID string,
	seatID string,
) (time.Duration, error) {
	key := s.BuildKey(showtimeID, seatID)
	return s.redisClient.TTL(ctx, key).Result()
}
