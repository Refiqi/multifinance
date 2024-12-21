package main

import (
	engineConfig "engine.multifinance.com/config"
	"engine.multifinance.com/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"multifinance.com/multifinance/config"
	"os"
)

func main() {

	envMode := os.Getenv("ENV_MODE")

	if len(envMode) == 0 {
		log.Fatalf("ENV_MODE is not set yet")
		return
	}

	var configMap config.ConfigMap
	engineConfig.LoadConfig(&configMap, envMode)
	postgresqlDB := sql.InitDB(configMap.DB)


	r := gin.Default()

	r.Use(gin.Logger(), gin.Recovery())


	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
