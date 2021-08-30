//Created by Goland
//@User: lenora
//@Date: 2021/2/5
//@Time: 2:13 下午
package server

import (
	"context"
	"fmt"
	"github.com/Lenora-Z/low-code/conf"
	"github.com/Lenora-Z/low-code/docs"
	"github.com/Lenora-Z/low-code/service/tritium"
	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strings"
	"time"
)

type Server interface {
	Run(configPath string) error
	Close() error
}

func NewCronServer(name string, task string) Server {
	s := new(cronServer)
	s.name = name
	s.task = task
	return s
}

type defaultServer struct {
	name        string
	version     string
	port        string
	conf        *conf.ServerConfig
	db          *gorm.DB
	businessDb  *gorm.DB
	dataDb      *gorose.Connection
	mongoDb     *mongo.Database
	engine      *gin.Engine
	minioClient *minio.Client
	//rpcService   micro.Service
	//rpcSmsClient smsProto.SmsService
	//registry     registry.Registry
	//reporter     go2sky.Reporter
}

func NewServer(name string, port string) Server {
	s := new(defaultServer)
	s.name = name
	s.port = port
	s.version = AppVersion
	return s
}

func (ds *defaultServer) Run(configPath string) error {
	// config
	if err := ds.config(configPath); err != nil {
		return fmt.Errorf("ds.config():%s", err.Error())
	}

	// db
	if err := ds.dbClient(); err != nil {
		return fmt.Errorf("ds.dbClient():%s", err.Error())
	}

	// mongodb
	if err := ds.mongoClient(); err != nil {
		return fmt.Errorf("ds.dbClient():%s", err.Error())
	}

	// db_business
	if err := ds.businessDbClient(); err != nil {
		return fmt.Errorf("ds.businessDbClient():%s", err.Error())
	}

	// router
	if err := ds.initRouter(); err != nil {
		return fmt.Errorf("rs.router(): %s", err.Error())
	}

	// skyWalking
	//var oTrace *go2sky.Tracer
	//var err error
	//if oTrace, err = ds.InitSkyWalking(); err != nil {
	//	return fmt.Errorf("ds.InitSkyWalking(): %s", err.Error())
	//}
	//
	//ds.rpcService = micro.NewService(
	//	micro.Name("go.micro.api.grpc.lowcode-backend"),
	//	micro.Version(ds.version),
	//	micro.WrapHandler(go2skyMicro.NewHandlerWrapper(oTrace, "User-Agent")),
	//	micro.WrapHandler(ratelimit.NewHandlerWrapper(100)),
	//	micro.WrapHandler(prometheus.NewHandlerWrapper()),
	//)
	//// rpcSmsClient
	//ds.newRpcSmsClient()
	//
	//// consul
	//if err := ds.InitConsul(); err == nil && ds.registry != nil {
	//	ds.rpcService.Init(micro.Registry(ds.registry))
	//}

	// otherInit
	if err := ds.init(); err != nil {
		return fmt.Errorf("ds.init(): %s", err.Error())
	}

	//port
	port := ds.port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	//go func() {
	//	err := ds.rpcService.Run()
	//	if err != nil {
	//		logrus.Fatal("rpc run err: ", err.Error())
	//	}
	//}()

	return ds.engine.Run(port)
}

func (ds *defaultServer) Close() error {
	if ds.db != nil {
		if err := ds.db.Close(); err != nil {
			return err
		}
	}
	if ds.dataDb != nil {
		if err := ds.dataDb.Close(); err != nil {
			return err
		}
	}

	if ds.mongoDb != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := ds.mongoDb.Client().Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ds *defaultServer) config(configPath string) error {
	ds.conf = new(conf.ServerConfig)
	err := configor.Load(ds.conf, configPath)
	if err != nil {
		return err
	}
	return nil
}

func (ds *defaultServer) dbClient() error {
	db, err := getDbConnection(ds.conf.DbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("mysql connect error:%+v", err)
	}
	fmt.Println("mysql connect successfully")
	ds.db = db

	dataDb, err := getDataDbConnection(ds.conf.BusinessDbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("data db connect error:%+v", err)
	}
	fmt.Println("data db connect successfully")
	ds.dataDb = dataDb
	return nil
}

func (ds *defaultServer) businessDbClient() error {
	db, err := getDbConnection(ds.conf.BusinessDbConfig, 0, 0)
	if err != nil {
		return fmt.Errorf("mysql_business connect error:%+v", err)
	}
	fmt.Println("mysql_business connect successfully")
	ds.businessDb = db
	return nil
}

func (ds *defaultServer) mongoClient() error {
	db, err := getMongoConnection(ds.conf.MongoConfig)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}
	ds.mongoDb = db.Database(ds.conf.MongoConfig.DbName)

	return nil
}

func (ds *defaultServer) initRouter() error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	ds.engine = r
	if ds.conf.Debug {
		swaggerDoc(r)
	}
	ds.routers()
	return nil
}

func (ds *defaultServer) init() error {
	if ds.conf.TritiumConfig.Switch {
		trimSrv := tritium.NewTritiumService(ds.conf.TritiumConfig.Api)
		trimSrv.ResourceSubmit("./resource/router.json")
	}

	minioClient, err := minio.New(ds.conf.MinioConf.EndpointIP, ds.conf.MinioConf.AccessKeyID, ds.conf.MinioConf.SecretAccessKey, false)
	if err != nil {
		return fmt.Errorf("minio: %+v", err)
	}
	ds.minioClient = minioClient
	return nil

}

//func (ds *defaultServer) InitSkyWalking() (*go2sky.Tracer, error) {
////var report go2sky.Reporter
//var err error
//if ds.conf.SkyWalking != nil {
//	ds.reporter, err = reporter.NewGRPCReporter(ds.conf.SkyWalking.OapServer)
//	logrus.Info("create gRpc reporter: oap server: ", ds.conf.SkyWalking.OapServer)
//	if err != nil {
//		return nil, fmt.Errorf("create gRpc reporter error: %v \n", err)
//	}
//} else {
//	ds.reporter, err = reporter.NewLogReporter()
//	logrus.Info("create log reporter: os.stderr")
//	if err != nil {
//		return nil, fmt.Errorf("create log reporter error: %v \n", err)
//	}
//}
//oTrace, err := go2sky.NewTracer(ds.name, go2sky.WithReporter(ds.reporter))
//if err != nil {
//	return nil, fmt.Errorf("create tracer error: %v \n", err)
//}
//return oTrace, nil
//}

func (ds *defaultServer) InitConsul() error {
	//if ds.conf.Consul != nil {
	//	ds.registry = consul.NewRegistry(
	//		registry.Addrs(fmt.Sprintf("%s:%d", ds.conf.Consul.Ip, ds.conf.Consul.Port)),
	//	)
	//}
	return nil
}

//go.micro.api.grpc.backend-basic
func (ds *defaultServer) newRpcSmsClient() {
	//ds.rpcSmsClient = smsProto.NewSmsService("go.micro.api.grpc.backend-basic", client.DefaultClient)
}

func getDbConnection(conf *conf.DBConfig, maxIdle, maxOpen int) (*gorm.DB, error) {
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	dsn := fmt.Sprintf(format, conf.User, conf.Password, conf.Name, conf.Port, conf.DbName, conf.Charset, url.QueryEscape(conf.Loc))
	logrus.Infof("dsn=%s", dsn)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.LogMode(conf.LogMode)
	idle := conf.MaxIdle
	if maxIdle > 0 {
		idle = maxIdle
	}
	open := conf.MaxOpen
	if maxOpen > 0 {
		open = maxOpen
	}
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	return db, nil
}

func getDataDbConnection(conf *conf.DBConfig, maxIdle, maxOpen int) (*gorose.Connection, error) {
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	dsn := fmt.Sprintf(format, conf.User, conf.Password, conf.Name, conf.Port, conf.DbName, conf.Charset, url.QueryEscape(conf.Loc))
	idle := conf.MaxIdle
	if maxIdle > 0 {
		idle = maxIdle
	}
	open := conf.MaxOpen
	if maxOpen > 0 {
		open = maxOpen
	}
	conn, err := gorose.Open(&gorose.DbConfigSingle{
		Driver:          "mysql",
		EnableQueryLog:  conf.LogMode,
		SetMaxOpenConns: open,
		SetMaxIdleConns: idle,
		Prefix:          "",
		Dsn:             dsn,
	})
	if err != nil {
		return nil, err
	}
	conn.Use(gorose.NewLogger())
	return conn, nil
}

// https://mongodb-documentation.readthedocs.io/en/latest/reference/connection-string.html#gsc.tab=0
func getMongoConnection(conf *conf.MongoConfig) (*mongo.Client, error) {
	format := "mongodb://%s:%s@%s:%s/%s?maxPoolSize=%s&minPoolSize=%s"
	uri := fmt.Sprintf(format, conf.User, conf.Password, conf.IP, conf.Port, conf.DbName, conf.MaxOpen, conf.MaxIdle)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func swaggerDoc(ctx *gin.Engine) {
	docs.SwaggerInfo.Title = AppTitle
	docs.SwaggerInfo.Description = AppName
	docs.SwaggerInfo.Version = AppVersion
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	ctx.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
