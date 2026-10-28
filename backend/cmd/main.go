package main

import (
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/egorvts/Guestbook/internal/handler"
	"github.com/egorvts/Guestbook/internal/model"
	"github.com/egorvts/Guestbook/internal/repository"
	ydb "github.com/ydb-platform/gorm-driver"
	ydbSDK "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/balancers"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func main() {
	initDB()
	
	if err := migrate(); err != nil {
		panic(err.Error())
	}

	handler := handler.NewHandler(repository.NewRepository(db))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/messages", handler.List)
	mux.HandleFunc("POST /api/messages", handler.Create)

	server := &http.Server{
		Addr:    		":" + getenvOrPanic("PORT"),
		Handler: 		mux,
		ReadTimeout:  	5 * time.Second,
		WriteTimeout: 	10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic("Server failed: " + err.Error())
	}
}

func initDB() {
	dbOnce.Do(func() {
		dsn := getenvOrPanic("YDB_CONNECTION_STRING")
		if dsn == "" {
			panic("YDB_CONNECTION_STRING environment variable is not set")
		}

		var err error
		credentials := yc.WithCredentials()
		if token := os.Getenv("YDB_TOKEN"); token != "" {
			credentials = ydbSDK.WithAccessTokenCredentials(token)
		}

		db, err = gorm.Open(ydb.Open(dsn, 
			ydb.With(
				ydbSDK.WithBalancer(balancers.SingleConn()), 
				credentials, 
				yc.WithInternalCA(),
			),
			ydb.WithMaxOpenConns(10),
			ydb.WithMaxIdleConns(5),
		))

		if err != nil {
			panic("Failed to connect to YDB: " + err.Error())
		}
	})
}

func migrate() error {
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		return errors.New("Failed to apply migrations: " + err.Error())
	}

	return nil
}

func getenvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Environment variable " + key + " is not set")
	}
	return value
}