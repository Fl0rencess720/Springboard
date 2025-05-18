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
	ID          uint      `gorm:"primarykey"`
	UID         string    `gorm:"unique;index;type:varchar(255)" json:"uid"`
	Openid      string    `gorm:"index;type:varchar(255)" json:"openid"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Projects    []Project `gorm:"foreignKey:PortfolioUID;references:UID" json:"projects"`
	TemplateUID string    `gorm:"index;type:varchar(255)" json:"template_uid"`
	Template    Template  `gorm:"foreignKey:TemplateUID;references:UID" json:"template"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Work struct {
	ID         uint   `gorm:"primarykey"`
	OSSKey     string `gorm:"unique;index;type:varchar(255)" json:"oss_key"`
	ProjectUID string `gorm:"type:varchar(255)" json:"project_uid"`
	// Size 格式为 axb 例如 1920x1080
	Size       string  `gorm:"type:varchar(255)" json:"size"`
	MarginTop  string  `gorm:"type:varchar(255)" json:"margin_top"`
	MarginLeft string  `gorm:"type:varchar(255)" json:"margin_left"`
	Scale      float64 `gorm:"type:double;default:1.0" json:"scale"` // 1.0 表示 不缩放
	PageNum    int     `gorm:"column:page;type:int" json:"page_num"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type Text struct {
	ID         uint   `gorm:"primarykey"`
	UID        string `gorm:"unique;index;type:varchar(255)" json:"uid"`
	ProjectUID string `gorm:"type:varchar(255)" json:"project_uid"`
	Content    string `gorm:"type:varchar(255)" json:"content"`
	FontSize   string `gorm:"type:varchar(255)" json:"font_size"`
	FontColor  string `gorm:"type:char(6);default:'000000'" json:"font_color"`
	Size       string `gorm:"type:varchar(255)" json:"size"` // 文本框大小
	MarginTop  string `gorm:"type:varchar(255)" json:"margin_top"`
	MarginLeft string `gorm:"type:varchar(255)" json:"margin_left"`
	PageNum    int    `gorm:"column:page;type:int" json:"page_num"`
}
type Template struct {
	ID         uint   `gorm:"primarykey"`
	UID        string `gorm:"unique;index;type:varchar(255)" json:"uid"`
	Name       string `gorm:"type:varchar(255)" json:"name"`
	FontOSSKey string `gorm:"type:varchar(255)" json:"font_oss_key"`
	Pages      []Page `gorm:"foreignKey:TemplateUID;references:UID" json:"pages"`
	CreatedAt  time.Time
}

// 模板中的固有页面
// oss key 格式：template_5_4.svg，其中5为模板id，4为页面在模板中的顺序
// margin, size表示矩形空位，用于放 work 的位置
type Page struct {
	ID            uint   `gorm:"primarykey"`
	UID           string `gorm:"unique;index;type:varchar(255)" json:"uid"`
	OSSKey        string `gorm:"unique;index;type:varchar(255)" json:"oss_key"`
	PreviewOSSKey string `gorm:"type:varchar(255)" json:"preview_oss_key"`
	// 出血线，4个字符分别代表svg的x、y、width、height
	Bleed         []string `gorm:"type:json;serializer:json" json:"bleed"`
	TemplateUID   string   `gorm:"type:varchar(255)" json:"template_uid"`
	MarginTop     string   `gorm:"type:varchar(255)" json:"margin_top"`
	MarginLeft    string   `gorm:"type:varchar(255)" json:"margin_left"`
	Size          string   `gorm:"type:varchar(255)" json:"size"`     // 图片容纳框大小
	BkgSize       string   `gorm:"type:varchar(255)" json:"bkg_size"` // 背景图大小
	IsContentPage bool     `gorm:"type:bool" json:"is_content_page"`
}
type Project struct {
	ID           uint   `gorm:"primarykey"`
	UID          string `gorm:"unique;index;type:varchar(255)" json:"uid"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	Order        int    `gorm:"type:int" json:"order"`
	PortfolioUID string `gorm:"type:varchar(255)" json:"portfolio_uid"`
	Works        []Work `gorm:"foreignKey:ProjectUID;references:UID" json:"works"`
	Texts        []Text `gorm:"foreignKey:ProjectUID;references:UID" json:"texts"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	if err := r.mysqlDB.Preload("Pages").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r PortfolioRepo) GetTemplatesFromDB(ctx context.Context, uids []string) ([]Template, error) {
	templates := []Template{}
	if err := r.mysqlDB.Preload("Pages").Where("uid IN ?", uids).Find(&templates).Error; err != nil {
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

func (r PortfolioRepo) GetTemplateByUIDFromDB(ctx context.Context, uid string) (Template, error) {
	template := Template{}
	if err := r.mysqlDB.Preload("Pages").Where("uid = ?", uid).First(&template).Error; err != nil {
		return Template{}, err
	}
	return template, nil
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

func (r PortfolioRepo) GetPortfoliosFromDB(ctx context.Context, openid string) ([]Portfolio, error) {
	portfolios := []Portfolio{}
	if err := r.mysqlDB.Preload("Projects.Works").Preload("Template").Where("openid = ?", openid).Find(&portfolios).Error; err != nil {
		return nil, err
	}
	return portfolios, nil
}

func (r PortfolioRepo) GetPortfoliosFromRedis(ctx context.Context, openid string) ([]Portfolio, error) {
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

func (r PortfolioRepo) GetPortfolioByUIDFromDB(ctx context.Context, uid string) (Portfolio, error) {
	portfolio := Portfolio{}
	if err := r.mysqlDB.Preload("Projects.Works").Preload("Template").Where("uid = ?", uid).First(&portfolio).Error; err != nil {
		return Portfolio{}, err
	}
	return portfolio, nil
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

func (r PortfolioRepo) SavePortfolioToDB(ctx context.Context, portfolio Portfolio) error {
	projects := portfolio.Projects
	works := []Work{}
	for _, project := range projects {
		for _, work := range project.Works {
			work.ProjectUID = project.UID
			works = append(works, work)
		}
	}
	texts := []Text{}
	for _, project := range projects {
		for _, text := range project.Texts {
			text.ProjectUID = project.UID
			texts = append(texts, text)
		}
	}
	err := r.mysqlDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uid"}},
			UpdateAll: true,
		}).Create(&portfolio).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uid"}},
			UpdateAll: true,
		}).Create(&projects).Error; err != nil {
			return err
		}

		if len(works) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "oss_key"}},
				UpdateAll: true,
			}).Create(&works).Error; err != nil {
				return err
			}
		}
		if len(texts) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "uid"}},
				UpdateAll: true,
			}).Create(&texts).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
