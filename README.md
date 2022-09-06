2021-09-22

1.env配置文件解析

2.新增mysql驱动,支持主从,连接池,读写分离

3.新增redis驱动,支持连接池

4.新增http服务器

2021-09-23

5.新增路由驱动,支持中间件

6.新增http参数解析器,支持get,post,json等参数接受

7.新增接口输出模块,支持json返回

2021-09-24

8.新增mysql,redis接口

9.新增日志系统

mysql操作
    
    //引入db-mysql库
    import (
        "fmt"
        "goframe/app/db"
        "log"
    )

    var mysql db.Mysql = new(db.MysqlClass)
    
    //查询多条
    res := mysql.Query("SELECT * FROM Magic limit 10")
    for i := 0; i <= len(res); i++ {
        value := res[i]
        fmt.Printf("1-key:%s MagID:%s MagName:%s\n", i, value["MagID"], value["MagName"])
    }
    
    //查询单条
    info := mysql.QueryRow("SELECT * FROM Magic limit 10")
    fmt.Printf("2 MagID:%s MagName:%s\n", info["MagID"], info["MagName"])

    //插入数据
    data := make(map[string]string)
    data["account"] = "test_001"
    data["username"] = "test_001"
    infoRows := mysql.Insert("account", data)
    fmt.Printf("插入成功: %d\n", infoRows)
    
    //插入数据-自增id
    data2 := make(map[string]string)
    data2["account"] = "test_001"
    data2["username"] = "test_001"
    infoId := mysql.InsertId("account", data2)
    fmt.Printf("插入成功 id: %d\n", infoId)
    
    //批量插入数据-返回受影响行数
    data3 := make(map[int]map[string]string)
    data3[0] = data
    data3[1] = data2
    infoRows2 := mysql.InsertAll("account", data3)
    fmt.Printf("插入成功: %d\n", infoRows2)
    
    //修改数据
    where := make(map[string]map[int]string)
    where["account"] = map[int]string{
        0: "=",
        1: "test_001",
    }
    
    updateData := make(map[string]string)
    updateData["password"] = "测试密码"
    updateData["birth_day"] = "2000/10/10"
    updateData["questions"] = "这是问题"
    
    updateRes := mysql.Update("account", where, updateData)
    fmt.Printf("更新成功: %d\n", updateRes)
    
    //删除数据
    whereDelete := make(map[string]map[int]string)
    whereDelete["account"] = map[int]string{
        0: "=",
        1: "test_001",
    }

    deleteRes := mysql.Delete("account", whereDelete)
    fmt.Printf("删除成功: %d\n", deleteRes)
    
    
    //原生sql执行
    execSql := "insert into account (`account`,`password`) values ('test_001','test_001') "
    execRes := mysql.ExecSql(execSql)
    fmt.Printf("执行成功: %d\n", execRes)

    execSql = "update account set questions='这是什么问题' where account = 'test_001' "
    execRes = mysql.ExecSql(execSql)
    fmt.Printf("执行成功: %d\n", execRes)
    
    //获取长度
    whereCount := make(map[string]map[int]string)
    whereCount["account"] = map[int]string{
        0: "=",
        1: "test_001",
    }
    
    countRes := mysql.GetCount("account", whereCount)
    fmt.Printf("长度为: %d\n", countRes)
    
    //事务处理
    mysql.BeginTransaction()
    dataTransaction := make(map[string]string)
    dataTransaction["account"] = "test_001"
    dataTransaction["username"] = "test_001"
    infoRows := mysql.Insert("account", dataTransaction)
    fmt.Printf("插入成功: %d\n", infoRows)
    //mysql.Commit()
    mysql.Rollback()
    
    
redis操作

    var redis db.Redis = new(db.RedisClass)

    //设置缓存,第三个参数为时间(秒),不设置为永久
    setRes := redis.Set("test_01", "value_022221", 60)
    fmt.Printf("redis 设置: %s\n", strconv.FormatBool(setRes))

    //获取缓存
    getRes := redis.Get("test_01")
    fmt.Printf("redis 获取: %s\n", getRes)

    //删除缓存
    delRes := redis.Del("test_01")
    fmt.Printf("redis 删除: %s\n", strconv.FormatBool(delRes))

    //获取redis连接
    redisConn := redis.GetConn()
    println(redisConn)
    
    //...其他操作自行实现 用到再说