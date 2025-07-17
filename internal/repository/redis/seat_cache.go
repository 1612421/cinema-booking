package redis

import (
	"context"
	"fmt"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	timesdk "github.com/1612421/cinema-booking/pkg/go-kit/time"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type ISeatCache interface {
	HoldSeat(ctx context.Context, dto HoldSeatCacheDTO) (int64, error)
	GetHoldSeatKey(showtimeID, SeatID uuid.UUID) string
	HDelHoldSeatsAll(ctx context.Context, dto ReleaseSeatsBulkDTO) error
	ReleaseSeat(ctx context.Context, dto ReleaseSeatDTO) error
	GetSeatInterest(ctx context.Context, dto SeatInterestDTO) (int64, error)
}

type seatCache struct {
	cli redis.UniversalClient
}

const (
	SeatCacheKey     = "seat"
	HoldSeatCacheKey = "seat_hold"
)

var (
	HoldSeatTTL          = timesdk.MinuteDuration(10)
	safeBulkDeleteScript = redis.NewScript(`
	  for _, key in ipairs(KEYS) do
		redis.call("DEL", key)
	  end
	  return 1
	`)
)

type HoldSeatCacheDTO struct {
	SeatID     uuid.UUID
	ShowtimeID uuid.UUID
	UserID     uuid.UUID
}

type ReleaseSeatCacheDTO struct {
	SeatID     uuid.UUID
	ShowtimeID uuid.UUID
	UserID     uuid.UUID
}

type ReleaseSeatsBulkDTO struct {
	SeatIDs    []uuid.UUID
	ShowtimeID uuid.UUID
}

type ReleaseSeatDTO struct {
	SeatID     uuid.UUID
	ShowtimeID uuid.UUID
	UserID     uuid.UUID
}

type SeatInterestDTO struct {
	SeatID     uuid.UUID
	ShowtimeID uuid.UUID
}

func NewSeatCache(redisCli redis.UniversalClient) ISeatCache {
	return &seatCache{
		cli: redisCli,
	}
}

func (sc *seatCache) GetHoldSeatKey(showtimeID, seatID uuid.UUID) string {
	return fmt.Sprintf("hold:showtime:%s:seat:%s", showtimeID, seatID)
}

func (sc *seatCache) GetHoldSeatUserField(userId uuid.UUID) string {
	return fmt.Sprintf("user:%s", userId)
}

func (sc *seatCache) HoldSeat(ctx context.Context, dto HoldSeatCacheDTO) (int64, error) {
	//if locked, err := sc.cli.SetNX(ctx, sc.GetHoldSeatKey(dto.ShowtimeID, dto.SeatID), dto.UserID.String(), HoldSeatTTL).Result(); err != nil || !locked {
	//	return errorx.New(http.StatusBadRequest, "seat is already held")
	//}

	key := sc.GetHoldSeatKey(dto.ShowtimeID, dto.SeatID)
	field := sc.GetHoldSeatUserField(dto.UserID)
	err := sc.cli.HSet(ctx, key, field, time.Now().Unix()).Err()
	sc.cli.Expire(ctx, key, HoldSeatTTL)
	if err != nil {
		return 0, errorx.New(http.StatusBadRequest, err.Error())
	}

	count, err := sc.cli.HLen(ctx, key).Result()
	if err != nil {
		return 0, errorx.New(http.StatusBadRequest, err.Error())
	}

	return count, nil
}

func (sc *seatCache) ReleaseSeat(ctx context.Context, dto ReleaseSeatDTO) error {
	key := sc.GetHoldSeatKey(dto.ShowtimeID, dto.SeatID)
	field := sc.GetHoldSeatUserField(dto.UserID)
	err := sc.cli.HDel(ctx, key, field).Err()
	if err != nil {
		return errorx.New(http.StatusBadRequest, err.Error())
	}

	count, err := sc.cli.HLen(ctx, key).Result()
	if err != nil {
		return errorx.New(http.StatusBadRequest, err.Error())
	}
	if count == 0 {
		sc.cli.Del(ctx, key)
	}

	return nil
}

func (sc *seatCache) HDelHoldSeatsAll(ctx context.Context, dto ReleaseSeatsBulkDTO) error {
	if err := safeBulkDeleteScript.Load(ctx, sc.cli).Err(); err != nil {
		return fmt.Errorf("script load failed: %w", err)
	}

	keys := make([]string, len(dto.SeatIDs))
	for i, seatID := range dto.SeatIDs {
		keys[i] = sc.GetHoldSeatKey(dto.ShowtimeID, seatID)
	}

	released, err := safeBulkDeleteScript.Run(ctx, sc.cli, keys).Int()
	if err != nil || released == 0 {
		return errorx.New(http.StatusBadRequest, "you are not the owner or already expired")
	}

	return nil
}

func (sc *seatCache) GetSeatInterest(ctx context.Context, dto SeatInterestDTO) (int64, error) {
	key := sc.GetHoldSeatKey(dto.ShowtimeID, dto.SeatID)
	count, err := sc.cli.HLen(ctx, key).Result()
	if err != nil {
		return 0, errorx.New(http.StatusBadRequest, err.Error())
	}

	return count, nil
}
