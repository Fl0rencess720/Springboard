package data

import (
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Portfolio struct {
	UID         string   `gorm:"primaryKey;type:varchar(255)"`
	Openid      string   `gorm:"index;type:varchar(255)"`
	Title       string   `gorm:"type:varchar(255)"`
	Works       []Work   `gorm:"foreignKey:PortfolioUID"`
	TemplateUID string   `gorm:"index;type:varchar(255)"`
	Template    Template `gorm:"foreignKey:TemplateUID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time
}

type Work struct {
	gorm.Model
	OSSKey       string    `gorm:"type:varchar(255)"`
	PortfolioUID string    `gorm:"index;type:varchar(255)"`
	Portfolio    Portfolio `gorm:"foreignKey:PortfolioUID;references:UID"`
}

type Template struct {
	UID        string      `gorm:"primaryKey;type:varchar(255)"`
	Name       string      `gorm:"type:varchar(255)"`
	OSSKey     string      `gorm:"type:varchar(255)"`
	Portfolios []Portfolio `gorm:"foreignKey:TemplateUID"`
	CreatedAt  time.Time
}

type PortfolioRepo struct {
	mysqlDB     *gorm.DB
	redisClient *redis.Client
}

func NewPortfolioRepo(mysqlDB *gorm.DB, redisClient *redis.Client) PortfolioRepo {
	return PortfolioRepo{
		mysqlDB:     mysqlDB,
		redisClient: redisClient,
	}
}

func (r PortfolioRepo) GetAllTemplatesFromDB() []Template {
	return nil
}
func (r PortfolioRepo) GetAllTemplatesFromRedis() []Template {
	return nil
}
func (r PortfolioRepo) GetHotTemplatesFromDB() []Template {
	return nil
}
func (r PortfolioRepo) GetHotTemplatesFromRedis() []Template {
	return nil
}
func (r PortfolioRepo) GetPortfolioFromDB(openid string) []Portfolio {
	return nil
}
func (r PortfolioRepo) GetPortfolioFromRedis(openid string) []Portfolio {
	return nil
}
