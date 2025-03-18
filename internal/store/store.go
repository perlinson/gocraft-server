package store

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"context"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var (
	dbpath = flag.String("db", "", "db file name (legacy, not used with MySQL)")
)

// var (
// 	store *Store
// )

func InitStore() (*Store, error) {
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

	// Open MySQL connection using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	store := &Store{DB: db}
	// 确保 DB 被正确赋值
	if store.DB == nil {
		return nil, fmt.Errorf("DB 未初始化")
	}
	// Initialize database tables
	err = store.initTables()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %v", err)
	}

	return store, nil
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Store ...
type Store struct {
	DB *gorm.DB
}

type Block struct {
	ChunkX int32 `gorm:"column:chunk_x"`
	ChunkZ int32 `gorm:"column:chunk_z"`
	BlockX int32 `gorm:"column:block_x"`
	BlockY int32 `gorm:"column:block_y"`
	BlockZ int32 `gorm:"column:block_z"`
	BlockType int32 `gorm:"column:block_type"`
}

type Chunk struct {
	ChunkX int32 `gorm:"column:chunk_x"`
	ChunkY int32 `gorm:"column:chunk_y"`
	ChunkZ int32 `gorm:"column:chunk_z"`
	Version string `gorm:"column:version"`
}

type Camera struct {
	ID int32 `gorm:"column:id"`
	X float32 `gorm:"column:x"`
	Y float32 `gorm:"column:y"`
	Z float32 `gorm:"column:z"`
	RX float32 `gorm:"column:rx"`
	RY float32 `gorm:"column:ry"`
}

type User struct {
	ID int32 `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Email string `gorm:"column:email"`
}

func (s *Store) initTables() error {
	// Create blocks table
	err := s.DB.AutoMigrate(&Block{})
	if err != nil {
		return err
	}

	// Create chunks table for version tracking
	err = s.DB.AutoMigrate(&Chunk{})
	if err != nil {
		return err
	}

	// Create camera table
	err = s.DB.AutoMigrate(&Camera{})
	if err != nil {
		return err
	}

	// Insert default camera if not exists
	var count int64
	err = s.DB.Model(&Camera{}).Where("id = ?", 1).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		camera := Camera{ID: 1, X: 0, Y: 16, Z: 0, RX: 0, RY: 0}
		err = s.DB.Create(&camera).Error
		if err != nil {
			return err
		}
	}

	// Create users table
	err = s.DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return nil
}

// UserExists 检查用户名是否已存在
func (s *Store) UserExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := s.DB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CreateUser 创建新用户
func (s *Store) CreateUser(ctx context.Context, username, password, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := User{Username: username, Password: string(hashedPassword), Email: email}
	return s.DB.Create(&user).Error
}

func (s *Store) UpdateBlock(id Vec3, w int) error {
	// Get chunk coordinates
	cid := id.Chunkid()

	// Log the update
	log.Printf("put %v -> %d", id, w)

	// Begin transaction
	tx := s.DB.Begin()
	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Insert or update block
	block := Block{ChunkX: cid.X, ChunkZ: cid.Z, BlockX: id.X, BlockY: id.Y, BlockZ: id.Z, BlockType: int32(w)}
	err := tx.Save(&block).Error
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

func (s *Store) UpdateCamera(x, y, z, rx, ry float32) error {
	camera := Camera{ID: 1, X: x, Y: y, Z: z, RX: rx, RY: ry}
	return s.DB.Save(&camera).Error
}

func (s *Store) GetCamera() (x, y, z, rx, ry float32) {
	// Default value if query fails
	y = 16

	var camera Camera
	err := s.DB.Where("id = ?", 1).First(&camera).Error
	if err != nil {
		log.Printf("Error getting camera: %v", err)
	} else {
		x = camera.X
		y = camera.Y
		z = camera.Z
		rx = camera.RX
		ry = camera.RY
	}

	return
}

func (s *Store) RangeBlocks(id Vec3, f func(bid Vec3, w int)) error {
	var blocks []Block
	err := s.DB.Where("chunk_x = ? AND chunk_z = ?", id.X, id.Z).Find(&blocks).Error
	if err != nil {
		return err
	}

	for _, block := range blocks {
		f(Vec3{block.BlockX, block.BlockY, block.BlockZ}, int(block.BlockType))
	}

	return nil
}

func (s *Store) UpdateChunkVersion(id Vec3, version string) error {
	chunk := Chunk{ChunkX: id.X, ChunkY: id.Y, ChunkZ: id.Z, Version: version}
	return s.DB.Save(&chunk).Error
}

func (s *Store) GetChunkVersion(id Vec3) string {
	var chunk Chunk
	err := s.DB.Where("chunk_x = ? AND chunk_y = ? AND chunk_z = ?", id.X, id.Y, id.Z).First(&chunk).Error
	if err != nil {
		// If no version found, return empty string
		if err == gorm.ErrRecordNotFound {
			return ""
		}
		log.Printf("Error getting chunk version: %v", err)
		return ""
	}

	return chunk.Version
}

func (s *Store) Close() {
	if s.DB != nil {
		db, err := s.DB.DB()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			db.Close()
		}
	}
}

func GenerateChunkVersion() string {
	return strconv.FormatInt(time.Now().UnixNano(), 16)
}

const (
	ChunkWidth = 32
)

type Vec3 struct {
	X, Y, Z int32
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
