package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"github.com/natefinch/lumberjack"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var initShortID uint = 10000

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("msg: %s", err.Error())
	}
}

func setupConfig() envConfig {
	var envConf envConfig
	env, err := util.Loadingenv(".env")
	fatalOnError(err, "Loadingenv failed")

	err = util.Setingenv(&envConf, env)
	fatalOnError(err, "Setingenv failed")

	return envConf
}

func setupLogger() {
	if os.Getenv("RUN_MODE") == "prod" {
		var logfile = os.Getenv("ACCESS_LOG")
		if logfile == "" {
			logfile = "./storage/log/access.log"
		}
		// 日志分割
		slog.SetDefault(slog.New(slog.NewJSONHandler(&lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})))
	} else {
		// 标准输出
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout)))
	}
}

// setupShortID 从数据库里link表获取最后一记录的hash，转为uint id，系统后续生成hash，基于此id递增
func setupShortID(db *gorm.DB) uint {
	var link models.Link
	err := db.Last(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return initShortID
	}
	fatalOnError(err, "init shortID server error")

	return util.DecodeBase62(link.ShortHash)
}

func setupDB() *gorm.DB {
	var times = 5
	var db *gorm.DB
	var err error
	mylog, colorful := setGormLog()
	for i := 1; i <= times; i++ {
		db, err = connectDB(mylog, colorful)
		if err != nil {
			fmt.Printf("第 %d 次尝试连接数据库失败\n", i)
			time.Sleep(2 * time.Second)
		} else {
			return db
		}
	}

	fatalOnError(err, "connect mysql failed")
	return db
}

func setGormLog() (*log.Logger, bool) {
	var colorful = true
	var logfile = os.Getenv("MYSQL_LOG")
	if logfile == "" {
		logfile = "./storage/log/db.log"
	}
	mode := os.Getenv("RUN_MODE")
	mylog := log.New(os.Stdout, "\r\n", log.LstdFlags)
	if mode == "prod" {
		colorful = true
		mylog = log.New(&lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    1, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}, "\r\n", log.LstdFlags)
	}

	return mylog, colorful
}

func connectDB(mylog *log.Logger, colorful bool) (*gorm.DB, error) {
	logConfig := logger.New(
		// io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		// log.New(os.Stdout, "\r\n", log.LstdFlags),
		mylog,
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  colorful,    // 彩色打印
		},
	)

	return gorm.Open(mysql.Open(os.Getenv("DB_SOURCE")), &gorm.Config{Logger: logConfig})
}

// setupCookieMgr 设置cookie管理
func setupSessionMgr(env *envConfig) *scs.SessionManager {
	var secure bool
	if env.RunMode == "prod" {
		secure = true
	}
	sessionMgr := scs.New()
	sessionMgr.Lifetime = env.SessionLifeTime
	sessionMgr.Cookie.Persist = true
	sessionMgr.Cookie.SameSite = http.SameSiteLaxMode
	sessionMgr.Cookie.Secure = secure

	return sessionMgr
}

func setupRedis(env *envConfig, options ...redis.DialOption) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 3 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", env.RedisAddress, options...)
		},
	}

	return pool
}
