package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Portfolio struct {
	UID         string   `gorm:"primaryKey;type:varchar(255)"`
	Openid      string   `gorm:"index;type:varchar(255)"`
	Title       string   `gorm:"type:varchar(255)"`
	Works       []Work   `gorm:"foreignKey:PortfolioUID;references:UID"`
	TemplateUID string   `gorm:"index;type:varchar(255)" json:"template_uid"`
	Template    Template `gorm:"foreignKey:TemplateUID;references:UID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Work struct {
	OSSKey       string `gorm:"primaryKey;type:varchar(255)"`
	PortfolioUID string `gorm:"type:varchar(255)" json:"portfolio_uid"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Template struct {
	UID       string `gorm:"primaryKey;type:varchar(255)"`
	Name      string `gorm:"type:varchar(255)"`
	OSSKey    string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
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
	result := r.redisClient.ZRangeWithScores(ctx, "zTemplates", 0, 5)
	if result.Err() != nil {
		return nil, result.Err()
	}
	templates := []Template{}
	zresults, err := result.Result()
	if err != nil {
		return nil, err
	}
	for _, zresult := range zresults {
		templates = append(templates, Template{OSSKey: zresult.Member.(string)})
	}

	return templates, nil
}

func (r PortfolioRepo) IncreTemplateScore(ctx context.Context, uid string) error {
	result := r.redisClient.ZIncrBy(ctx, "zTemplates", 1, uid)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (r PortfolioRepo) GetPortfolioFromDB(ctx context.Context, openid string) ([]Portfolio, error) {
	return nil, nil
}
func (r PortfolioRepo) GetPortfolioFromRedis(ctx context.Context, openid string) ([]Portfolio, error) {
	return nil, nil
}

func (r PortfolioRepo) SavePortfolioToDB(ctx context.Context, portfolio Portfolio, works []Work, openid string) error {
	err := r.mysqlDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uid"}},
			UpdateAll: true,
		}).Create(&portfolio).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "oss_key"}},
			UpdateAll: true,
		}).Create(&works).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
