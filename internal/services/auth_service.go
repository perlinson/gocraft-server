package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perlinson/gocraft-server/internal/proto/auth"
	Store "github.com/perlinson/gocraft-server/internal/store"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer
	mu       sync.RWMutex
	sessions map[string]*UserSession
	store    *Store.Store
	jwtKey   []byte // 添加 JWT 密钥
}

type UserSession struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	User    *User  `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 修改 NewAuthService 方法以接受 Store 作为参数
func NewAuthService(store *Store.Store) *AuthService {
	return &AuthService{
		sessions: make(map[string]*UserSession),
		store:    store,
	}
}

// 生成随机令牌
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// 登录实现
func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// TODO: 实际用户验证逻辑
	userID := "user123" // 示例用户ID

	// 生成会话令牌
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// 创建会话
	session := &UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()

	expiresTime := time.Now().Add(24 * time.Hour)
	user := &auth.User{
		Id:   userID,
		Name: "测试用户",
	}
	return &auth.LoginResponse{
		Token:   token,
		Expires: expiresTime.Unix(),
		User:    user,
	}, nil
}

// 登出实现
func (s *AuthService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	s.mu.Lock()
	delete(s.sessions, req.Token)
	s.mu.Unlock()

	return &auth.LogoutResponse{}, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 1. 检查用户名是否已存在
	exists, err := s.store.UserExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, err
	}

	// 2. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. 创建用户
	userID := uuid.New().String()

	// 4. 保存到数据库
	if err := s.store.CreateUser(ctx, req.Username, string(hashedPassword), req.Email); err != nil {
		return nil, err
	}

	// 5. 生成 token
	token := uuid.New().String()
	expires := time.Now().Add(24 * time.Hour).Unix()

	// 6. 保存会话
	s.mu.Lock()
	s.sessions[token] = &UserSession{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.mu.Unlock()

	user := &User{
		ID:   userID,
		Name: req.Username,
	}
	return &RegisterResponse{
		Token:   token,
		Expires: expires,
		User:    user,
	}, nil
}

// RegisterRoutes 注册 HTTP 路由
func (s *AuthService) RegisterRoutes(r *gin.Engine) {
	// 登录路由
	r.POST("/api/auth/login", s.httpLogin)
	// 注册路由
	r.POST("/api/auth/register", s.httpRegister)
}

// httpLogin 处理登录请求
func (s *AuthService) httpLogin(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := s.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// httpRegister 处理注册请求
func (s *AuthService) httpRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := s.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "register failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
