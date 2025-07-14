package main

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/rizky-ardiansah/event-api/docs"
	"github.com/rizky-ardiansah/event-api/internal/database"
	"github.com/rizky-ardiansah/event-api/internal/env"
)

// @title Event API
// @version 1.0
// @description A rest API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	dbURL := env.GetEnvString("DATABASE_URL", " ")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "defaultsecret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
