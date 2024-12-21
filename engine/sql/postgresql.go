package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

type PostgresqlConfig struct {
	Host                string
	Port 				string
	DbName              string
	Username            string
	Password            string
	Ssl                 string
	MaxOpenConnections  int
	MaxIdleConnections  int
}

type ClientConfig interface {
	GetConnectionConfig() ConnectionConfig
	GetConnectionSpec() ConnectionSpec
}

// ConnectionConfig ...
type ConnectionConfig struct {
	MaxOpenConnections int
	MaxIdleConnections int
}

// ConnectionSpec ...
type ConnectionSpec struct {
	DriverName     string
	DataSourceName string
}

type sqlClient struct {
	GormDB *gorm.DB
}

// GetConnectionSpec ...
func (ths PostgresqlConfig) GetConnectionSpec() ConnectionSpec {
	return ConnectionSpec{
		DriverName:     "postgres",
		DataSourceName: "host=" + ths.Host + "user=" + ths.Username + "password=" + ths.Password + "dbname=" + ths.DbName + "port=" + ths.Port + "sslmode=" + ths.Ssl + "TimeZone=Asia/Jakarta"}
}
// GetConnectionConfig ...
func (ths PostgresqlConfig) GetConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		MaxOpenConnections: ths.MaxOpenConnections,
		MaxIdleConnections: ths.MaxIdleConnections}
}

func InitDB(config ClientConfig) *gorm.DB {
	var err error
	var once sync.Once
	var gormDB *gorm.DB

	connectionSpec := config.GetConnectionSpec()
	once.Do(func() {
		gormDB, err = gorm.Open(postgres.Open(connectionSpec.DataSourceName), &gorm.Config{})
		if err != nil {
			log.Fatalf("Gagal terhubung ke database: %v", err)
		}
		log.Println("berhasil terhubung ke database")
	})
	return gormDB
}
