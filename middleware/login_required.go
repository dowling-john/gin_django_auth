package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flexin1981/gin_django_auth/datalayer"
	"github.com/gin-gonic/gin"
)

const (
	DjangoSessionCookie string = "sessionid"
)

func sessionTokenExistsInRequest(c *gin.Context) bool {
	if _, err := c.Cookie(DjangoSessionCookie); err != nil {
		log.Println(DjangoSessionCookieNotFound)
		return false
	}
	return true
}

func sessionTokenExistsInDatabase(c *gin.Context, d datalayer.SessionServiceInterface) bool {
	cookie, _ := c.Cookie(DjangoSessionCookie)
	if _, err := d.Get(cookie); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func sessionTokenBlank(c *gin.Context) bool {
	cookie, _ := c.Cookie(DjangoSessionCookie)
	if cookie == "" {
		log.Println(DjangoSessionCookieIsBlank)
		return true
	}
	return false
}

// LoginRequired uses a cookie named "sessionid", this session id is compared to the django_sessions table
// If the session is in the table the user is considered as authenticated.
func LoginRequired(c *gin.Context) {
	d := datalayer.NewSessionService()
	if sessionTokenExistsInRequest(c) && !sessionTokenBlank(c) && sessionTokenExistsInDatabase(c, d) {
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedJson)
		return
	}
}
