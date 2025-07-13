package main

import (
	"database/sql"
	"log"

	_ "github.com/rizky-ardiansah/event-api/docs"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rizky-ardiansah/event-api/internal/database"
	"github.com/rizky-ardiansah/event-api/internal/env"
)

// title: Event API
// @version 1.0
// @description ini adalah simple API untuk mengelola event
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description masukan token JWT di header Authorization dengan format **Bearer &lt;token&gt;**


type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
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
