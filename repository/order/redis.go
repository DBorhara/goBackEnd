package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/DBorhara/goBackEnd/model"
)

type RedisRepository struct {
	Client *redis.Client
}

func orderIDKey(orderID uint64) string {
	return fmt.Sprintf("order:%d", orderID)
}

func (r *RedisRepository) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order: %w", err)
	}

	key := orderIDKey(order.OrderId)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)

	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("error inserting order: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("error inserting order: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("error inserting order: %w", err)
	}
	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepository) FindById(ctx context.Context, orderID uint64) (model.Order, error) {
	key := orderIDKey(orderID)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("error getting order: %w", err)
	}

	var order model.Order

	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("error unmarshalling order: %w", err)
	}
	return order, nil
}

func (r *RedisRepository) DeleteById(ctx context.Context, orderID uint64) error {
	key := orderIDKey(orderID)
	txn := r.Client.TxPipeline()
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("error deleting order: %w", err)
	}
	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("error deleting order: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}
	return nil
}

func (r *RedisRepository) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error marshalling order: %w", err)
	}
	key := orderIDKey(order.OrderId)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("error updating order: %w", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindAllResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepository) FindAll(ctx context.Context, page FindAllPage) (FindAllResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()

	if err != nil {
		return FindAllResult{}, fmt.Errorf("error getting orders: %w", err)
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()

	if err != nil {
		return FindAllResult{}, fmt.Errorf("error getting orders: %w", err)
	}

	if len(keys) == 0 {
		return FindAllResult{
			Orders: []model.Order{},
		}, nil
	}
	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindAllResult{}, fmt.Errorf("error unmarshalling order: %w", err)
		}
		orders[i] = order
	}
	return FindAllResult{
		Orders: orders,
		Cursor: uint64(cursor),
	}, nil
}
