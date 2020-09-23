package dao

import (
	"blc-iot-demo/web/model"
	"fmt"
	"time"
)

//Create action table
func CreateTableWithAction() {
	sqlStr := `CREATE TABLE IF NOT EXISTS action (
				action_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
				user_id BIGINT NOT NULL,
				action_type VARCHAR (64) NOT NULL,
				action_target_name VARCHAR (64) NOT NULL,
				action_target_id VARCHAR (64) NOT NULL,
				action_time VARCHAR (64) NOT NULL,
				FOREIGN KEY (user_id) REFERENCES user (id)
			);
			alter table action default character set utf8;
			alter table action change action_type action_type varchar(64) character set utf8;
			alter table action change action_target_name action_target_name varchar(64) character set utf8;
			alter table action change action_target_id action_target_id varchar(64) character set utf8;
			alter table action change action_time action_time varchar(64) character set utf8;`
	Exec(sqlStr)
	fmt.Println("---------------------------------------------")
	fmt.Println("Action table created")
}

//插入
func InsertAction(action model.Action) (int64, error) {
	return Exec("insert into action(user_id,action_type,action_target_name,action_target_id,action_time) values (?,?,?,?,?)",
		action.UserID, action.ActionType, action.ActionTargetName, action.ActionTargetID, action.ActionTime)
}

//遍历
func QueryAllAction(userID int) ([]*model.Action, error) {
	sqlStr := fmt.Sprintf("select action_id, user_id, action_type, action_target_name, action_target_id, action_time from action where user_id='%d'", userID)

	fmt.Println("--------------------------准备查询所有职员-------------------")
	fmt.Println(sqlStr)
	rows, err := db.Query(sqlStr)

	if err != nil {
		return nil, err
	}

	fmt.Println("-------------------------创建切片-------------------")
	var actions []*model.Action

	for rows.Next() {
		var actionID int
		var userID int
		var actionType string
		var actionTargetName string
		var actionTargetID string
		var actionTime string

		fmt.Println("-------------------------写入行-------------------")
		err := rows.Scan(&actionID, &userID, &actionType, &actionTargetName, &actionTargetID, &actionTime)
		if err != nil {
			return nil, err
		}
		action := &model.Action{
			ActionID:         actionID,
			UserID:           userID,
			ActionType:       actionType,
			ActionTargetName: actionTargetName,
			ActionTargetID:   actionTargetID,
			ActionTime:       actionTime,
		}

		actions = append(actions, action)
	}
	fmt.Println("查询到staff")
	for k, v := range actions {
		fmt.Printf("---%v---%v----\n", k+1, v)
	}
	return actions, nil
}

func CreateSuperAdminActionInAction() {
	userID := 1
	actionType := "任命管理员"
	actionTargetName := "a1"
	actionTargetID := "2"
	actionTime := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into action(user_id,action_type,action_target_name,action_target_id,action_time) values (?,?,?,?,?)",
		userID, actionType, actionTargetName, actionTargetID, actionTime)

	fmt.Println("---------------------------------------------")
	fmt.Println("Super Admin Action created")

}

func CreateAdminActionInAction() {
	userID := 2
	actionType := "任命员工"
	actionTargetName := "s1"
	actionTargetID := "5"
	actionTime := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into action(user_id,action_type,action_target_name,action_target_id,action_time) values (?,?,?,?,?)",
		userID, actionType, actionTargetName, actionTargetID, actionTime)

	fmt.Println("---------------------------------------------")
	fmt.Println("Admin Action created")
}

func CreateUserActionInAction() {
	userID := 3
	actionType := "查询"
	actionTargetName := "川红工夫"
	actionTargetID := "54df5g4d5g68tuk466a8wg46"
	actionTime := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into action(user_id,action_type,action_target_name,action_target_id,action_time) values (?,?,?,?,?)",
		userID, actionType, actionTargetName, actionTargetID, actionTime)

	fmt.Println("---------------------------------------------")
	fmt.Println("User Action created")
}

func CreateStaffActionInAction() {
	userID := 5
	actionType := "上链"
	actionTargetName := "川红工夫"
	actionTargetID := "54df5g4d5g68tuk466a8wg46"
	actionTime := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into action(user_id,action_type,action_target_name,action_target_id,action_time) values (?,?,?,?,?)",
		userID, actionType, actionTargetName, actionTargetID, actionTime)

	fmt.Println("---------------------------------------------")
	fmt.Println("Staff Action created")
}
