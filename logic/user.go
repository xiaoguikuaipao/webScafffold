package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) {
	//0. estimate the existent of user
	mysql.QueryUserByUsername()
	//1. generate the UID
	snowflake.GenID()
	//2. stores in the db
	mysql.InsertUser()
}
