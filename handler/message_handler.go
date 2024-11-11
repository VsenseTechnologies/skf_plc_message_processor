package handler

import (
	"database/sql"

	"github.com/VsenseTechnologies/skf_mqtt_message_processor/controller"
	"github.com/VsenseTechnologies/skf_mqtt_message_processor/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
)

func Handler(redisConn *redis.Client, postgresConn *sql.DB) mqtt.MessageHandler {
	redisRepo := repository.NewRedisRepository(redisConn)
	postgresRepo := repository.NewPostgresRepository(postgresConn)
	return func(c mqtt.Client, m mqtt.Message) {
		go controller.MessageProcessor(c, m, redisRepo, postgresRepo)
	}
}
