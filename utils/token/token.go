package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"
)

var apiSecret = os.Getenv("API_SECRET")

// GenerateTokenはユーザーのIDを基に新しいJWTトークンを生成します。
// トークンの有効期限は環境変数`TOKEN_HOUR_LIFESPAN`で設定された時間です。
func GenerateToken(user_id string) (string, error) {
	// 環境変数からトークンの有効期限を読み込み、数値に変換
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", fmt.Errorf("failed to parse TOKEN_HOUR_LIFESPAN: %v", err)
	}

	// JWTクレームセットを作成
	claims := jwt.MapClaims{
		"authorized": true,                                                            // 認証済みフラグ
		"user_id":    user_id,                                                         // ユーザーID
		"exp":        time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix(), // 有効期限
	}

	// 新しいJWTトークンを作成し、クレームセットを設定
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 環境変数から読み込んだ秘密鍵でトークンに署名
	signedToken, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// ExtractTokenはリクエストからJWTトークンを抽出します。
// Authorizationヘッダー、もしくはクエリパラメータからトークンを取得します。
func ExtractToken(c *gin.Context) string {
	// AuthorizationヘッダーからBearerトークンを取得
	bearerToken := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(bearerToken, "Bearer ")
	if token != "" {
		return token
	}
	// ヘッダーにトークンがない場合、クエリパラメータからトークンを取得
	return c.Query("token")
}

// ParseTokenは与えられたJWTトークン文字列を解析し、トークンオブジェクトを返します。
// この関数はトークンの署名方法を確認し、設定されたAPI秘密鍵でトークンを検証します。
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
}

// TokenValidはリクエスト内のJWTトークンの有効性を検証します。
// トークンが無効である場合、エラーを返します。
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := ParseToken(tokenString)
	return err
}

// ExtractTokenIDはリクエストからJWTトークンを抽出し、そのトークン内のユーザーIDを取得します。
// トークンが無効、またはユーザーIDの抽出に失敗した場合、エラーを返します。
func ExtractTokenID(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("user_id is not a valid UUID")
		}
		// UUID形式のuserIDをそのまま返します。
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}
