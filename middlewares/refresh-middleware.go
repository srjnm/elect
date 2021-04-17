package middlewares

import (
	"elect/dto"
	"elect/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func EnsureValidity(jwtService services.JWTService) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

		cookie, err := cxt.Cookie("token")
		if err != nil {
			cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized User",
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

		accessToken, err := jwtService.ValidateAccessToken(value["access_token"])
		if err != nil && err.Error() != "Token is expired" {
			cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Session Invalid",
			})
			return
		}
		if accessToken.Valid {
			cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: "Session still valid!",
			})
			return
		}
	}
}
