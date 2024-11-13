package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client,
	}
}

func (repo *RedisRepository) GetRegisterType(plcId string, regAddress string) (string, error) {
	key := fmt.Sprintf("rg_ty_%s_%s", plcId, regAddress)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) GetDrierId(plcId string, regAddress string) (string, error) {
	key := fmt.Sprintf("dr_id_%s_%s", plcId, regAddress)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) GetRegisterValue(plcId string, regAddress string) (string, error) {
	key := fmt.Sprintf("rg_vl_%s_%s", plcId, regAddress)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) UpdateRegisterValue(plcId string, regAddress string, value string) error {
	key := fmt.Sprintf("rg_vl_%s_%s", plcId, regAddress)
	_, err := repo.client.Set(context.Background(), key, value, 0).Result()
	return err
}

func (repo *RedisRepository) GetDrierRecipeStepCount(drierId string) (string, error) {
	key := fmt.Sprintf("rcp_stp_ct_%s", drierId)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) UpdateRecipeStepCompleteStatus(drierId string, status string) error {
	key := fmt.Sprintf("rcp_stp_cmp_%s", drierId)
	_, err := repo.client.Set(context.Background(), key, status, 0).Result()
	return err
}

func (repo *RedisRepository) UpdateRecipeSetTime(drierId string, time string) error {
	key := fmt.Sprintf("rcp_stp_stm_%s", drierId)
	_, err := repo.client.Set(context.Background(), key, time, 0).Result()
	return err
}

func (repo *RedisRepository) UpdateDrierRecipeRealTimeTemperature(drierId string, temp string) error {
	key := fmt.Sprintf("rcp_stp_rtp_%s", drierId)
	_, err := repo.client.Set(context.Background(), key, temp, 0).Result()
	return err
}

func (repo *RedisRepository) UpdateDrierRecipeStepCount(drierId string, count string) error {
	key := fmt.Sprintf("rcp_stp_ct_%s", drierId)
	_, err := repo.client.Set(context.Background(), key, count, 0).Result()
	return err
}

func (repo *RedisRepository) GetDrierRecipeRealTimeTemperature(drierId string) (string, error) {
	key := fmt.Sprintf("rcp_stp_rtp_%s", drierId)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) GetDrierRecipeSetTime(drierId string) (string, error) {
	key := fmt.Sprintf("rcp_stp_stm_%s", drierId)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}

func (repo *RedisRepository) UpdateDrierPid(drierId string, pid string) error {
	key := fmt.Sprintf("pid_%s", drierId)
	_, err := repo.client.Set(context.Background(), key, pid, 0).Result()
	return err
}

func (repo *RedisRepository) GetDrierPid(drierId string) (string, error) {
	key := fmt.Sprintf("pid_%s", drierId)
	result, err := repo.client.Get(context.Background(), key).Result()
	return result, err
}
