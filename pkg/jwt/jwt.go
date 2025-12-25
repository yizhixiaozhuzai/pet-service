package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

// GenerateToken 生成JWT token
func (m *JWTManager) GenerateToken(userID uint, username string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(m.tokenDuration)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "pet-service",
			Subject:   "user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(m.tokenDuration.Seconds()), nil
}

// ValidateToken 验证JWT token
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的token")
	}

	return claims, nil
}

// RefreshToken 刷新token
func (m *JWTManager) RefreshToken(tokenString string) (string, int64, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", 0, err
	}

	// 检查token是否即将过期（30分钟内）
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", 0, errors.New("token还未到刷新时间")
	}

	return m.GenerateToken(claims.UserID, claims.Username)
}
