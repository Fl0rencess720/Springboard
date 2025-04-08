package data

import (
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func Init() {
	mysqlInit()
	redisInit()
}

func mysqlInit() {
	mysqlDB, err := gorm.Open(mysql.Open(viper.GetString("data.mysql.source")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect mysql")
	}
	db = mysqlDB
}

func redisInit() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("data.redis.addr"),
		Password:     viper.GetString("data.redis.password"),
		DB:           viper.GetInt("data.redis.db"),
		DialTimeout:  viper.GetDuration("data.redis.dial_timeout"),
		WriteTimeout: viper.GetDuration("data.redis.write_timeout"),
		ReadTimeout:  viper.GetDuration("data.redis.read_timeout"),
	})
	redisClient.AddHook(redisotel.TracingHook{})
	rdb = redisClient
}

func GetDB() *gorm.DB {
	return db
}

func GetRedis() *redis.Client {
	return rdb
}

func Close() {
	dbSQL, err := db.DB()
	if err != nil {
		panic("failed to close mysql")
	}
	if err := dbSQL.Close(); err != nil {
		panic("failed to close mysql")
	}
	if err := rdb.Close(); err != nil {
		panic("failed to close redis")
	}
}
