package repository

import (
	"database/sql"
	"fmt"

	"github.com/VsenseTechnologies/skf_mqtt_message_processor/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db,
	}
}

func (repo *PostgresRepository) UpdateRegisterValue(plcId string, regAddress string, value string) error {
	query := fmt.Sprintf(`UPDATE %s SET value=$2 WHERE reg_address=$1`, plcId)
	_, err := repo.db.Exec(query, regAddress, value)
	return err
}

func (repo *PostgresRepository) CreateBatch(batch *model.Batch) error {
	query := `INSERT INTO batches (drier_id,recipe_step,set_time,real_time,real_temp,real_pid) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := repo.db.Exec(query, batch.DrierId, batch.RecipeStep, batch.SetTime, batch.RealTimeTime, batch.RealTimeTemperature, batch.RealTimePid)
	return err
}
