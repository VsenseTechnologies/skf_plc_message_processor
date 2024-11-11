package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/VsenseTechnologies/skf_mqtt_message_processor/cache"
	"github.com/VsenseTechnologies/skf_mqtt_message_processor/db"
	"github.com/VsenseTechnologies/skf_mqtt_message_processor/handler"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load env variable Error -> %v\n", err.Error())
	}

	//initializing the logger
	initLogger()

	cacheConn, err := cache.Connect()

	if err != nil {
		log.Fatalf("error occurred while connecting to redis, Error -> %v\n", err.Error())
	}

	fmt.Println("connected to redis")

	dbConn, err := db.Connect()

	fmt.Println("connected to database ")

	if err != nil {
		log.Fatalf("error occurred while connecting to database Error -> %v\n", err.Error())
	}

	var brokerHost = os.Getenv("BROKER_HOST")

	if brokerHost == "" {
		log.Fatalf("missing or empty env variable S2_BROKER_HOST \n")
	}

	var brokerPort = os.Getenv("BROKER_PORT")

	if brokerPort == "" {
		log.Fatalf("missing or empty env variable BROKER_PORT\n")
	}

	var clientId = os.Getenv("CLIENT_ID")

	if clientId == "" {
		log.Fatalf("missing or empty env variable CLIENT_ID\n")
	}

	var brokerAddress = fmt.Sprintf("tcp://%s:%s", brokerHost, brokerPort)

	opts := mqtt.NewClientOptions()

	opts.AddBroker(brokerAddress)

	opts.SetClientID(clientId)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("connected to broker")
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("disconnected from the broker, Error -> %v\n", err.Error())
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("error occurred while connecting to broker, Error -> %v\n", token.Error())
	}

	if client.IsConnected() {
		client.Subscribe("+/message/processor", 1, handler.Handler(cacheConn, dbConn))
	}

	for {
		if !client.IsConnected() {

			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Printf("error occurred while connecting to broker, Error -> %v\n", token.Error())
				continue
			}

			client.Unsubscribe("+/message/processor")
			client.Subscribe("+/message/processor", 1, handler.Handler(cacheConn, dbConn))

			log.Println("reconnected to broker")

		}

		time.Sleep(time.Second * 1)
	}

}
