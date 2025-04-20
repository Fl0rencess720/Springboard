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
	ID          uint     `gorm:"primarykey"`
	UID         string   `gorm:"unique;index;type:varchar(255)" json:"uid"`
	Openid      string   `gorm:"index;type:varchar(255)" json:"openid"`
	Title       string   `gorm:"type:varchar(255)" json:"title"`
	Works       []Work   `gorm:"foreignKey:PortfolioUID;references:UID" json:"works"`
	TemplateUID string   `gorm:"index;type:varchar(255)" json:"template_uid"`
	Template    Template `gorm:"foreignKey:TemplateUID;references:UID" json:"template"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Work struct {
	ID           uint   `gorm:"primarykey"`
	OSSKey       string `gorm:"unique;index;type:varchar(255)" json:"oss_key"`
	PortfolioUID string `gorm:"type:varchar(255)" json:"portfolio_uid"`
	// Size 格式为 axb 例如 1920x1080
	Size       string `gorm:"type:varchar(255)" json:"size"`
	MarginTop  string `gorm:"type:varchar(255)" json:"margin_top"`
	MarginLeft string `gorm:"type:varchar(255)" json:"margin_left"`
	// 出血线，4个字符分别代表上、左、下、右的裁剪位
	Bleed     []string `gorm:"type:json;serializer:json" json:"bleed"`
	Page      int      `gorm:"type:varchar(255)" json:"page"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Template struct {
	ID        uint   `gorm:"primarykey"`
	UID       string `gorm:"unique;index;type:varchar(255)" json:"uid"`
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
	data, err := r.redisClient.Get(ctx, "templates").Bytes()
	if err != nil {
		return nil, err
	}
	var templates []Template
	if err = json.Unmarshal(data, &templates); err != nil {
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
	portfolios := []Portfolio{}
	if err := r.mysqlDB.Preload("Works").Preload("Template").Where("openid = ?", openid).Find(&portfolios).Error; err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (r PortfolioRepo) GetPortfolioFromRedis(ctx context.Context, openid string) ([]Portfolio, error) {
	data, err := r.redisClient.Get(ctx, "portfolios:"+openid).Bytes()
	if err != nil {
		return nil, err
	}
	var portfolios []Portfolio
	if err = json.Unmarshal(data, &portfolios); err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (r PortfolioRepo) SavePortfoliosToRedis(ctx context.Context, portfolios []Portfolio, openid string) error {
	portfoliosJson, err := json.Marshal(portfolios)
	if err != nil {
		return err
	}
	if err := r.redisClient.Set(ctx, "portfolios:"+openid, portfoliosJson, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r PortfolioRepo) SavePortfolioToDB(ctx context.Context, portfolio Portfolio, works []Work) error {
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
