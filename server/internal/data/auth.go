package data

import "gorm.io/gorm"

type AppUser struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique;index;type:varchar(255)" json:"username"`
	Password string `json:"password"`
	Openid   string `json:"openid"`
}

type AuthRepo struct {
	mysqlDB *gorm.DB
}

func NewAuthRepo(mysqlDB *gorm.DB) AuthRepo {
	return AuthRepo{mysqlDB: mysqlDB}
}

func (r AuthRepo) RegisterAppUser(username, password string) error {
	return r.mysqlDB.Create(&AppUser{
		Username: username,
		Password: password,
	}).Error
}
func (r AuthRepo) VerifyLogin(username, password string) error {
	return r.mysqlDB.Where("username = ? AND password = ?", username, password).First(&AppUser{}).Error
}
