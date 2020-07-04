package controller

import "fab-sdk-go-sample/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	IsAdmin   string
}

var users []User

func init() {
	admin := User{LoginName: "Hanxiaodong", Password: "123456", IsAdmin: "T"}
	alice := User{LoginName: "ChainDesk", Password: "123456", IsAdmin: "T"}
	bob := User{LoginName: "alice", Password: "123456", IsAdmin: "F"}
	jack := User{LoginName: "1", Password: "1", IsAdmin: "F"}

	users = append(users, admin)
	users = append(users, alice)
	users = append(users, bob)
	users = append(users, jack)

}
