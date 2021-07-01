package middlewares

import (
	"elect/services"
	"errors"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/qor/admin"
)

func AdminMiddleware(jwtService services.JWTService) *admin.Middleware {
	return &admin.Middleware{
		Name: "refresh",
		Handler: func(context *admin.Context, middleware *admin.Middleware) {
			reqCookie, err := context.Request.Cookie("token")
			if err != nil {
				context.AddError(errors.New("Bad Request!"))
				return
			}

			var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
			value := make(map[string]string)
			err = s.Decode("tokens", reqCookie.String(), &value)
			if err != nil {
				context.AddError(errors.New("Bad Request!"))
				return
			}

			accessToken, err := jwtService.ValidateAccessToken(value["access_token"])
			if err != nil && err.Error() != "Token is expired" {
				context.AddError(errors.New("Session Invalid!"))
				return
			}

			if !accessToken.Valid {
				refreshToken, err := jwtService.ValidateRefreshToken(value["refresh_token"])
				if err != nil && err.Error() != "Token is expired" {
					context.AddError(errors.New("Session Invalid!"))
					return
				}

				if !refreshToken.Valid {
					context.AddError(errors.New("Session Expired!"))
					return
				}

				return
			}
		},
	}
}
