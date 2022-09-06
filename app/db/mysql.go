package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goframe/app/logger"
	"goframe/app/utils"
	"goframe/config"
	"os"
	"strings"
)

var mysqlConfig *config.MysqlConfig
var mysqlReadClient *sql.DB
var mysqlWriteClient *sql.DB

type Mysql interface {
	DbStar()
	Query(strSql string) map[int]map[string]string
	QueryRow(strSql string) map[string]string
	Insert(table string, data map[string]string) int64
	InsertId(table string, data map[string]string) int64
	InsertAll(table string, data map[int]map[string]string) int64
	Update(table string, where map[string]map[int]string, data map[string]string) int64
	Delete(table string, where map[string]map[int]string) int64
	ExecSql(strSql string) int64
	GetCount(table string, where map[string]map[int]string) int64
	BeginTransaction()
	Commit()
	Rollback()
}

type MysqlClass struct {
}

/**
 * [初始化]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func init() {
	mysqlConfig = config.LoadMysqlConfig()

	//连接读库
	readDns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", mysqlConfig.DB_READ_USERNAME, mysqlConfig.DB_READ_PASSWORD,
		mysqlConfig.DB_READ_HOST, mysqlConfig.DB_READ_PORT, mysqlConfig.DB_DATABASE, mysqlConfig.DB_CHARSET)

	mysqlReadClient, _ = sql.Open("mysql", readDns)

	mysqlReadClient.SetMaxOpenConns(mysqlConfig.DB_POOL_MAX)
	mysqlReadClient.SetMaxIdleConns(mysqlConfig.DB_POOL_MIN)

	err := mysqlReadClient.Ping()
	if err != nil {
		logger.Logger.Fatalln("Failed to connect read-mysql , err :" + err.Error())
		os.Exit(1)
	} else {
		logger.Logger.Println("SUCCESS: read-mysql连接成功")
	}

	//连接写库
	writeDns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", mysqlConfig.DB_WRITE_USERNAME, mysqlConfig.DB_WRITE_PASSWORD,
		mysqlConfig.DB_WRITE_HOST, mysqlConfig.DB_WRITE_PORT, mysqlConfig.DB_DATABASE, mysqlConfig.DB_CHARSET)

	mysqlWriteClient, _ = sql.Open("mysql", writeDns)

	mysqlWriteClient.SetMaxOpenConns(mysqlConfig.DB_POOL_MAX)
	mysqlWriteClient.SetMaxIdleConns(mysqlConfig.DB_POOL_MIN)

	err = mysqlWriteClient.Ping()
	if err != nil {
		logger.Logger.Fatalln("Failed to connect write-mysql , err :" + err.Error())
		os.Exit(1)
	} else {
		logger.Logger.Println("SUCCESS: write-mysql连接成功 ")
	}
}

func (mysqlClass MysqlClass) DbStar() {

}

/**
 * [查询集合]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Query(strSql string) map[int]map[string]string {
	strSql = strings.Trim(strSql, " ")

	rows, err := mysqlReadClient.Query(strSql)

	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return nil
	}

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rowList := make(map[int]map[string]string)

	var num = 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// 这个map用来存储一行数据，列名为map的key，map的value为列的值
		rowMap := make(map[string]string)

		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
		}

		rowList[num] = rowMap
		num++
	}

	rows.Close()
	logger.Logger.Println(fmt.Sprintf("sql: %s", strSql))
	return rowList
}

/**
 * [查询行]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) QueryRow(strSql string) map[string]string {
	strSql = strings.Trim(strSql, " ")

	rows, err := mysqlReadClient.Query(strSql)

	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error:%s", err.Error()))
		return nil
	}

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error:%s", err.Error()))
		return nil
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 这个map用来存储一行数据，列名为map的key，map的value为列的值
	rowMap := make(map[string]string)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
		}
		break
	}

	rows.Close()
	logger.Logger.Println(fmt.Sprintf("sql: %s", strSql))
	return rowMap
}

/**
 * [单条新增 返回影响行数]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Insert(table string, data map[string]string) int64 {

	columns := utils.GetMpaKeys(data)

	values := utils.GetMpaValues(data)

	strColumns := fmt.Sprintf("`%s`", strings.Join(columns, "`,`"))

	strValues := fmt.Sprintf("\"%s\"", strings.Join(values, "\",\""))

	strSql := fmt.Sprintf("insert into `%s` (%s) VALUES (%s)", table, strColumns, strValues)

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		logger.Logger.Println(fmt.Sprintf("succes:%d sql: %s", affected, strSql))
		return affected
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return 0
	}
}

/**
 * [单条新增 返回自增id]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) InsertId(table string, data map[string]string) int64 {

	columns := utils.GetMpaKeys(data)

	values := utils.GetMpaValues(data)

	strColumns := fmt.Sprintf("`%s`", strings.Join(columns, "`,`"))

	strValues := fmt.Sprintf("\"%s\"", strings.Join(values, "\",\""))

	strSql := fmt.Sprintf("insert into `%s` (%s) VALUES (%s)", table, strColumns, strValues)

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		logger.Logger.Println(fmt.Sprintf("succes:%d sql: %s", affected, strSql))

		id, err := exec.LastInsertId()
		if err == nil {
			return id
		}

		return 0
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s ", err.Error(), strSql))
		return 0
	}
}

/**
 * [批量新增 返回影响行数]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) InsertAll(table string, data map[int]map[string]string) int64 {
	columns := utils.GetMpaKeys(data[0])

	strColumns := fmt.Sprintf("`%s`", strings.Join(columns, "`,`"))

	strValues := make([]string, len(data))

	i := 0
	for _, v := range data {
		values := utils.GetMpaValues(v)
		strValues[i] = fmt.Sprintf("(\"%s\")", strings.Join(values, "\",\""))
		i++
	}

	strAllValues := fmt.Sprintf("%s", strings.Join(strValues, ","))

	strSql := fmt.Sprintf("insert into `%s` (%s) VALUES %s", table, strColumns, strAllValues)

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		logger.Logger.Println(fmt.Sprintf("succes:%d sql: %s", affected, strSql))
		return affected
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return 0
	}
}

/**
 * [修改 返回影响行数]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Update(table string, where map[string]map[int]string, data map[string]string) int64 {

	columnsWhere := make([]string, len(where))

	i := 0
	for k, v := range where {
		columnsWhere[i] = fmt.Sprintf("`%s` %s \"%s\"", k, v[0], v[1])
		i++
	}

	strColumnsWhere := fmt.Sprintf("%s", strings.Join(columnsWhere, " and "))

	columnsData := make([]string, len(data))

	j := 0
	for k, v := range data {
		columnsData[j] = fmt.Sprintf("`%s` = \"%s\"", k, v)
		j++
	}

	strColumnsData := fmt.Sprintf("%s", strings.Join(columnsData, ","))

	strSql := fmt.Sprintf("update `%s` set %s where %s", table, strColumnsData, strColumnsWhere)

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		logger.Logger.Println(fmt.Sprintf("succes:%d sql: %s", affected, strSql))
		return affected
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return 0
	}
}

/**
 * [删除 返回影响行数]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Delete(table string, where map[string]map[int]string) int64 {
	columnsWhere := make([]string, len(where))

	i := 0
	for k, v := range where {
		columnsWhere[i] = fmt.Sprintf("`%s` %s \"%s\"", k, v[0], v[1])
		i++
	}

	strColumnsWhere := fmt.Sprintf("%s", strings.Join(columnsWhere, " and "))

	strSql := fmt.Sprintf("delete from `%s` where %s", table, strColumnsWhere)

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		logger.Logger.Println(fmt.Sprintf("succes:%d sql: %s", affected, strSql))
		return affected
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return 0
	}
}

/**
 * [执行原生语句 返回影响行数]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) ExecSql(strSql string) int64 {

	strSql = strings.Trim(strSql, " ")
	tmpStrSql := strings.ToUpper(strSql)

	logger.Logger.Println(fmt.Sprintf("ExecSql: %s", strSql))

	flag := false

	switch {
	case strings.Index(tmpStrSql, "INSERT") == 0:
		flag = true
	case strings.Index(tmpStrSql, "UPDATE") == 0:
		flag = true
	case strings.Index(tmpStrSql, "DELETE") == 0:
		flag = true
	}

	if flag != true {
		logger.Logger.Println("ExecSql: 没有可执行方法")
		return 0
	}

	exec, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		affected, _ := exec.RowsAffected()
		return affected
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
		return 0
	}
}

/**
 * [执行原生语句]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) GetCount(table string, where map[string]map[int]string) int64 {
	columnsWhere := make([]string, len(where))

	i := 0
	for k, v := range where {
		columnsWhere[i] = fmt.Sprintf("`%s` %s \"%s\"", k, v[0], v[1])
		i++
	}

	strColumnsWhere := fmt.Sprintf("%s", strings.Join(columnsWhere, " and "))

	strSql := fmt.Sprintf("select count(1) as count from `%s` where %s", table, strColumnsWhere)

	logger.Logger.Println(fmt.Sprintf("ExecSql: %s", strSql))

	rows, err := mysqlReadClient.Query(strSql)

	if err != nil {
		logger.Logger.Println(fmt.Sprintf("error:%s", err.Error()))
		return 0
	}

	var count int64
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			logger.Logger.Println(fmt.Sprintf("error:%s", err.Error()))
		}
		break
	}

	rows.Close()

	return count
}

/**
 * [开启事务]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) BeginTransaction() {
	strSql := "START TRANSACTION;"

	_, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		logger.Logger.Println(fmt.Sprintf("开启事务 sql: %s", strSql))
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
	}
}

/**
 * [提交事务]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Commit() {
	strSql := "COMMIT;"

	_, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		logger.Logger.Println(fmt.Sprintf("提交事务 sql: %s", strSql))
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
	}
}

/**
 * [事务回滚]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func (mysqlClass MysqlClass) Rollback() {
	strSql := "ROLLBACK;"

	_, err := mysqlWriteClient.Exec(strSql)

	if err == nil {
		logger.Logger.Println(fmt.Sprintf("事务回滚 sql: %s", strSql))
	} else {
		logger.Logger.Println(fmt.Sprintf("error:%s sql: %s", err.Error(), strSql))
	}
}
