package server

import (
	"fmt"
	redis "github.com/go-redis/redis/v7"
)

type redisServer struct{
	redisClient  *redis.Client
	data  chan string
}

func NewRedisServer(addr string , db int, channeldata  chan string) (*redisServer, error){
	redis_server := redisServer{}
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
		DB: db,
	})
	_, err := rdb.Ping().Result()
	if err != nil{
		return &redis_server, err
	}
	redis_server.redisClient = rdb
	redis_server.data = channeldata
	return &redis_server, nil
}

func (server *redisServer) StartScribe(){
	pubsub := server.redisClient.Subscribe("name")
	_, err := pubsub.Receive()
    if err != nil {

    }
	for {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Println(msg.Channel, ":", msg.Payload)
			server.data  <- string(msg.Payload)
		}
	}
}