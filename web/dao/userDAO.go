package dao

import (
	"blc-iot-demo/web/model"
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

//Create user table
func CreateTableWithUser() {
	sqlStr := `CREATE TABLE IF NOT EXISTS user (
				id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
				username VARCHAR (64),
				PASSWORD VARCHAR (64),
				role varchar(64),
				email VARCHAR (64),
				phone VARCHAR (64),
				STATUS varchar(64),
				createtime VARCHAR (64)
			);
			alter table user default character set utf8;
			alter table user change role role varchar(64) character set utf8;
			alter table user change STATUS status varchar(64) character set utf8;`
	Exec(sqlStr)
	fmt.Println("---------------------------------------------")
	fmt.Println("user table created")
}

func TimeStampToData(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

func CreateSuperAdminInUser(){
	un := "sa"
	psw := MD5("1")
	role := "超级管理员"
	em := "elijah0501@outlook.com"
	tel := 17778340501
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		un, psw, role, em, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("Super Admin created")

}


func CreateAdminInUser(){
	un := "a1"
	psw := MD5("1")
	role := "管理员"
	em := "elijah0501@outlook.com"
	tel := 17778340501
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		un, psw, role, em, tel, st, ct)


	fmt.Println("---------------------------------------------")
	fmt.Println("Admin created")
}

func CreateUser0InUser(){
	un := "u1"
	psw := MD5("1")
	role := "用户"
	em := "elijah0501@outlook.com"
	tel := 17778340501
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		un, psw, role, em, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}

func CreateUser1InUser(){
	un := "u2"
	psw := MD5("1")
	role := "用户"
	em := "elijah0501@outlook.com"
	tel := 17778340501
	st := "异常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		un, psw, role, em, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}

func CreateStaffInUser(){
	un := "s1"
	psw := MD5("1")
	role := "员工"
	em := "elijah0501@outlook.com"
	tel := 17778340501
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		un, psw, role, em, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}


//插入
func InsertUser(user model.User) (int64, error) {
	return Exec("insert into user(username,password,role,email,phone,status,createtime) values (?,?,?,?,?,?,?)",
		user.Username, user.Password, user.Role, user.Email, user.Phone, user.Status, user.Createtime)
}

//按条件查询
func QueryUserWightCon(con string) int {
	sqlStr := fmt.Sprintf("select id from user %s", con)
	fmt.Println(sqlStr)
	row := QueryRowDB(sqlStr)

	fmt.Println("Row is:", row)
	id := 0
	row.Scan(&id)

	log.Println(id)

	fmt.Println("查到的id为", id)
	return id
}

//根据用户名查询id
func QueryUserWithUsername(username string) int {
	sqlStr := fmt.Sprintf("where username='%s'", username)
	return QueryUserWightCon(sqlStr)
}

// 通过 username 和 password 查找 User全部信息
func FindUserByUsernameAndPassword(username string, password string) (user *model.User) {

	var id			int
	var role		string			//0 普通， 1 管理员
	var email		string
	var phone		string
	var status     	string			// 0 正常状态， 1 删除
	var createtime 	string

	sqlStr := fmt.Sprintf("select id, role, email, phone, status, createtime from user where username='%s' and password='%s'", username, password)
	row := QueryRowDB(sqlStr)
	_ = row.Scan(&id, &role, &email, &phone, &status, &createtime)

	user = &model.User{
		Id:			id,
		Username:	username,
		Password:	password,
		Role:		role,
		Email:		email,
		Phone:		phone,
		Status:		status,
		Createtime:	createtime,
	}
	return
}

//查询所有管理员
func QueryAllAdmin() ([]*model.User, error) {
	sqlStr := "select id, username, password, role, email, phone, status, createtime from user where role = '管理员'"

	fmt.Println("--------------------------准备查询所有管理员-------------------")
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)

	if err != nil{
		return nil, err
	}

	fmt.Println("-------------------------创建切片-------------------")
	var admins []*model.User

	for rows.Next() {
		var id			int
		var username	string
		var password	string
		var role		string			//0 普通， 1 管理员
		var email		string
		var phone		string
		var status     	string			// 0 正常状态， 1 删除
		var createtime 	string

		fmt.Println("-------------------------写入行-------------------")
		err := rows.Scan(&id, &username, &password, &role, &email, &phone, &status, &createtime)
		if err != nil {
			return nil, err
		}
		admin := &model.User{
			Id:			id,
			Username:	username,
			Password:	password,
			Role:		role,
			Email:		email,
			Phone:		phone,
			Status:		status,
			Createtime:	createtime,
		}

		admins = append(admins, admin)
	}
	fmt.Println("查询到admin")
	for k, v := range admins {
		fmt.Printf("---%v---%v----\n", k+1, v)
	}
	return admins, nil
}

//查询所有用户
func QueryAllUser() ([]*model.User, error) {
	sqlStr := "select id, username, password, role, email, phone, status, createtime from user where role = '用户'"

	fmt.Println("--------------------------准备查询所有用户-------------------")
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)

	if err != nil{
		return nil, err
	}

	fmt.Println("-------------------------创建切片-------------------")
	var users []*model.User

	for rows.Next() {
		var id			int
		var username	string
		var password	string
		var role		string			//0 普通， 1 管理员
		var email		string
		var phone		string
		var status     	string			// 0 正常状态， 1 删除
		var createtime 	string

		fmt.Println("-------------------------写入行-------------------")
		err := rows.Scan(&id, &username, &password, &role, &email, &phone, &status, &createtime)
		if err != nil {
			return nil, err
		}
		user := &model.User{
			Id:			id,
			Username:	username,
			Password:	password,
			Role:		role,
			Email:		email,
			Phone:		phone,
			Status:		status,
			Createtime:	createtime,
		}

		users = append(users, user)
	}
	fmt.Println("查询到user")
	for k, v := range users {
		fmt.Printf("---%v---%v----\n", k+1, v)
	}
	return users, nil
}

//查询所有职员
func QueryAllStaff() ([]*model.User, error) {
	sqlStr := "select id, username, password, role, email, phone, status, createtime from user where role = '员工'"

	fmt.Println("--------------------------准备查询所有职员-------------------")
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)

	if err != nil{
		return nil, err
	}

	fmt.Println("-------------------------创建切片-------------------")
	var staffs []*model.User

	for rows.Next() {
		var id			int
		var username	string
		var password	string
		var role		string			//0 普通， 1 管理员
		var email		string
		var phone		string
		var status     	string
		var createtime 	string

		fmt.Println("-------------------------写入行-------------------")
		err := rows.Scan(&id, &username, &password, &role, &email, &phone, &status, &createtime)
		if err != nil {
			return nil, err
		}
		staff := &model.User{
			Id:			id,
			Username:	username,
			Password:	password,
			Role:		role,
			Email:		email,
			Phone:		phone,
			Status:		status,
			Createtime:	createtime,
		}

		staffs = append(staffs, staff)
	}
	fmt.Println("查询到staff")
	for k, v := range staffs {
		fmt.Printf("---%v---%v----\n", k+1, v)
	}
	return staffs, nil
}
