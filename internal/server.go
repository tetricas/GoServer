package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Start(port string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", newUser)
	e.Logger.Fatal(e.Start(":" + port))
}

type User struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func newUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return
	}

	user := UserInternal{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: false,
	}
	AddUserToDB(&user)

	return c.JSON(http.StatusOK, u)
}
