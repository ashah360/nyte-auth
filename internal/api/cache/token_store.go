package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ashah360/nyte-auth/internal/api/cerror"
	"github.com/ashah360/nyte-auth/internal/api/model"
)

const (
	UserSnapshotFieldName = "user_state"
)

type TokenStore interface {
	Create(ctx context.Context, userId, token string, state *model.UserSnapshot, expiry time.Duration) error
	Get(ctx context.Context, userId, token string) (*model.UserSnapshot, error)
	GetByUserID(ctx context.Context, userId string) ([]*model.UserSnapshot, error)
	Update(ctx context.Context, userId, token string, snapshot *model.UserSnapshot) error
	Delete(ctx context.Context, userId, token string) error
	DeleteByUserID(ctx context.Context, userId string) error
}

type tokenStore struct {
	redisClient *redis.Client
}

func (ts *tokenStore) Create(ctx context.Context, userId, token string, state *model.UserSnapshot, expiry time.Duration) error {
	k := hashKey(userId, token)
	if _, err := ts.redisClient.HSetNX(ctx, k, UserSnapshotFieldName, state).Result(); err != nil {
		return fmt.Errorf("create: redis hset error: %w", err)
	}

	if _, err := ts.redisClient.Expire(ctx, k, expiry).Result(); err != nil {
		return fmt.Errorf("create: redis expire error: %w", err)
	}

	return nil
}

func (ts *tokenStore) Get(ctx context.Context, userId, token string) (*model.UserSnapshot, error) {
	b, err := ts.redisClient.HGet(ctx, hashKey(userId, token), UserSnapshotFieldName).Bytes()
	if err == redis.Nil {
		return nil, cerror.ErrUserSnapshotDoesNotExist
	} else if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("find: not found")
	}

	us := new(model.UserSnapshot)

	if err = us.UnmarshalBinary(b); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}

	return us, nil
}

func (ts *tokenStore) GetByUserID(ctx context.Context, userId string) ([]*model.UserSnapshot, error) {
	var snapshots []*model.UserSnapshot

	keys, err := ts.redisClient.Keys(ctx, hashKey(userId, "*")).Result()
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		b, err := ts.redisClient.HGet(ctx, k, UserSnapshotFieldName).Bytes()
		if err != nil {
			return nil, err
		}

		var s model.UserSnapshot

		if err := json.Unmarshal(b, &s); err != nil {
			return nil, err
		}

		snapshots = append(snapshots, &s)
	}

	return snapshots, nil
}

func (ts *tokenStore) Update(ctx context.Context, userId, token string, snapshot *model.UserSnapshot) error {
	ttl, err := ts.redisClient.TTL(ctx, hashKey(userId, token)).Result()
	if err != nil {
		return err
	}

	if err := ts.redisClient.HSet(ctx, hashKey(userId, token), UserSnapshotFieldName, snapshot).Err(); err != nil {
		return err
	}

	return ts.redisClient.Expire(ctx, hashKey(userId, token), ttl).Err()
}

func (ts *tokenStore) Delete(ctx context.Context, userId, token string) error {
	return ts.redisClient.Del(ctx, hashKey(userId, token)).Err()
}

func (ts *tokenStore) DeleteByUserID(ctx context.Context, userId string) error {
	keys, err := ts.redisClient.Keys(ctx, hashKey(userId, "*")).Result()
	if err != nil {
		return err
	}

	// user has no tokens so there's nothing to delete
	if len(keys) == 0 {
		return nil
	}

	err = ts.redisClient.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func hashKey(userId string, token string) string {
	return fmt.Sprintf("user:%s:token:%s", userId, token)
}

func NewTokenStore(rds *redis.Client) TokenStore {
	return &tokenStore{
		redisClient: rds,
	}
}
