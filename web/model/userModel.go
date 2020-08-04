package model

type User struct {
	Id				int
	Username		string
	Password		string
	//Logincounter	int64
	Role			string			//0 普通， 1 管理员
	Email			string
	Phone			string
	Status     		string			// 0 正常状态， 1 删除
	Createtime 		string
}
