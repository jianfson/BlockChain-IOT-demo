package dao

import (
	"blc-iot-demo/web/model"
	"fmt"
)

//Create session table
func CreateTableWithSession() {
	sqlStr := `CREATE TABLE IF NOT EXISTS session (
				session_id VARCHAR (64) PRIMARY KEY NOT NULL,
				user_id BIGINT NOT NULL,
				username VARCHAR (64) NOT NULL,
				PASSWORD VARCHAR (64) NOT NULL,
				role VARCHAR (64) NOT NULL,
				email VARCHAR (64) NOT NULL,
				phone VARCHAR (64) NOT NULL,
				STATUS VARCHAR (64) NOT NULL,
				createtime VARCHAR (64) NOT NULL,
				FOREIGN KEY (user_id) REFERENCES user (id)
			);
			alter table session default character set utf8;
			alter table session change role role varchar(64) character set utf8;
			alter table session change status status varchar(64) character set utf8;`
	Exec(sqlStr)
	fmt.Println("---------------------------------------------")
	fmt.Println("session table created")
}

// 添加session记录
func AddSession(sess *model.Session) error{
	sqlStr := `INSERT INTO session VALUES(?,?,?,?,?,?,?,?,?);`
	fmt.Println("---------------------------------------------")
	fmt.Println("正在写入Session表")
	_, err := Exec(sqlStr, sess.SessionID, sess.UserID, sess.UserName, sess.PassWord, sess.Role, sess.Email, sess.Phone, sess.Status, sess.CreateTime)
	if err != nil{
		return err
	}
	return nil
}

func GetSession(sessID string) (sess *model.Session) {
	var sessionID	string		//存到cookie中
	var userID		int			//设置外键
	var userName	string
	var passWord	string
	var role		string
	var email		string
	var phone		string
	var status    	string
	var createTime 	string

	sqlStr := fmt.Sprintf("select session_id, user_id, username, password, role, email, phone, status, createtime from session where session_id='%s'", sessID)
	row := QueryRowDB(sqlStr)
	_ = row.Scan(&sessionID, &userID, &userName, &passWord, &role, &email, &phone, &status, &createTime)
	sess = &model.Session{
		SessionID:		sessionID,
		UserID:			userID,
		UserName:		userName,
		PassWord:		passWord,
		Role:			role,
		Email:			email,
		Phone:			phone,
		Status:			status,
		CreateTime:		createTime,
	}
	return
}


// 删除session记录
func DeleteSession(sessID string)  error{
	sqlStr := `DELETE FROM session WHERE session_id = ?;`
	_, err := Exec(sqlStr, sessID)
	if err != nil{
		return err
	}
	return nil
}
