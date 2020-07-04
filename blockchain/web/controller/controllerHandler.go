package controller

import (
	"encoding/json"
	"fab-sdk-go-sample/service"
	"fmt"
	"net/http"
)

//
var cuser User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "login.html", nil)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "index.html", nil)
}

func (app *Application) Help(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
	}{
		CurrentUser: cuser,
	}
	ShowView(w, r, "help.html", data)
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

// 用户登出
func (app *Application) LoginOut(w http.ResponseWriter, r *http.Request) {
	cuser = User{}
	ShowView(w, r, "login.html", nil)
}

// 显示添加信息页面
func (app *Application) AddTeaShow(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "addTea.html", data)
}

// 显示修改信息页面
func (app *Application) ModifyTea(w http.ResponseWriter, r *http.Request) {

		teaId := r.FormValue("new_id")

		nextOwner := r.FormValue("new_owner")
fmt.Println(teaId, nextOwner)


/*	teaID := r.FormValue("id")
	result, err := app.Setup.FindTeaInfoByID(teaID)
	var teaInfo = service.Tea{}
	json.Unmarshal(result, &teaInfo)

	data := &struct {
		Tea         service.Tea
		CurrentUser User
		Msg         string
		Flag        bool
		History     bool
	}{
		Tea:         teaInfo,
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
		History:     true,
	}

	if err != nil {
		data.Msg = err.Error()
		data.Flag = true
	}*/

	app.Setup.ModifyTea(teaId, nextOwner)


	ShowView(w, r, "modifySuccess.html", teaId)
}

// 添加 tea
func (app *Application) SaveTea(w http.ResponseWriter, r *http.Request) {

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

func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		CurrentUser User
		Msg         string
		Flag        bool
	}{
		CurrentUser: cuser,
		Msg:         "",
		Flag:        false,
	}
	ShowView(w, r, "query.html", data)
}


// 进入修改查询页面
func (app *Application) ModifyQueryPage(w http.ResponseWriter, r *http.Request) {

	ShowView(w, r, "modifyQuery.html", nil)
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
	ShowView(w, r, "modify.html", data)
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



// 扫描二维码查询
/*func (app *Application) FindTeaByQrcode(w http.ResponseWriter, r *http.Request) {
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
}*/
