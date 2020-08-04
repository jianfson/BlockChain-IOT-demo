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
	admin := User{LoginName: "0", Password: "0", IsAdmin: "T"}
	jack := User{LoginName: "1", Password: "1", IsAdmin: "F"}
	alice := User{LoginName: "2", Password: "2", IsAdmin: "T"}
	bob := User{LoginName: "3", Password: "3", IsAdmin: "F"}

	users = append(users, admin)
	users = append(users, alice)
	users = append(users, bob)
	users = append(users, jack)

}
