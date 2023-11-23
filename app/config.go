package app

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   8081,
	}
	if redisAddress, ok := os.LookupEnv("REDIS_ADDRESS"); ok {
		cfg.RedisAddress = redisAddress
	}
	if serverPort, ok := os.LookupEnv("SERVER_PORT"); ok {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}
