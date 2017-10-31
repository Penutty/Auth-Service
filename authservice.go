package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/penutty/authservice/user"
	"github.com/penutty/util"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	e := echo.New()

	e.POST("/user/:user/create", createUser)
	e.POST("/auth", authUser)

	e.Logger.Fatal(e.Start(":8080"))
}

// createUser is a POST endpoint that accepts
// Content-Type: [application/json; charset=UTF-8]
// Body: {
//			UserID: UserID
//			Email: Email
//			Password: Password
//		 }
// on success returns
// Status: 201 - Created
func createUser(c echo.Context) error {
	eq, err := util.DeepEqual(c.ParamNames(), []string{"UserID", "Email", "Password"})
	if err != nil {
		c.Logger().Printf("Error: %v\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if !eq {
		c.Logger().Printf("Request Body Invalid. StatusCode: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	u := user.NewUser(c.Param("UserID"), c.Param("Email"), c.Param("Password"))
	user.Create(u, user.MomentDB())

	switch u.Err() {
	case nil:
		return c.NoContent(http.StatusCreated)
	default:
		return c.NoContent(http.StatusInternalServerError)
	}
}

// authUser is a POST endpoint that accepts
// Body: {
//			UserID: UserID
//			Password: Password
//		 }
// on success returns
// Status: 200
func authUser(c echo.Context) error {
	eq, err := util.DeepEqual(c.ParamNames(), []string{"UserID", "Password"})
	if !eq {
		c.Logger().Printf("Request Body Invalid. StatusCode: %v", http.StatusBadRequest)
		return c.NoContent(http.StatusBadRequest)
	}

	u := user.Fetch(c.Param("UserID"), user.MomentDB())
	switch u.Err() {
	case nil:
		break
	default:
		return c.NoContent(http.StatusInternalServerError)
	}

	if u.Password() != c.Param("Password") {
		c.Logger().Printf("Failed Login - User: %v", c.Param("UserID"))
		return c.NoContent(http.StatusUnauthorized)
	}

	token, err := generateJwt(c.FormValue("UserID"))
	if err != nil {
		c.Logger().Printf("JWT failed to be generated. ERROR: %v\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("jwt", token)
	return c.NoContent(http.StatusOK)

}

// generateJwt uses a requests UserID and a []byte secret to generate a JSON web token.
func generateJwt(UserID string) (string, error) {

	p, err := ioutil.ReadFile("/home/tjp/.ssh/jwt_private.pem")
	if err != nil {
		return "", err
	}

	t := jwt.New(jwt.SigningMethodRS256)
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	claims["iss"] = "Auth-Service"
	claims["sub"] = UserID
	claims["aud"] = "Moment-Service"
	claims["exp"] = time.Now().UTC().AddDate(0, 0, 7).Unix()
	claims["iat"] = time.Now().UTC().Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM(p)
	if err != nil {
		return "", err
	}

	token, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}
