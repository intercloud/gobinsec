package gobinsec

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteConfig is the configuration for SQLite
type SQLiteConfig struct {
	File       string `yaml:"file"`
	Expiration int32  `yaml:"expiration"`
}

// NewSQLiteConfig builds configuration for SQLite cache
func NewSQLiteConfig(config *SQLiteConfig) *SQLiteConfig {
	if config == nil {
		config = &SQLiteConfig{}
	}
	if config.File == "" {
		config.File = "~/.gobinsec.db"
	}
	if strings.HasPrefix(config.File, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		config.File = filepath.Join(home, config.File[2:])
	}
	if config.Expiration == 0 {
		config.Expiration = 86400
	}
	return config
}

type SQLiteCache struct {
	Expiration int32
	Database   *sql.DB
}

// NewSQLiteCache builds a cache using SQLite
func NewSQLiteCache(config *SQLiteConfig) Cache {
	database, err := sql.Open("sqlite3", config.File)
	if err != nil {
		return nil
	}
	cache := &SQLiteCache{
		Database:   database,
		Expiration: config.Expiration,
	}
	return cache
}

// Get returns NVD response for given dependency
func (sc *SQLiteCache) Get(d *Dependency) []byte {
	vulnerabilities, err := sc.RetrieveDependency(d)
	if err != nil {
		fmt.Printf("ERROR getting dependency: %v\n", err)
	}
	return vulnerabilities
}

// Set put NVD response for given dependency in cache
func (sc *SQLiteCache) Set(d *Dependency, v []byte) {
	err := sc.InsertDependency(d, v)
	if err != nil {
		fmt.Printf("ERROR setting dependency: %v\n", err)
	}
}

// Ping does nothing
func (sc *SQLiteCache) Ping() error {
	return sc.CreateTable()
}

// Clean deletes expired entries
func (sc *SQLiteCache) Clean() {
	err := sc.CleanTable()
	if err != nil {
		fmt.Printf("ERROR cleaning dependencies: %v\n", err)
	}
}

// CreateTable creates dependencies table
func (sc *SQLiteCache) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS dependencies (
		dependency TEXT,
		date TEXT,
		vulnerabilities TEXT
	)`
	_, err := sc.Database.Exec(query)
	return err
}

// InsertDependency insert given dependency in database
func (sc *SQLiteCache) InsertDependency(d *Dependency, v []byte) error {
	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	query := `INSERT INTO dependencies
		(dependency, date, vulnerabilities)
		VALUES
		(?, ?, ?)`
	_, err := sc.Database.Exec(query, d.Key(), now, string(v))
	return err
}

// RetrieveDependency insert given dependency in database
func (sc *SQLiteCache) RetrieveDependency(d *Dependency) ([]byte, error) {
	duration := time.Duration(sc.Expiration) * time.Second
	limit := time.Now().UTC().Add(-duration).Format("2006-01-02T15:04:05")
	query := `SELECT vulnerabilities
		FROM dependencies
		WHERE dependency = ?
		AND date > ?`
	row, err := sc.Database.Query(query, d.Key(), limit)
	if err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}
	defer row.Close()
	if row.Next() {
		var vulnerabilities string
		if err := row.Scan(&vulnerabilities); err != nil {
			return nil, err
		}
		return []byte(vulnerabilities), nil
	}
	return nil, nil
}

// CleanTable deletes old expired entries
func (sc *SQLiteCache) CleanTable() error {
	duration := time.Duration(sc.Expiration) * time.Second
	limit := time.Now().UTC().Add(-duration).Format("2006-01-02T15:04:05")
	query := `DELETE
		FROM dependencies
		WHERE date < ?`
	_, err := sc.Database.Exec(query, limit)
	return err
}
