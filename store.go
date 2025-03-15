package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	dbpath = flag.String("db", "", "db file name (legacy, not used with MySQL)")
)

var (
	store *Store
)

func InitStore() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get MySQL connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "app_db")
	dbCharset := getEnv("DB_CHARSET", "utf8mb4")
	dbParseTime := getEnv("DB_PARSE_TIME", "True")
	dbLoc := getEnv("DB_LOC", "Local")

	// Build MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset, dbParseTime, dbLoc)

	log.Printf("Connecting to MySQL database at %s:%s/%s", dbHost, dbPort, dbName)

	store, err = NewStore(dsn)
	return err
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Store struct {
	db *sql.DB
}

func NewStore(dsn string) (*Store, error) {
	// Open MySQL connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create the store instance
	store := &Store{
		db: db,
	}

	// Initialize database tables
	err = store.initTables()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %v", err)
	}

	return store, nil
}

func (s *Store) initTables() error {
	// Create blocks table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS blocks (
			chunk_x INT NOT NULL,
			chunk_z INT NOT NULL,
			block_x INT NOT NULL,
			block_y INT NOT NULL,
			block_z INT NOT NULL,
			block_type INT NOT NULL,
			PRIMARY KEY (chunk_x, chunk_z, block_x, block_y, block_z)
		)
	`)
	if err != nil {
		return err
	}

	// Create chunks table for version tracking
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS chunks (
			chunk_x INT NOT NULL,
			chunk_y INT NOT NULL,
			chunk_z INT NOT NULL,
			version VARCHAR(32) NOT NULL,
			PRIMARY KEY (chunk_x, chunk_y, chunk_z)
		)
	`)
	if err != nil {
		return err
	}

	// Create camera table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS camera (
			id INT NOT NULL DEFAULT 1,
			x FLOAT NOT NULL,
			y FLOAT NOT NULL,
			z FLOAT NOT NULL,
			rx FLOAT NOT NULL,
			ry FLOAT NOT NULL,
			PRIMARY KEY (id)
		)
	`)
	if err != nil {
		return err
	}

	// Insert default camera if not exists
	_, err = s.db.Exec(`
		INSERT IGNORE INTO camera (id, x, y, z, rx, ry) VALUES (1, 0, 16, 0, 0, 0)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateBlock(id Vec3, w int) error {
	// Get chunk coordinates
	cid := id.Chunkid()

	// Log the update
	log.Printf("put %v -> %d", id, w)

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert or update block
	_, err = tx.Exec(
		"REPLACE INTO blocks (chunk_x, chunk_z, block_x, block_y, block_z, block_type) VALUES (?, ?, ?, ?, ?, ?)",
		cid.X, cid.Z, id.X, id.Y, id.Z, w,
	)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (s *Store) UpdateCamera(x, y, z, rx, ry float32) error {
	_, err := s.db.Exec(
		"UPDATE camera SET x = ?, y = ?, z = ?, rx = ?, ry = ? WHERE id = 1",
		x, y, z, rx, ry,
	)
	return err
}

func (s *Store) GetCamera() (x, y, z, rx, ry float32) {
	// Default value if query fails
	y = 16

	err := s.db.QueryRow("SELECT x, y, z, rx, ry FROM camera WHERE id = 1").Scan(&x, &y, &z, &rx, &ry)
	if err != nil {
		log.Printf("Error getting camera: %v", err)
	}

	return
}

func (s *Store) RangeBlocks(id Vec3, f func(bid Vec3, w int)) error {
	rows, err := s.db.Query(
		"SELECT block_x, block_y, block_z, block_type FROM blocks WHERE chunk_x = ? AND chunk_z = ?",
		id.X, id.Z,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var x, y, z, w int
		err := rows.Scan(&x, &y, &z, &w)
		if err != nil {
			return err
		}
		f(Vec3{x, y, z}, w)
	}

	return rows.Err()
}

func (s *Store) UpdateChunkVersion(id Vec3, version string) error {
	_, err := s.db.Exec(
		"REPLACE INTO chunks (chunk_x, chunk_y, chunk_z, version) VALUES (?, ?, ?, ?)",
		id.X, id.Y, id.Z, version,
	)
	return err
}

func (s *Store) GetChunkVersion(id Vec3) string {
	var version string
	err := s.db.QueryRow(
		"SELECT version FROM chunks WHERE chunk_x = ? AND chunk_y = ? AND chunk_z = ?",
		id.X, id.Y, id.Z,
	).Scan(&version)

	if err != nil {
		// If no version found, return empty string
		if err == sql.ErrNoRows {
			return ""
		}
		log.Printf("Error getting chunk version: %v", err)
		return ""
	}

	return version
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func GenrateChunkVersion() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}

const (
	ChunkWidth = 32
)

type Vec3 struct {
	X, Y, Z int
}

func (v Vec3) Left() Vec3 {
	return Vec3{v.X - 1, v.Y, v.Z}
}

func (v Vec3) Right() Vec3 {
	return Vec3{v.X + 1, v.Y, v.Z}
}

func (v Vec3) Up() Vec3 {
	return Vec3{v.X, v.Y + 1, v.Z}
}

func (v Vec3) Down() Vec3 {
	return Vec3{v.X, v.Y - 1, v.Z}
}

func (v Vec3) Front() Vec3 {
	return Vec3{v.X, v.Y, v.Z + 1}
}

func (v Vec3) Back() Vec3 {
	return Vec3{v.X, v.Y, v.Z - 1}
}

func (v Vec3) Chunkid() Vec3 {
	x := v.X
	z := v.Z
	if x < 0 {
		x = x - ChunkWidth + 1
	}
	if z < 0 {
		z = z - ChunkWidth + 1
	}
	return Vec3{x / ChunkWidth, 0, z / ChunkWidth}
}
