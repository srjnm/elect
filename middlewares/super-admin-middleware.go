package middlewares

import (
	"elect/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func SuperAdminMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

		cookie, err := cxt.Cookie("token")
		if err != nil {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/logout")
			return
		}

		value := make(map[string]string)
		err = s.Decode("tokens", cookie, &value)
		if err != nil {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/logout")
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

		accessToken, err := jwtService.ValidateAccessToken(value["access_token"])
		if err != nil && err.Error() != "Token is expired" {
			cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/logout")
			return
		}
		if !accessToken.Valid {
			refreshToken, err := jwtService.ValidateRefreshToken(value["refresh_token"])
			if err != nil && err.Error() != "Token is expired" {
				cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/logout")
				return
			}

			if !refreshToken.Valid {
				cxt.Redirect(http.StatusTemporaryRedirect, "https://e1ect.herokuapp.com/logout")
				return
			}
		}

		return
	}
}
