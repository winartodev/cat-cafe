package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	postgresDNS = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
)

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
}

func (d *Database) SetupConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf(postgresDNS, d.Username, d.Password, d.Host, d.Port, d.Name, d.SSLMode)

	db, err := sql.Open(d.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error open db connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
