package dbrepo

import (
	"testing"

	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
)

func TestCreateUser(t *testing.T) {
	email := util.RandomString(6) + "@" + util.RandomString(2) + "." + util.RandomString(2)
	avatar := "http://" + util.RandomString(6)
	name := util.RandomString(4)
	password := util.RandomString(6)
	user := models.User{
		Email:    email,
		Avatar:   avatar,
		Name:     name,
		Password: password,
	}

	err := testRepo.CreateUser(&user)
	if err != nil {
		t.Error(err)
	}
	if user.ID <= 0 {
		t.Error("create user failed, id=", user.ID)
	}
}
