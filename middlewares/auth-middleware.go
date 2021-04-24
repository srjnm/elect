package middlewares

import (
	"elect/dto"
	"elect/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func Authorization(jwtService services.JWTService) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

		cookie, err := cxt.Cookie("token")
		if err != nil {
			cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: "Invalid Request",
			})
			return
		}

		value := make(map[string]string)
		err = s.Decode("tokens", cookie, &value)
		if err != nil {
			cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: "Invalid Request",
			})
			return
		}

		accessToken, err := jwtService.ValidateAccessToken(value["access_token"])
		if err != nil && err.Error() != "Token is expired" {
			cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Session Invalid",
			})
			return
		}
		if !accessToken.Valid {
			refreshToken, err := jwtService.ValidateRefreshToken(value["refresh_token"])
			if err != nil && err.Error() != "Token is expired" {
				cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
					Message: "Session Invalid",
				})
				return
			}

			if !refreshToken.Valid {
				cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
					Message: "Session Expired",
				})
				return
			}

			cxt.AbortWithStatusJSON(http.StatusNotAcceptable, dto.Response{
				Message: "Refresh Required",
			})
			return
		}

		return
	}
}
