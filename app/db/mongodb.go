package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goframe/app/logger"
	"goframe/config"
)

var mongoDbConfig *config.MongoDbConfig
var MongoDbClient *mongo.Client

/**
 * [初始化]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func init() {
	return
	mongoDbConfig = config.LoadMongoDbConfig()

	// 设置mongoDB客户端连接信息
	mongoDbDns := fmt.Sprintf("mongodb://%s:%s@%s:%s,%s:%s/%s", mongoDbConfig.MONGODB_USERNAME, mongoDbConfig.MONGODB_PASSWORD,
		mongoDbConfig.MONGODB_WRITE_HOST, mongoDbConfig.MONGODB_PORT, mongoDbConfig.MONGODB_READ_HOST, mongoDbConfig.MONGODB_PORT,
		mongoDbConfig.MONGODB_AUTHDB)

	//设置副本集
	if len(mongoDbConfig.REPLICASETNAME) > 0 {
		mongoDbDns += "?replicaSet=" + mongoDbConfig.REPLICASETNAME
	}

	clientOptions := options.Client().ApplyURI(mongoDbDns)

	// 建立客户端连接
	MongoDbClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error: %s  \n ", err.Error()))
		return
	}

	//设置读偏好TODO

	//是否连接成功
	err = MongoDbClient.Ping(context.TODO(), nil)

	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error: %s  \n ", err.Error()))
		return
	}

	logger.Logger.Println("SUCCESS: mongodb连接成功")
}

/**
 * [新增数据]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func InsertMongo(table string, data map[string]string) (insertID primitive.ObjectID) {

	collection := MongoDbClient.Database(mongoDbConfig.MONGODB_DATABASE).Collection(table)

	insertRest, err := collection.InsertOne(context.TODO(), bson.D{{"name", "Alice"}})
	if err != nil {
		fmt.Println(err)
		return
	}

	insertID = insertRest.InsertedID.(primitive.ObjectID)

	// 断开客户端连接
	err = MongoDbClient.Disconnect(context.TODO())
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error: %s  \n ", err.Error()))
	}

	return insertID
}
