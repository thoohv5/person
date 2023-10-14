package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"

	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
)

type IRedis interface {
	Execute(ctx context.Context, commands ...func(ctx context.Context, client *redis.Client) error) error
}

type defaultRedis struct {
	r   *redis.Client
	log log.Logger
}

func NewRedis(rc *Config, log log.Logger) (IRedis, func(), error) {

	dr := &defaultRedis{
		log: log,
	}

	r := redis.NewClient(&redis.Options{
		Network:      rc.GetNetwork(),
		Addr:         rc.GetAddr(),
		Password:     rc.GetPassword(),
		DB:           int(rc.GetDB()),
		DialTimeout:  time.Second * time.Duration(rc.GetConnectionTimeout()),
		ReadTimeout:  time.Second * time.Duration(rc.GetReadTimeout()),
		WriteTimeout: time.Second * time.Duration(rc.GetWriteTimeout()),
		MinIdleConns: int(rc.GetMaxIdle()),
	})

	r.AddHook(NewHook(otel.Tracer("github.com/go-redis/redis")))
	dr.r = r

	return dr, func() {
		if err := r.Close(); nil != err {
			dr.log.Errorc(context.Background(), "redis close err", logger.FieldError(err))
		}
	}, nil

}

func (dr *defaultRedis) Execute(ctx context.Context, commands ...func(ctx context.Context, client *redis.Client) error) error {

	for _, command := range commands {
		if err := command(ctx, dr.r); nil != err {
			// dr.log.Errorc(ctx, "redis command err", logger.FieldError(err))
			return err
		}
	}

	return nil

}
