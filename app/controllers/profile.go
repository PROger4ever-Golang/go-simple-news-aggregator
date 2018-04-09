package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"

	"github.com/PROger4ever/go-simple-news-aggregator/app/models"
	"github.com/PROger4ever/go-simple-news-aggregator/app/routes"
)

type Profile struct {
	Application
}

func (c Profile) Settings() revel.Result {
	return c.Render()
}

func (c Profile) SaveSettings(password, verifyPassword string) revel.Result {
	models.ValidatePassword(c.Validation, password)
	c.Validation.Required(verifyPassword).
		Message("Please verify your password")
	c.Validation.Required(verifyPassword == password).
		Message("Your password doesn't match")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(routes.Profile.Settings())
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := c.Txn.ExecUpdate(c.Db.SqlStatementBuilder.
		Update("User").Set("HashedPassword", bcryptPassword).
		Where("UserId=?", c.connected().UserId))
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Password updated")
	return c.Redirect(routes.Sources.Index(20, 1))
}
