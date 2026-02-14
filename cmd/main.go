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
	"github.com/namkatcedrickjumtock/travel-planner/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Failed:", err)
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
			DisableTLS     bool   `conf:"env:DB_DISABLE_TLS,default:false"`
			AllowedOrigins string `conf:"env:ALLOWED_ORIGINS,required"`
			MigrationsPath string `conf:"env:DB_MIGRATIONS_PATH,required"`
		}
	}

	// loadDevEnv loads .env file if present
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
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

	tslConfig := ""
	if cfg.DB.DisableTLS {
		tslConfig = "sslmode=disable"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, tslConfig)

	sqlDBInstance, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDBInstance,
	}), &gorm.Config{TranslateError: true})
	if err != nil {
		return err
	}

	defer sqlDBInstance.Close()

	// ctx := context.Background()

	// if err := persistence.Migrate(sqlDBInstance, cfg.DB.MigrationsPath, cfg.DB.Name); err != nil {
	// 	return err
	// }

	repo, err := persistence.NewRepository(gormDB)
	if err != nil {
		return err
	}

	listener, err := api.NewAPIListener(repo)
	if err != nil {
		return err
	}

	listenAddress := fmt.Sprintf("0.0.0.0:%s", cfg.API.ListenPort)

	return listener.Run(listenAddress)
}
