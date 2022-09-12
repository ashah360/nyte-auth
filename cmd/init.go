package main

import (
	"github.com/ashah360/nyte-auth/internal/captcha"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/ashah360/nyte-auth/internal/api/app"
	"github.com/ashah360/nyte-auth/internal/api/cache"
	"github.com/ashah360/nyte-auth/internal/api/repository"
	"github.com/ashah360/nyte-auth/internal/api/service"
)

var a app.Application

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg := app.NewConfig()
	cfg.RequireRecaptchaToRegister = false

	// DB
	db := &repository.DBInfo{
		Name:     os.Getenv("SQL_DBNAME"),
		Host:     os.Getenv("SQL_HOSTNAME"),
		Port:     os.Getenv("SQL_PORT"),
		User:     os.Getenv("SQL_USER"),
		Password: os.Getenv("SQL_PASSWORD"),
	}
	conn, err := db.Connect("postgres")
	if err != nil {
		log.Fatal(err)
	}
	urepo := repository.NewUserRepository(conn.Unsafe())

	// Token Store
	rds := cache.NewClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0)
	ts := cache.NewTokenStore(rds)

	// Service
	srv := service.NewAuthService(urepo, ts)
	us := service.NewUserService(urepo)
	rs := captcha.NewRecaptchaService(os.Getenv("RECAP_VERIFY_URL"), os.Getenv("RECAP_CLIENT_KEY"), os.Getenv("RECAP_SERVER_KEY"))

	a = app.NewApplication(srv, us, rs, cfg)
}
