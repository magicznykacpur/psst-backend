package api

import (
	"github.com/magicznykacpur/psst-backend/internal/database"
	"github.com/magicznykacpur/psst-backend/ws"
)

type ApiConfig struct {
	DB   *database.Queries
	Port string
	Hub  *ws.Hub
}
