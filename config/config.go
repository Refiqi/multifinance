package config

import (
	"engine.multifinance.com/cache"
	"engine.multifinance.com/sql"
)

type ConfigMap struct {
	DB 		 			  sql.PostgresqlConfig
	DoubleBufferLruConfig cache.DoubleBufferLruConfig
}
