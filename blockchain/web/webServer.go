/**
  author: kevin
 */
package web

import (
	"fab-sdk-go-sample/web/controller"
	"fmt"
	"net/http"
)

func WebStart(app *controller.Application)  {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定路由信息(匹配请求)
	http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/loginout", app.LoginOut)

	http.HandleFunc("/index", app.Index)
	http.HandleFunc("/help", app.Help)

	http.HandleFunc("/qrcodesearch", app.QueryPage)


	http.HandleFunc("/addTeaInfo", app.AddTeaShow)	// 显示添加信息页面
	http.HandleFunc("/addTea", app.SaveTea)	// 提交信息请求

	//http.HandleFunc("/modifyTea", app.ModifyTea)	// Modify Tea Page
	http.HandleFunc("/modifyQuery", app.ModifyQuery)	// Modify Tea Page
	http.HandleFunc("/modifyQueryPage", app.ModifyQueryPage)	// Modify Tea Page
	http.HandleFunc("/modifyResult", app.ModifyTea)	// Modify Tea Page
	//http.HandleFunc("/modifyTea", app.ModifyTea)	// submit modify request

	http.HandleFunc("/queryPage", app.QueryPage)	// 转至根据 id 查询信息页面
	http.HandleFunc("/query", app.FindTeaByID)	// 根据 id 查询信息

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}