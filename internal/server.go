package internal

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Start(port string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/user", newUser)
	e.GET("/user/:email", getUser)
	e.POST("/delete-user/:email", deleteUser)
	e.Logger.Fatal(e.Start(":" + port))
}

type (
	User struct {
		Name  string `json:"name"  form:"name"  query:"name"  validate:"required"`
		Email string `json:"email" form:"email" query:"email" validate:"required,email"`
	}
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func newUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(u); err != nil {
		return err
	}

	user := UserInternal{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: false,
	}
	AddUserToDB(&user)

	return c.JSON(http.StatusOK, u)
}

func getUser(c echo.Context) (err error) {
	email := c.Param("email")
	user := GetUserFromDB(email)

	return c.String(http.StatusOK, fmt.Sprintf("{name: %s, email: %s}", user.Name, user.Email))
}

func deleteUser(c echo.Context) (err error) {
	email := c.Param("email")
	DeleteUserFromDB(email)

	return c.String(http.StatusOK, email)
}
