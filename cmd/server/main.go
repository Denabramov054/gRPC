package main

import (
	//"crypto/tls"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"grpc/internal/repository"
	"grpc/internal/service"
	"grpc/transport"
	_ "grpc/transport"
	"grpc/transport/handler"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error while reading config, %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	clickDB := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{
			os.Getenv("CLICKHOUSE_HOST") +
				":" +
				os.Getenv("CLICKHOUSE_PORT"),
		},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: true,
	})
	clickDB.SetMaxIdleConns(5)
	clickDB.SetMaxOpenConns(10)
	clickDB.SetConnMaxLifetime(time.Hour)

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  viper.GetString("db.sslMode"),
	})
	if err != nil {
		log.Fatal("Connect to db err: ", err)
	}

	rdb, err := repository.NewRedisClient(repository.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: viper.GetString("rdb.password"),
		DB:       viper.GetInt("rdb.db"),
	})
	if err != nil {
		log.Fatal("Connect to rdb err: ", err)
	}

	repos := repository.NewRepository(db, rdb,clickDB)
	services := service.NewService(repos)

	deps := transport.Deps{
		UserHandler: handler.NewUserHandler(services.User),
	}
	grpcServer := transport.NewServer(deps)
	grpcConfig := transport.ServerConfig{
		Host: os.Getenv("APP_HOST"),
		Port: os.Getenv("APP_PORT"),
	}
	go func() {
		if err = grpcServer.ListenAndServe(grpcConfig); err != nil {
			log.Println("grpc ListenAndServe error", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Println("Shutdown serv...")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
