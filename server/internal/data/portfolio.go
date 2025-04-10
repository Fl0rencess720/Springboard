package data

import (
	"context"
	"encoding/json"
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

func (r PortfolioRepo) GetAllTemplatesFromDB(ctx context.Context) ([]Template, error) {
	templates := []Template{}
	if err := r.mysqlDB.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
func (r PortfolioRepo) GetAllTemplatesFromRedis(ctx context.Context) ([]Template, error) {
	result := r.redisClient.Get(ctx, "templates")
	if result.Err() != nil {
		return nil, result.Err()
	}
	var templates []Template
	if err := result.Scan(&templates); err != nil {
		return nil, err
	}
	return templates, nil

}

func (r PortfolioRepo) SaveAllTemplatesToRedis(ctx context.Context, templates []Template) error {
	templatesJson, err := json.Marshal(templates)
	if err != nil {
		return err
	}
	if err := r.redisClient.Set(ctx, "templates", templatesJson, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r PortfolioRepo) GetHotTemplatesFromDB(ctx context.Context) ([]Template, error) {
	return nil, nil
}
func (r PortfolioRepo) GetHotTemplatesFromRedis(ctx context.Context) ([]Template, error) {
	return nil, nil
}
func (r PortfolioRepo) GetPortfolioFromDB(ctx context.Context, openid string) ([]Portfolio, error) {
	return nil, nil
}
func (r PortfolioRepo) GetPortfolioFromRedis(ctx context.Context, openid string) ([]Portfolio, error) {
	return nil, nil
}
