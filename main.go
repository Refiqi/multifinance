package main

import (
	"engine.multifinance.com/cache"
	engineConfig "engine.multifinance.com/config"
	"engine.multifinance.com/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"multifinance.com/multifinance/config"
	"multifinance.com/multifinance/handlers"
	"multifinance.com/multifinance/middleware"
	"multifinance.com/multifinance/model"
	"net/http"
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
	postgresDB := sql.InitDB(configMap.DB)
	postgresDB.AutoMigrate(&model.User{}, &model.Transaction{}, &model.LoanLimit{})

	lruCache := cache.NewDoubleBufferLru(configMap.DoubleBufferLruConfig)

	r := gin.Default()

	r.Use(middleware.InjectDBToContext(postgresDB), middleware.InjectCacheToContext(lruCache), middleware.HeaderPolicy, middleware.ValidateParams)

	handlers.SetupLoanLimitRoutes(r)
	handlers.SetupTransactionRoutes(r)
	handlers.SetupUserRoutes(r)

	// Start the server
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
