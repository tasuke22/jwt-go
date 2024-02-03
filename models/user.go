package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/tasuke/go-auth/utils/token"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
	"time"
)

// User 構造体は、ユーザー情報を表します。
// IDフィールドをUUID形式に変更します。
type User struct {
	ID        string `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string `gorm:"size:255;not null;unique" json:"name"`
	Email     string `gorm:"size:100;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null;" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate はGORMのフックで、Userが保存される前に呼び出されます。
// これを使用して、UUIDを生成しIDフィールドに設定します。
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("ID", uuid.String())
}

// VerifyPassword は提供されたパスワードとハッシュ化されたパスワードを比較します。
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// SaveUser は新しいユーザーをデータベースに保存します。
// パスワードのハッシュ化を含む前処理を行った後、ユーザーを保存します。
func (u *User) SaveUser() (*User, error) {
	if err := u.beforeSave(); err != nil {
		return nil, err
	}

	if err := DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// beforeSave はSaveUserから呼び出され、
// ユーザーのパスワードをハッシュ化し、名前を適切にエスケープします。
func (u *User) beforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	return nil
}

// LoginCheck はユーザーのログインを検証します。
// 正しい認証情報が提供された場合、ユーザーに対してトークンを生成します。
func LoginCheck(email, password string) (string, error) {
	u := User{}
	if err := DB.Where("email = ?", email).Take(&u).Error; err != nil {
		return "", errors.New("login failed")
	}

	if err := VerifyPassword(password, u.Password); err != nil {
		return "", errors.New("invalid login credentials")
	}

	return token.GenerateToken(u.ID)
}

// GetUserByID は指定されたIDに基づいてユーザーを取得します。
// ユーザーが見つかった場合は、パスワード情報を除去して返します。
func GetUserByID(uid string) (User, error) {
	var u User
	// "id = ?"を使用して、UIDに基づいてユーザーを検索します。
	if err := DB.Where("id = ?", uid).First(&u).Error; err != nil {
		return u, err
	}

	u.PrepareGive()
	return u, nil
}

// PrepareGive はユーザー情報を返す前に、センシティブな情報を除去します。
func (u *User) PrepareGive() {
	u.Password = "" // パスワードを空にする
}
