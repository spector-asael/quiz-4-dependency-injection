package handler

import (
	"log/slog"
	"github.com/spector-asael/lab4-crud/internal/data"
)

type ServerConfig struct {
	Port int
	Environment  string
	DB struct {
        DSN string
    }

}

type ApplicationDependencies struct {
	Config ServerConfig
	Logger *slog.Logger
	Models data.Models
}