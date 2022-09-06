package controller

import (
	"fmt"
	"goframe/app/db"
	"goframe/app/utils"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var IndexParam struct {
		MagID    string `form:"magid"`
		MagName  string `form:"magname"`
		Page     int64  `form:"page"`
		PageSize int64  `form:"pagesize"`
	}

	if err := utils.BindParam(r, &IndexParam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var mysql db.Mysql = new(db.MysqlClass)
	where := ""
	if len(IndexParam.MagID) > 0 {
		where += fmt.Sprintf(" MagID = %s", IndexParam.MagID)
	}

	where2 := ""
	if len(IndexParam.MagID) > 0 {
		where2 += fmt.Sprintf(" MagName = '%s'", IndexParam.MagName)
	}

	sql := "select * from Magic"

	if len(where) > 0 {
		sql += " where " + where
	}

	if len(where2) > 0 {
		sql += " and " + where2
	}

	page := int64(1)
	if IndexParam.Page > 0 {
		page = IndexParam.Page
	}

	pageSize := int64(20)
	if IndexParam.PageSize > 0 {
		pageSize = IndexParam.PageSize
	}

	star := (page - 1) * pageSize

	sql += fmt.Sprintf(" limit %d offset %d", pageSize, star)

	res := mysql.Query(sql)

	utils.ReturnJson(w, 200, "succes", res)
}
