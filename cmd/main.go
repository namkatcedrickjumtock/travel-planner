package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/joho/godotenv"
	"github.com/namkatcedrickjumtock/travel-planner/api"
	"github.com/namkatcedrickjumtock/travel-planner/internal/services"
	"github.com/namkatcedrickjumtock/travel-planner/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		// Use fmt.Fprintf so the formatted message is printed correctly.
		fmt.Fprintf(os.Stderr, "error starting the application: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg struct {
		API struct {
			ListenPort string `conf:"env:LISTEN_PORT,required"`
		}
		DB struct {
			User           string `conf:"env:DB_USER,mask,required"`
			Password       string `conf:"env:DB_PASSWORD,mask,required"`
			Host           string `conf:"env:DB_HOST,required"`
			Port           int    `conf:"env:DB_PORT,required"`
			Name           string `conf:"env:DB_NAME,required"`
			AllowedOrigins string `conf:"env:ALLOWED_ORIGINS,required"`
			MigrationsPath string `conf:"env:DB_MIGRATIONS_PATH,required"`
		}
	}

	// Load .env file when present (development convenience).
	// In production the variables should be injected directly into the environment.
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file")
		}
	}

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Printf("%v\n", help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	// Build the PostgreSQL DSN from config values.
	sslMode := "sslmode=disable"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s %s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, sslMode,
	)

	// Open the standard library DB first so we can pass it to the migrator.
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("opening sql connection: %w", err)
	}
	defer sqlDB.Close()

	// Wrap the same connection in GORM for the ORM layer.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return fmt.Errorf("opening gorm connection: %w", err)
	}

	// Run any pending SQL migrations before accepting traffic.
	if err := persistence.Migrate(sqlDB, cfg.DB.MigrationsPath, cfg.DB.Name); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	// Wire up the dependency chain: persistence → service → api.
	repo, err := persistence.NewRepository(gormDB)
	if err != nil {
		return fmt.Errorf("creating repository: %w", err)
	}

	svc, err := services.NewTravelPlannerService(repo)
	if err != nil {
		return fmt.Errorf("creating service: %w", err)
	}

	listener, err := api.NewAPIListener(svc)
	if err != nil {
		return fmt.Errorf("creating api listener: %w", err)
	}

	listenAddress := fmt.Sprintf("0.0.0.0:%s", cfg.API.ListenPort)
	log.Printf("server listening on %s", listenAddress)

	return listener.Run(listenAddress)
}
