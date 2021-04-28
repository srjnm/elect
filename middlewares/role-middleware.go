package middlewares

import (
	"elect/dto"
	"elect/roles"
	"elect/services"
	"net/http"
	"os"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func Authorizer(jwtService services.JWTService, e *casbin.Enforcer) gin.HandlerFunc {
	return func(cxt *gin.Context) {
		role := ""

		var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

		cookie, err := cxt.Cookie("otp")
		if err == nil {
			value := make(map[string]string)
			err = s.Decode("token", cookie, &value)

			_, err := jwtService.ValidateOTPToken(value["otp_token"])
			if err != nil {
				if err.Error() == "Token is expired" {
					cxt.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
						Message: "OTP time expired!",
					})
					return
				} else {
					cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
						Message: "Invalid Request",
					})
					return
				}
			}

			role = strconv.Itoa(roles.Authenticated)
		}

		if role == "" {
			cookie, err = cxt.Cookie("token")
			if err == http.ErrNoCookie {
				role = strconv.Itoa(roles.Anonymous)
			}
			if err != nil && err != http.ErrNoCookie {
				cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
					Message: "Invalid Request",
				})
				return
			}
		}

		if role == "" {
			value := make(map[string]string)
			err = s.Decode("tokens", cookie, &value)

			roleInt, err := jwtService.GetRole(value["access_token"])
			if err != nil {
				if err.Error() == "Token is expired" {
					cxt.AbortWithStatusJSON(http.StatusNotAcceptable, dto.Response{
						Message: "Refresh Required",
					})
					return
				} else {
					cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
						Message: "Invalid Request",
					})
					return
				}
			}

			role = strconv.Itoa(roleInt)
		}

		res, err := e.Enforce(role, cxt.Request.URL.Path, cxt.Request.Method)
		if err != nil {
			cxt.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Message: "Invalid Request",
			})
			return
		}

		if !res {
			if role == "-1" || role == "-2" {
				cxt.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, dto.Response{
					Message: "Not logged in!",
				})
				return
			}

			cxt.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
				Message: "Forbidden",
			})
			return
		}

		return
	}
}
