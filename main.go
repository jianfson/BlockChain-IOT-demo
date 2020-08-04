package main

import (
	"blc-iot-demo/web"
	"blc-iot-demo/web/cliInit"
	"blc-iot-demo/web/controller"
	"blc-iot-demo/web/dao"
)

func main() {

	//DB Conn
	dao.InitMysql()

	//Web
	app := controller.Application{
		cliInit.CliInit(),
	}

	defer cliInit.SDK.Close()

	web.WebStart(&app)
}
