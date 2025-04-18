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
	UID         string   `gorm:"primaryKey;type:varchar(255)" json:"uid"`
	Openid      string   `gorm:"index;type:varchar(255)" json:"openid"`
	Title       string   `gorm:"type:varchar(255)" json:"title"`
	Works       []Work   `gorm:"foreignKey:PortfolioUID;references:UID" json:"works"`
	TemplateUID string   `gorm:"index;type:varchar(255)" json:"template_uid"`
	Template    Template `gorm:"foreignKey:TemplateUID;references:UID" json:"template"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Work struct {
	OSSKey       string `gorm:"primaryKey;type:varchar(255)" json:"oss_key"`
	PortfolioUID string `gorm:"type:varchar(255)" json:"portfolio_uid"`
	// Size 格式为 axb 例如 1920x1080
	Size       string `gorm:"type:varchar(255)" json:"size"`
	MarginTop  string `gorm:"type:varchar(255)" json:"margin_top"`
	MarginLeft string `gorm:"type:varchar(255)" json:"margin_left"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Template struct {
	UID       string `gorm:"primaryKey;type:varchar(255)" json:"uid"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	OSSKey    string `gorm:"type:varchar(255)" json:"oss_key"`
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

func (r PortfolioRepo) GetTemplatesFromDB(ctx context.Context, uids []string) ([]Template, error) {
	templates := []Template{}
	if err := r.mysqlDB.Where("uid IN ?", uids).Find(&templates).Error; err != nil {
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
	//目前查询最热的5个
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
		templates = append(templates, Template{UID: zresult.Member.(string)})
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
	portfolio := []Portfolio{}
	if err := r.mysqlDB.Where("uid IN ?", openid).Find(&portfolio).Error; err != nil {
		return nil, err
	}
	return portfolio, nil
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
