package configs

import "github.com/jackc/pgx/v5/pgxpool"


func InitDbConfig()*pgxpool.Pool{
	pg := InitDB()

	return pg
}