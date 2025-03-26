package api

import "github.com/magicznykacpur/psst-backend/internal/database"

type ApiConfig struct {
	DB   *database.Queries
	Port string
}
