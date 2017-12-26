package main

import (
	"fmt"
	"orm-engine/entities"
	"reflect"
)

// user Dao
//通过组合模式继承
type IUserDao interface {
	entities.IBaseDao
}

type userDao struct {
	entities.BaseDao
}

// 初始化dao

var userDaoImpl IUserDao

func UserDao() IUserDao {
	if userDaoImpl == nil {
		userDaoImpl = &userDao{entities.BaseDao{EntityType: reflect.TypeOf(new(entities.UserInfo)).Elem()}}
		userDaoImpl.Init()
	}
	return userDaoImpl
}

func main() {
	userDaoImpl = UserDao()
	var u = entities.UserInfo{UserName: "sysu0dd"}
	user := entities.NewUserInfo(u)
	user.DepartName = "depart0001"
	err := userDaoImpl.Save(user)
	if err != nil {
		panic(err)
	}
	//data := make([]byte, 0)
	pEveryOne, err := userDaoImpl.Find()
	for index := 0; index < pEveryOne.Len(); index++ {
		item := reflect.ValueOf(pEveryOne.Front().Value).Interface()
		fmt.Printf("result: %v\n", reflect.ValueOf(item))
		pEveryOne.Remove(pEveryOne.Front())
	}

}
