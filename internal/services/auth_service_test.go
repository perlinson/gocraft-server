package services_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/perlinson/gocraft-server/internal/services"
	"github.com/perlinson/gocraft-server/internal/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 测试认证服务
func TestAuthService(t *testing.T) {
	// 初始化 GORM 数据库连接
	// 添加更多的错误日志以帮助诊断问题
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		t.Fatalf("failed to load environment variables: %v", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_CHARSET"), os.Getenv("DB_PARSE_TIME"), os.Getenv("DB_LOC"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	store := &store.Store{DB: db}
	authService := services.NewAuthService(store)
	// 初始化 Gin 引擎
	router := gin.Default()
	authService.RegisterRoutes(router)

	// 测试注册
	t.Run("Register User", func(t *testing.T) {
		payload := `{"username":"testuser","password":"password123","email":"test@example.com"}`
		req, _ := http.NewRequest("POST", "/api/auth/register", strings.NewReader(payload))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	// 测试登录
	t.Run("Login User", func(t *testing.T) {
		payload := `{"username":"testuser","password":"password123"}`
		req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(payload))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// 测试玩家服务
func TestPlayerService(t *testing.T) {
	// TODO: 添加玩家服务的 gRPC 测试
}

// 测试方块服务
func TestBlockService(t *testing.T) {
	// TODO: 添加方块服务的 gRPC 测试
}
