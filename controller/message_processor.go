package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VsenseTechnologies/skf_mqtt_message_processor/model"
	"github.com/VsenseTechnologies/skf_mqtt_message_processor/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MessageProcessor(c mqtt.Client, m mqtt.Message, cacheRepo *repository.RedisRepository, dbRepo *repository.PostgresRepository) {

	clientId := strings.Split(m.Topic(), "/")[0]
	var rawMessage model.RawMessage

	if err := json.Unmarshal(m.Payload(), &rawMessage); err != nil {
		log.Printf("error occurred while decoding the mqtt message to json, Error -> %v\n", err.Error())
	}

	start := time.Now()

	channel1 := make(chan any)
	channel2 := make(chan any)
	channel3 := make(chan any)

	go func() {
		registerType, err := cacheRepo.GetRegisterType(clientId, rawMessage.RegisterAddress)
		if err != nil {
			channel1 <- err
		}
		channel1 <- registerType
	}()

	go func() {
		drierId, err := cacheRepo.GetDrierId(clientId, rawMessage.RegisterAddress)
		if err != nil {
			channel2 <- err
		}
		channel2 <- drierId
	}()

	go func() {
		registerValue, err := cacheRepo.GetRegisterValue(clientId, rawMessage.RegisterAddress)

		if err != nil {
			channel3 <- err
		}
		channel3 <- registerValue
	}()

	res1 := <-channel1
	res2 := <-channel2
	res3 := <-channel3

	value, ok := res1.(error)

	if ok {
		log.Printf("error occurred with redis while getting the register type, Error -> %v\n", value.Error())
		return
	}

	registerType := res1.(string)

	value, ok = res2.(error)

	if ok {
		log.Printf("error occurred with redis while getting the drier id, Error -> %v\n", value.Error())
		return
	}

	drierId := res2.(string)

	value, ok = res3.(error)

	if ok {
		log.Printf("error occurred with redis while getting the register value, Error -> %v\n", value.Error())
		return
	}

	registerValue := res3.(string)

	end := time.Since(start)

	fmt.Println("took -> ", end)

	splittedRegisterType := strings.Split(registerType, "_")

	fmt.Println(registerType)

	if splittedRegisterType[0] == "rt" {

		message := map[string]string{
			registerType: rawMessage.Data,
		}

		jsonMessage, err := json.Marshal(message)

		if err != nil {
			log.Printf("error occurred while encoding the message to json format, Error -> %v\n", err.Error())
			return
		}

		c.Publish(drierId, 0, false, jsonMessage)

		if registerValue != rawMessage.Data {
			if registerType == "rt_pid" {
				if err := cacheRepo.UpdateDrierPid(drierId, rawMessage.Data); err != nil {
					log.Printf("error occurred with redis while updating the drier pid value, Error -> %v\n", err.Error())
					return
				}
			}
			if err := cacheRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
				log.Printf("error occurred with redis while updating the register value, Error -> %v\n", err.Error())
				return
			}
			if err := dbRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
				log.Printf("error occurred with database while updating the register value, Error -> %v\n", err.Error())
				return
			}
		}
		return
	}

	if splittedRegisterType[0] == "rcp" {
		if splittedRegisterType[3] == "tp" {
			recipeStepCount, err := cacheRepo.GetDrierRecipeStepCount(drierId)

			if err != nil {
				log.Printf("error occurred with redis while getting drier recipe step count, Error -> %v\n", err.Error())
				return
			}

			if splittedRegisterType[2] == recipeStepCount {
				if err := cacheRepo.UpdateDrierRecipeTemperature(drierId, rawMessage.Data); err != nil {
					log.Printf("error occurred with redis while updating the drier recipe temperature, Error -> %v\n", err.Error())
					return
				}
			}

			if registerValue != rawMessage.Data {
				if err := cacheRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
					log.Printf("error occurred with redis while updating the register value, Error -> %v\n", err.Error())
					return
				}
				if err := dbRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
					log.Printf("error occurred with database while updating the register value, Error -> %v\n", err.Error())
					return
				}
			}
			return
		}

		maxTimeStr := os.Getenv("DRIER_RECIPE_MAX_TIME")

		if maxTimeStr == "" {
			log.Printf("missing or empty env variable DRIER_RECIPE_MAX_TIME")
			return
		}

		maxTime, err := strconv.Atoi(maxTimeStr)

		if err != nil {
			log.Printf("error occurred while parsing drier recipe max time from string to integer, Error -> %v\n", err.Error())
			return
		}

		currentTime, err := strconv.Atoi(rawMessage.Data)

		if err != nil {
			log.Printf("error occurred while parsing drier recipe current time from string to integer, Error -> %v\n", err.Error())
			return
		}

		if currentTime > 0 && currentTime < maxTime {
			recipeStepCount, err := cacheRepo.GetDrierRecipeStepCount(drierId)
			if err != nil {
				log.Printf("error occurred with redis while getting drier recipe step count, Error -> %v\n", err.Error())
				return
			}

			if splittedRegisterType[2] != recipeStepCount {
				if err := cacheRepo.UpdateDrierRecipeStepCount(drierId, splittedRegisterType[2]); err != nil {
					log.Printf("error occurred with redis while updating drier recipe step count, Error -> %v\n", err.Error())
					return
				}
			}

			recipeTemperature, err := cacheRepo.GetDrierRecipeTemperature(drierId)

			if err != nil {
				log.Printf("error occurred with redis while getting the drier recipe temperature, Error -> %v\n", err.Error())
				return
			}

			recipeStepMessage := &model.RecipeStep{
				StepCount:   splittedRegisterType[2],
				Time:        rawMessage.Data,
				Temperature: recipeTemperature,
			}

			recipeStepJsonMessage, err := json.Marshal(recipeStepMessage)

			if err != nil {
				log.Printf("error occurred while encoding recipe step message to json, Error -> %v\n", err.Error())
				return
			}

			c.Publish(drierId, 1, false, recipeStepJsonMessage)

			drierPid, err := cacheRepo.GetDrierPid(drierId)

			if err != nil {
				log.Printf("error occurred with redis while getting the drier pid, Error -> %v\n", err)
				return
			}

			batch := &model.Batch{
				DrierId:     drierId,
				RecipeStep:  splittedRegisterType[2],
				Time:        rawMessage.Data,
				Temperature: recipeTemperature,
				Pid:         drierPid,
			}

			if err := dbRepo.CreateBatch(batch); err != nil {
				log.Printf("error occurred with database while creating the batch, Error -> %v\n", err.Error())
				return
			}

		}

		if registerValue != rawMessage.Data {
			if err := cacheRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
				log.Printf("error occurred with redis while updating the register value, Error -> %v\n", err.Error())
				return
			}
			if err := dbRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
				log.Printf("error occurred with database while updating the register value, Error -> %v\n", err.Error())
				return
			}
		}

		return
	}

	if registerValue != rawMessage.Data {
		if err := cacheRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
			log.Printf("error occurred with redis while updating the register value, Error -> %v\n", err.Error())
		}
		if err := dbRepo.UpdateRegisterValue(clientId, rawMessage.RegisterAddress, rawMessage.Data); err != nil {
			log.Printf("error occurred with database while updating the register value, Error -> %v\n", err.Error())
		}

		message := map[string]string{
			registerType: rawMessage.Data,
		}

		jsonMessage, err := json.Marshal(message)

		if err != nil {
			log.Printf("error occurred while encoding message to json format, Error -> %v\n", err.Error())
			return
		}

		c.Publish(drierId, 1, false, jsonMessage)

	}
}
