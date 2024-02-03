package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tasuke/go-auth/models"
	"github.com/tasuke/go-auth/utils/token"
	"net/http"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	// JSONのバインディングを試み、エラーがあればレスポンスを返す
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// User構造体のインスタンスを作成し、入力値を設定
	u := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// ユーザーを保存し、エラーがあればレスポンスを返す
	if _, err := u.SaveUser(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	// 登録成功のメッセージを返す
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	var input LoginInput

	// リクエストボディのバインドを試みる
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// ユーザー認証を試みる
	token, err := models.LoginCheck(input.Email, input.Password)
	if err != nil {
		// 認証失敗のエラーメッセージは、情報漏洩を避けるために一般的にする
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 認証成功、トークンを返す
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CurrentUser(c *gin.Context) {
	// トークンからユーザーIDを抽出
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		fmt.Println("current")
		// トークン関連のエラーは、不正なトークンやトークンの解析に関連する可能性があるため、
		// HTTPステータスコードとしては401 Unauthorizedを使用するのが適切
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// ユーザー情報を取得
	u, err := models.GetUserByID(userId)
	if err != nil {
		// ユーザーが見つからない場合やデータベース関連のエラーには、
		// 404 NotFoundや500 Internal Server Errorなど、エラーの種類に応じた
		// ステータスコードを設定することを検討する
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{"message": "User retrieved successfully", "data": u})
}
