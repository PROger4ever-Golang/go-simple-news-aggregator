package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Sources.checkUser, revel.BEFORE)
	revel.InterceptMethod(Articles.checkUser, revel.BEFORE)
}
