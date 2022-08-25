package ctx

import (
	"context"
	"literate-barnacle/service/models"
	"net/http"
	"time"

	jwt "github.com/Viva-Victoria/bear-jwt"
)

type ContextProvider func() context.Context
type TimeoutContextProvider func(duration time.Duration) (context.Context, context.CancelFunc)

type Context struct {
	context.Context
	Authorized    bool
	Authorization models.Authorization
}

func GetContext(r *http.Request, authorization bool) (Context, error) {
	ctx := Context{
		Context: r.Context(),
	}

	if authorization {
		authorizationHeader := r.Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			return Context{}, ErrUnauthorized
		}

		token, err := jwt.Parse([]byte(authorizationHeader[7:]))
		if err != nil {
			return Context{}, NewParseTokenError(err)
		}

		switch token.ValidateNow() {
		case jwt.StateValid:
			claims := models.TokenClaims{}
			if err = token.Claims.Get(&claims); err != nil {
				return Context{}, NewParseTokenError(err)
			}

			ctx.Authorized = true
			ctx.Authorization = claims.Authorization
		case jwt.StateExpired, jwt.StateNotIssued, jwt.StateInactive:
			return Context{}, ErrTokenExpired
		default:
			return Context{}, NewParseTokenError(err)
		}
	}

	return ctx, nil
}
