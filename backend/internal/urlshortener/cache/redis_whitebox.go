package redis

import "context"

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
