package main

import (
	"encoding/gob"
	"html/template"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/short-net/dbrepo"
	"github.com/lightsaid/short-net/mailer"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/token"
)

type application struct {
	env           envConfig
	shortID       uint
	store         dbrepo.Repository
	templateCache map[string]*template.Template
	sessionMgr    *scs.SessionManager
	mailer        mailer.Mailer
	tokenMaker    token.Maker
	wg            sync.WaitGroup
	mutex         sync.RWMutex
	redisPool     *redis.Pool
}

func main() {
	// 加载配置
	envConf := setupConfig()

	// 设置日志
	setupLogger()

	// 连接、迁移 db
	db := setupDB()

	conn, err := db.DB()
	fatalOnError(err, "Get *sql.DB failed")

	// 设置最大连接数、空闲数，重用时间
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)
	defer conn.Close()

	// 执行迁移
	err = db.AutoMigrate(&models.User{}, &models.Link{}, &models.Book{}, &models.Order{}, &models.OrderDetail{})
	fatalOnError(err, "db.AutoMigrate failed")

	// redis 连接
	redisPool := setupRedis(&envConf)

	// session 创建
	sessionMgr := setupSessionMgr(&envConf)
	// 存储 session
	sessionMgr.Store = redisstore.New(redisPool)

	// test redis
	// pool := redisPool.Get()
	// defer pool.Close()
	// res, err := pool.Do("SET", "short_version", "v1.0")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf(">>>> redis: %v\n", res)

	// 获取link表最后一条记录的short_hash值对应数值，用于生成后续短网址
	shortID := setupShortID(db)

	// tokenMaker
	tokenMaker, err := token.NewTokenMaker(envConf.TokenSecretKey)
	fatalOnError(err, "token maker create failed")

	var app = application{
		env:        envConf,
		shortID:    shortID,
		sessionMgr: sessionMgr,
		tokenMaker: tokenMaker,
		wg:         sync.WaitGroup{},
		mutex:      sync.RWMutex{},
		redisPool:  redisPool,
	}

	app.mailer = mailer.NewMailSender(
		envConf.SmtpAuthAddress,
		envConf.SmtpServerAddress,
		envConf.MailSenderName,
		envConf.MailSenderAddress,
		envConf.MailSenderPassword,
	)

	app.store = dbrepo.NewRepository(db)

	// 注册 gob 数据，用于序列化 cookie 存储的数据
	gob.Register(renderData{})
	gob.Register(models.User{})

	// 加载模板
	err = app.genTemplateCache()
	fatalOnError(err, "genTemplateCache failed")

	err = app.serve()
	fatalOnError(err, "app.serve failed")
}
