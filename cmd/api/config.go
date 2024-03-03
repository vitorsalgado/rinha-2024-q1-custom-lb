package main

import "github.com/vitorsalgado/rinha-2024-q1/internal/sys"

const (
	EnvAddr         = "ADDR"
	EnvDBConnString = "DB_CONN_STRING"
)

type Config struct {
	Addr         string
	DBConnString string
}

func newConfig() Config {
	config := Config{}
	config.Addr = sys.EnvStr(EnvAddr, ":8080")
	config.DBConnString = sys.EnvStr(EnvDBConnString, "postgresql://rinha:rinha@db:5432/rinha?sslmode=disable")

	return config
}
