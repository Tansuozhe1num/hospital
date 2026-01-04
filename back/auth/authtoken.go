package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthToken(token string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			authorization = strings.TrimSpace(authorization[7:])
		}
		if authorization == "" {
			authorization = request.URL.Query().Get("token")
		}
		if authorization != token {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func IssueToken(userID string, role string, ttl time.Duration) (string, error) {
	userID = strings.TrimSpace(userID)
	role = strings.TrimSpace(role)
	if userID == "" {
		return "", errors.New("user id cannot be empty")
	}
	if !isValidRole(role) {
		return "", errors.New("invalid role")
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	claims := Claims{
		UserID: userID,
		Role:   role,
		Exp:    time.Now().Add(ttl).Unix(),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	sig := sign(payload)
	return payload + "." + sig, nil
}

func ParseToken(token string) (*Claims, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}
	payload, sig := parts[0], parts[1]
	if payload == "" || sig == "" {
		return nil, errors.New("invalid token")
	}
	if !hmac.Equal([]byte(sig), []byte(sign(payload))) {
		return nil, errors.New("invalid token")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token")
	}
	if strings.TrimSpace(claims.UserID) == "" || !isValidRole(claims.Role) {
		return nil, errors.New("invalid token")
	}
	if claims.Exp <= 0 || time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func GinAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		r = strings.TrimSpace(r)
		if r != "" {
			allowed[r] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		token := bearerTokenFromRequest(c.Request)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if len(allowed) > 0 {
			if _, ok := allowed[claims.Role]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}

		c.Set("authClaims", claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) (*Claims, bool) {
	v, ok := c.Get("authClaims")
	if !ok {
		return nil, false
	}
	claims, ok := v.(*Claims)
	return claims, ok
}

func bearerTokenFromRequest(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return strings.TrimSpace(authorization[7:])
	}
	return strings.TrimSpace(authorization)
}

func sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(secret()))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func secret() string {
	s := strings.TrimSpace(os.Getenv("AUTH_SECRET"))
	if s != "" {
		return s
	}
	return "hospital-system-dev-secret"
}

func isValidRole(role string) bool {
	switch role {
	case "patient", "doctor", "admin":
		return true
	default:
		return false
	}
}
