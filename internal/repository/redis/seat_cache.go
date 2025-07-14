package redis

import (
	"context"
	"fmt"
	"github.com/1612421/cinema-booking/pkg/go-kit/errorx"
	timesdk "github.com/1612421/cinema-booking/pkg/go-kit/time"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type ISeatCache interface {
	HoldSeat(ctx context.Context, dto HoldSeatCacheDTO) error
	GetHoldSeatKey(showtimeID, SeatID uuid.UUID) string
	ReleaseSeats(ctx context.Context, dto ReleaseSeatsBulkDTO) error
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
	  local user = ARGV[1]
	  for _, key in ipairs(KEYS) do
		if redis.call("GET", key) ~= user then
		  return 0
		end
	  end
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
	UserID     uuid.UUID
}

func NewSeatCache(redisCli redis.UniversalClient) ISeatCache {
	return &seatCache{
		cli: redisCli,
	}
}

func (sc *seatCache) GetHoldSeatKey(showtimeID, SeatID uuid.UUID) string {
	return fmt.Sprintf("hold:showtime:%s:seat:%s", showtimeID, SeatID)
}

func (sc *seatCache) HoldSeat(ctx context.Context, dto HoldSeatCacheDTO) error {
	if locked, err := sc.cli.SetNX(ctx, sc.GetHoldSeatKey(dto.ShowtimeID, dto.SeatID), dto.UserID.String(), HoldSeatTTL).Result(); err != nil || !locked {
		return errorx.New(http.StatusBadRequest, "seat is already held")
	}

	return nil
}

func (sc *seatCache) ReleaseSeats(ctx context.Context, dto ReleaseSeatsBulkDTO) error {
	if err := safeBulkDeleteScript.Load(ctx, sc.cli).Err(); err != nil {
		return fmt.Errorf("script load failed: %w", err)
	}

	keys := make([]string, len(dto.SeatIDs))
	for i, seatID := range dto.SeatIDs {
		keys[i] = sc.GetHoldSeatKey(dto.ShowtimeID, seatID)
	}

	released, err := safeBulkDeleteScript.Run(ctx, sc.cli, keys, dto.UserID.String()).Int()
	if err != nil || released == 0 {
		return errorx.New(http.StatusBadRequest, "you are not the owner or already expired")
	}

	return nil
}
