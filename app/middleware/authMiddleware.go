package middleware

import (
	"goframe/app/logger"
	"goframe/app/utils"
	"net/http"
)

func Process(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		//获取ip 可做白名单限制TODO
		ip := utils.GetClientIp(r)

		logger.Logger.Printf("%s %s %s %s %d", r.Method, ip, r.URL.String(), r.Proto, r.ContentLength)

		//设置跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Printf("Error: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
