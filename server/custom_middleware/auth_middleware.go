package custommiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/momoh-yusuf/note-app/utils"
)

// HTTP middleware setting a value on the request context
func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]

		claims, err := utils.ValidateJwtToken(token)

		if err != nil {
			utils.CustomResponseInJson(w, http.StatusBadRequest, "Error occurred", err)
			return
		}

		// lets verify token
		// isIssuerValid, err := jwt.MapClaims.GetExpirationTime(jwt.ErrTokenExpired)

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		authUser := claims

		fmt.Println(claims)
		ctx := context.WithValue(r.Context(), "user", authUser)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
