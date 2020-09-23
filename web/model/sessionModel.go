package model

type Session struct {
	SessionID	string		//存到cookie中
	UserID		int			//设置外键
	UserName	string
	PassWord	string
	Role		string
	Email		string
	Phone		string
	Status    	string
	CreateTime 	string
}
