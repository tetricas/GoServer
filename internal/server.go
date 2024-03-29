package internal

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Start(port string) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	admin := e.Group("/user")
	admin.POST("/", newUser)
	admin.GET("/:email", getUser)
	admin.POST("/:email", deleteUser)
	admin.Use(middleware.BasicAuth(login))

	e.Logger.Fatal(e.Start(":" + port))
}

type (
	User struct {
		Name     string `json:"name"  form:"name"  query:"name"  validate:"required"`
		Email    string `json:"email" form:"email" query:"email" validate:"required,email"`
		Password string `json:"password" form:"password" query:"password" validate:"required"`
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	user := UserInternal{
		Name:    u.Name,
		Email:   u.Email,
		Secret:  string(hashedPassword),
		IsAdmin: false,
	}
	AddUserToDB(&user)

	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) (err error) {
	email := c.Param("email")
	user := GetUserFromDB(email)
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{name: %s, email: %s}", user.Name, user.Email))
}

func deleteUser(c echo.Context) (err error) {
	email := c.Param("email")
	DeleteUserFromDB(email)

	return c.String(http.StatusOK, email)
}

func login(username, password string, c echo.Context) (bool, error) {
	user := GetUserFromDB(username)
	if user == nil {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
