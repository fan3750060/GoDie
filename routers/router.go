package routers

import (
	"myproject/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.AutoRouter(&controllers.MainController{})
    beego.Router("/", &controllers.MainController{})
    beego.Router("/user", &controllers.MainController{})	
}
