package controller

import (
	"encoding/json"
	"fab-sdk-go-sample/service"
	"net/http"
)

//实例化User
var cuser User

// 返回首页
func (app *Application) BackToHome(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "index.html", nil)
}

// 进入登录界面
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "login.html", nil)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	loginName := r.FormValue("loginName")
	password := r.FormValue("password")

	var flag bool
	for _, user := range users {
		if user.LoginName == loginName && user.Password == password {
			cuser = user
			flag = true
			break
		}
	}

	data := &struct {
		CurrentUser User
		Flag        bool
	}{
		Flag: false,
	}

	if flag {
		// 登录成功
		ShowView(w, r, "index.html", data)
	} else {
		// 登录失败
		data.Flag = true
		data.CurrentUser.LoginName = "123"
		ShowView(w, r, "login.html", data)
	}
}

// 进入查询页面
func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "queryPage.html", nil)
}

// 根据teaID查询信息
func (app *Application) FindTeaByID(w http.ResponseWriter, r *http.Request) {
	teaID := r.FormValue("id")
	result, err := app.Setup.FindTeaInfoByID(teaID)
	var tea = service.Tea{}
	json.Unmarshal(result, &tea)

	data := &struct {
		Tea         service.Tea
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Tea:         tea,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	ShowView(w, r, "queryResult.html", data)
}

// 显示添加信息页面
func (app *Application) AddTeaPage(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "addTeaPage.html", data)
}

// 添加信息
func (app *Application) AddTea(w http.ResponseWriter, r *http.Request) {
	tea := service.Tea{
		Id:     r.FormValue("new_id"),
		Maker:  r.FormValue("new_maker"),
		Owner:  r.FormValue("new_owner"),
		Weight: r.FormValue("new_weight"),
	}
	data := &struct {
		Tea         service.Tea
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Tea:         tea,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	app.Setup.SaveTea(tea)
	ShowView(w, r, "addSuccess.html", data)
}

// 进入修改查询页面
func (app *Application) ModifyQueryPage(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "modifyQueryPage.html", nil)
}

// 根据teaID查询信息
func (app *Application) ModifyQuery(w http.ResponseWriter, r *http.Request) {
	teaID := r.FormValue("id")
	result, err := app.Setup.FindTeaInfoByID(teaID)
	var tea = service.Tea{}
	json.Unmarshal(result, &tea)

	data := &struct {
		Tea         service.Tea
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Tea:         tea,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}
	ShowView(w, r, "modifyPage.html", data)
}

// 修改信息传回上一层并显示修改成功页面
func (app *Application) ModifyResult(w http.ResponseWriter, r *http.Request) {

	//待优化
	teaId := r.FormValue("new_id")
	nextOwner := r.FormValue("new_owner")

	tea := service.Tea{
		Id:     r.FormValue("new_id"),
		Maker:  r.FormValue("new_maker"),
		Owner:  r.FormValue("new_owner"),
		Weight: r.FormValue("new_weight"),
	}

	data := &struct {
		Tea         service.Tea
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Tea:         tea,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	app.Setup.ModifyTea(teaId, nextOwner)
	ShowView(w, r, "modifySuccess.html", data)
}
