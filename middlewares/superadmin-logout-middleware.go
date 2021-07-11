package middlewares

import (
	"elect/dto"
	"elect/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func SuperAdminLogoutMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

		cookie, err := cxt.Cookie("token")
		if err != nil {
			cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: "Not Logged in!",
			})
			return
		}

		value := make(map[string]string)
		err = s.Decode("tokens", cookie, &value)
		if err != nil {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/")
			return
		}

		_, role, err := jwtService.GetUserIDAndRole(value["access_token"])
		if err != nil {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/")
			return
		}

		if role != 2 {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/")
			return
		}

		return
	}
}
