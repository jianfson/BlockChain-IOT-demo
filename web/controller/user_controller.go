package controller

import (
	"blc-iot-demo/web/service"
	"blc-iot-demo/web/dao"
	"blc-iot-demo/web/model"
	"blc-iot-demo/web/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//进入首页
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "index.html", data)
}

// 返回首页
func (app *Application) BackToHome(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "index.html", data)
}

// 进入注册界面
func (app *Application) RegisterPage(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "AccountRelated/register.html", nil)
}

// 注册添加用户信息
func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("loginName")
	Password := r.FormValue("password")
	email := r.FormValue("email")
	phone := r.FormValue("tel")
	role := "用户"
	statue := "正常"

	if username == "" || Password == "" || email == "" || phone == "" {
		ShowView(w, r, "AccountRelated/register.html", nil)
		return
	} else {

		ID := dao.QueryUserWithUsername(username)
		fmt.Println("id", ID)
		if ID > 0 {
			ShowView(w, r, "AccountRelated/register.html", nil)
			return
		}
		password := utils.MD5(Password)
		createtime := utils.SwitchTimeStampToData(time.Now().Unix())
		fmt.Println(createtime)

		user := model.User{Username: username, Password: password, Role: role, Email: email, Phone: phone, Status: statue, Createtime: createtime}
		_, err := dao.InsertUser(user) //user表插入记录
		//为该用户建表

		if err != nil {
			ShowView(w, r, "AccountRelated/register.html", nil)
			return
		} else {
			ShowView(w, r, "AccountRelated/login.html", nil)
			return
		}

	}
}

// 进入登录界面
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "index.html", data)
		return
	}
	ShowView(w, r, "AccountRelated/login.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Sess         *model.Session
		FailedLogin  bool
		IsLogin      bool
		IsSuperAdmin bool
		IsAdmin      bool
		IsUser       bool
		IsStaff      bool
		Msg          string
	}{
		Sess:         nil,
		FailedLogin:  false,
		IsLogin:      false,
		IsSuperAdmin: false,
		IsAdmin:      false,
		IsUser:       false,
		IsStaff:      false,
		Msg:          "",
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("默认参数已就绪")

	//获取表格信息
	username := r.FormValue("loginName")
	Password := r.FormValue("password")
	password := utils.MD5(Password)
	fmt.Println("---------------------------------------------")
	fmt.Println("前端表格读取完成")

	//返回完整的用户信息
	user := dao.FindUserByUsernameAndPassword(username, password)
	fmt.Println("---------------------------------------------")
	fmt.Println("用户", user.Username, "查询结果已传回，正在核查")

	if user.Id == 0 {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println("用户名或密码错误，登陆失败，以未登录状态返回首页")
		data.Msg = "用户名或密码错误"
		ShowView(w, r, "AccountRelated/login.html", data)
		return

	} else if user.Status == "异常" {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println(user.Role, user.Username, "账户受限，登陆失败，以未登录状态返回首页")
		data.Msg = user.Role + user.Username + "账户受限，登陆失败，请联系管理员"
		ShowView(w, r, "index.html", data)
		return

	} else if user.Status == "正常" {
		uuid := utils.CreateUUID()
		session := &model.Session{
			SessionID:  uuid,
			UserID:     user.Id,
			UserName:   user.Username,
			PassWord:   user.Password,
			Role:       user.Role,
			Email:      user.Email,
			Phone:      user.Phone,
			Status:     user.Status,
			CreateTime: user.Createtime,
		}

		dao.AddSession(session)

		fmt.Println("---------------------------------------------")
		fmt.Println("Session已设置")

		cookie := http.Cookie{
			Name:     "user",
			Value:    uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		fmt.Println("---------------------------------------------")
		fmt.Println("Cookie已送往浏览器")

		data.IsLogin = true
		data.Sess = session
		fmt.Println("---------------------------------------------")
		fmt.Println("默认参数已更新")

		if user.Role == "超级管理员" {
			data.IsSuperAdmin = true
			ShowView(w, r, "SuperBackStage/superBackStage.html", data)
			fmt.Println("---------------------------------------------")
			fmt.Println("超级管理员登陆成功！页面已跳转")
			return

		} else if user.Role == "管理员" {
			data.IsAdmin = true
			ShowView(w, r, "BackStage/backStage.html", data)
			fmt.Println("---------------------------------------------")
			fmt.Println("管理员登陆成功！页面已跳转")
			return

		} else if user.Role == "用户" {
			data.IsUser = true
			ShowView(w, r, "index.html", data)
			fmt.Println("---------------------------------------------")
			fmt.Println("普通用户：", user.Username, "登陆成功！页面已跳转")
			return

		} else if user.Role == "员工" {
			data.IsStaff = true
			ShowView(w, r, "index.html", data)
			fmt.Println("---------------------------------------------")
			fmt.Println("员工：", user.Username, "登陆成功！页面已跳转")
			return

		} else {
			data.FailedLogin = true
			data.IsLogin = false
			data.Sess = nil
			fmt.Println("---------------------------------------------")
			fmt.Println("登陆失败，以未登录状态返回首页")
			data.Msg = "登陆失败，原因不明，请重试，或联系管理员"
			ShowView(w, r, "index.html", data)
			return
		}
	}
}

// 退出登陆
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	data := utils.DeleteSession(r)
	ShowView(w, r, "index.html", data)
}

// 进入个人界面
func (app *Application) ProfilePage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "AccountRelated/profilePage.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

//-----------------------------------------------------------------------------------------
// 进入超管后台
func (app *Application) SuperBackStage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "SuperBackStage/superBackStage.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入数据管理
func (app *Application) SbsDataMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		result, _ := app.Setup.QueryTeaByMaker("杭州龙井茶业集团有限公司")
		var teas []service.Tea
		json.Unmarshal(result, &teas)
		var Teas []*service.Tea
		for i := 0; i < len(teas); i++ {
			tea := teas[i]
			Teas = append(Teas, &tea)
		}
		data.Teas = Teas
		ShowView(w, r, "SuperBackStage/sbsDataMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入管理员任命
func (app *Application) SbsAdminMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	fmt.Println("---------------------------------------------")
	fmt.Println("查询已有管理员")
	admins, _ := dao.QueryAllAdmin()
	data.Admin = admins

	fmt.Println("---------------------------------------------")
	fmt.Println("分流")
	if data.IsLogin {
		ShowView(w, r, "SuperBackStage/sbsAdminMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

//添加新管理员
func (app *Application) SbsAddNewAdmin(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		username := r.FormValue("loginName")
		Password := r.FormValue("password")
		email := r.FormValue("email")
		phone := r.FormValue("tel")
		role := "管理员"
		statue := "正常"

		ID := dao.QueryUserWithUsername(username)
		fmt.Println("id", ID)

		fmt.Println("---------------------------------------------")
		fmt.Println("分流")

		if ID > 0 {
			data.Msg = "管理员已存在，请重试"
			ShowView(w, r, "SuperBackStage/sbsAdminMana.html", data)
			return
		}
		password := utils.MD5(Password)
		createtime := utils.SwitchTimeStampToData(time.Now().Unix())
		user := model.User{Username: username, Password: password, Role: role, Email: email, Phone: phone, Status: statue, Createtime: createtime}
		_, err := dao.InsertUser(user) //user表插入记录
		//为该用户建表

		if err != nil {
			data.Msg = "出现错误！"
			ShowView(w, r, "SuperBackStage/sbsAdminMana.html", data)
			return
		} else {
			data.Msg = "任命成功！"
			fmt.Println("---------------------------------------------")
			fmt.Println("查询已有用户")
			admins, _ := dao.QueryAllAdmin()
			data.Admin = admins
			ShowView(w, r, "SuperBackStage/sbsAdminMana.html", data)
			return
		}
	} else {
		data.Msg = "请登录"
		fmt.Println("请登录")
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入用户管理
func (app *Application) SbsUserMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		fmt.Println("---------------------------------------------")
		fmt.Println("查询已有用户")
		users, _ := dao.QueryAllUser()
		data.User = users

		fmt.Println("---------------------------------------------")
		fmt.Println("分流")
		ShowView(w, r, "SuperBackStage/sbsUserMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入操作指南
func (app *Application) SbsInstruction(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		ShowView(w, r, "SuperBackStage/sbsInstruction.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

//-----------------------------------------------------------------------------------------
// 进入后台
func (app *Application) BackStage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "BackStage/backStage.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入用户管理
func (app *Application) BsUserMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	fmt.Println("---------------------------------------------")
	fmt.Println("查询已有用户")
	users, _ := dao.QueryAllUser()
	data.User = users

	fmt.Println("---------------------------------------------")
	fmt.Println("分流")

	if data.IsLogin {
		ShowView(w, r, "BackStage/bsUserMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入数据管理
func (app *Application) BsDataMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		result, _ := app.Setup.QueryTeaByMaker("dragonwell")
		var teas []service.Tea
		json.Unmarshal(result, &teas)
		var Teas []*service.Tea
		for i := 0; i < len(teas); i++ {
			tea := teas[i]
			Teas = append(Teas, &tea)
		}
		data.Teas = Teas
		ShowView(w, r, "BackStage/bsDataMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 进入职员管理
func (app *Application) BsStaffMana(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		fmt.Println("---------------------------------------------")
		fmt.Println("查询已有职员")
		staffs, _ := dao.QueryAllStaff()
		data.Staff = staffs
		ShowView(w, r, "BackStage/bsStaffMana.html", data)
		return
	} else {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

//添加新职员
func (app *Application) BsAddNewStaff(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsLogin {
		username := r.FormValue("loginName")
		Password := r.FormValue("password")
		email := r.FormValue("email")
		phone := r.FormValue("tel")
		role := "员工"
		statue := "正常"

		ID := dao.QueryUserWithUsername(username)
		fmt.Println("id", ID)

		fmt.Println("---------------------------------------------")
		fmt.Println("分流")

		if ID > 0 {
			data.Msg = "职员已存在，请重试"
			ShowView(w, r, "BackStage/bsStaffMana.html", data)
			return
		}
		password := utils.MD5(Password)
		createtime := utils.SwitchTimeStampToData(time.Now().Unix())
		user := model.User{Username: username, Password: password, Role: role, Email: email, Phone: phone, Status: statue, Createtime: createtime}
		_, err := dao.InsertUser(user) //user表插入记录
		//为该用户建表

		if err != nil {
			data.Msg = "出现错误！"
			ShowView(w, r, "BackStage/bsStaffMana.html", data)
			return
		} else {
			data.Msg = "任命成功！"
			fmt.Println("---------------------------------------------")
			fmt.Println("查询已有职员")
			staffs, _ := dao.QueryAllAdmin()
			data.Staff = staffs
			ShowView(w, r, "BackStage/bsStaffMana.html", data)
			return
		}
	} else {
		data.Msg = "请登录"
		fmt.Println("请登录")
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}
