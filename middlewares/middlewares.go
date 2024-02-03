package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tasuke/go-auth/utils/token"
	"net/http"
)

// JwtAuthMiddleware はJWTを用いた認証のミドルウェアを提供します。
// このミドルウェアは、リクエストのAuthorizationヘッダーに含まれるトークンが有効かどうかを検証します。
// トークンが無効な場合、リクエストはUnauthorizedエラーで中断されます。
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token.TokenValid関数を用いてトークンの検証を行う
		err := token.TokenValid(c)
		if err != nil {
			// トークン検証に失敗した場合、HTTPステータスコード401とともにUnauthorizedエラーメッセージを返す
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // リクエスト処理を中断し、以降のハンドラーが実行されないようにする
			return
		}
		c.Next() // トークン検証が成功した場合、次のハンドラーへ処理を渡す
	}
}
