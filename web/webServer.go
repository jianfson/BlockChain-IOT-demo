/**
  author: kevin
 */
package web

import (
	"blc-iot-demo/web/controller"
	"fmt"
	"net/http"
)

func WebStart(app *controller.Application)  {


	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 指定第一次打开系统进入的页面
	http.HandleFunc("/", app.Home)

	http.HandleFunc("/backToHome", app.BackToHome)			// 返回首页

	// 登陆
	http.HandleFunc("/loginPage", app.LoginView)
	http.HandleFunc("/login", app.Login)

	http.HandleFunc("/logout", app.Logout)

	http.HandleFunc("/registerPage", app.RegisterPage)
	http.HandleFunc("/register", app.Register)

	//用户操作页面路由
	http.HandleFunc("/profilePage", app.ProfilePage)

	//超管后台
	http.HandleFunc("/superBackStage", app.SuperBackStage)
	http.HandleFunc("/sbsInstruction", app.SbsInstruction)
	http.HandleFunc("/sbsDataMana", app.SbsDataMana)
	http.HandleFunc("/sbsAdminMana", app.SbsAdminMana)
	http.HandleFunc("/sbsUserMana", app.SbsUserMana)
	http.HandleFunc("/sbsAddNewAdmin", app.SbsAddNewAdmin)

	//后台
	http.HandleFunc("/backStage", app.BackStage)
	http.HandleFunc("/bsDataMana", app.BsDataMana)
	http.HandleFunc("/bsUserMana", app.BsUserMana)
	http.HandleFunc("/bsStaffMana", app.BsStaffMana)
	http.HandleFunc("/bsAddNewStaff", app.BsAddNewStaff)
	http.HandleFunc("/modifyUserStatus", app.ModifyStatus)

	// 添加
	http.HandleFunc("/addTeaPage", app.AddTeaPage) // 显示添加信息页面
	http.HandleFunc("/addTea", app.AddTea)         // 提交修改请求并跳转添加成功提示页面
	http.HandleFunc("/bulkAddTea", app.BulkAddTea)

	// 修改
	//http.HandleFunc("/modifyQueryPage", app.ModifyQueryPage) // 进入修改查询页面
	http.HandleFunc("/modifyQuery", app.ModifyQuery)         // 显示查询结果并修改
	http.HandleFunc("/modifyResult", app.ModifyResult)       // 显示修改结果

	// 查询
	http.HandleFunc("/queryPage", app.QueryPage)		// 转至查询信息页面
	http.HandleFunc("/findTeaByID", app.FindTeaByID)	// 根据id查询并转至查询结果页面

	//http.HandleFunc("/backStagePage", app.BackStagePage)


	fmt.Println("---------------------------------------------")
	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
