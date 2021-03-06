package model

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

// User represents a registered user.
type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}

func (c *UserModel) TableName() string {
	return "tb_users"
}

// Create creates a new user account.
func (u *UserModel) Create() error {
	return DB.Write.Create(&u).Error
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(id int64) error {
	user := UserModel{}
	user.Id = id
	return DB.Write.Unscoped().Delete(&user).Error
}

// Update updates an user account information.
func (u *UserModel) Update() error {
	return DB.Write.Save(u).Error
}

// GetUser gets an user by the user identifier.
func GetUser(username string) (*UserModel, error) {

	u := &UserModel{}
	d := DB.Read.Where("username = ?", username).First(u)
	return u, d.Error
}

// ListUser List all users
func ListUser(username string, offset, limit int64) ([]*UserModel, uint64, error) {
	// if limit == 0 {
	// 	limit = constvar.DefaultLimit
	// }

	users := make([]*UserModel, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Read.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Read.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

//Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
