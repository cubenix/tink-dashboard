package redis

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

const (
	envRedisAddr = "REDIS_ADDRESS"
	envRedisPass = "REDIS_PASSWORD"

	infoRedisInit = "initializing redis client"
	errRedisAddr  = "required env variable REDIS_ADDRESS not defined"
	errConnection = "unable to connect with redis server"
)

var client *Cache

// Cache represents the redis cache client
type Cache struct {
	client *redis.Client
}

// Instance returns an existing client instance; creates new otherwise
func Instance() *Cache {
	if client == nil {
		client = new()
	}
	return client
}

// new creates a new redis cache client
func new() *Cache {
	log.Info(infoRedisInit)
	addr := os.Getenv(envRedisAddr)
	if addr == "" {
		log.Fatal(errRedisAddr)
	}
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv(envRedisPass),
	})

	if _, err := rc.Ping().Result(); err != nil {
		log.Fatal(errConnection)
	}
	return &Cache{client: rc}
}

func init() {
	_ = Instance()
}
