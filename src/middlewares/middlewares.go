package middleware

import (
	util "duke/init/src/helpers"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

//EnableCors :- enable cors

var env = util.EnvConfig()

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With"
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		// for name, values := range r.Header {
		// 	// Loop over all values for the name.
		// 	for _, value := range values {
		// 		fmt.Println(name, value)
		// 	}
		// }

		next.ServeHTTP(w, r)
	})
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		w.Header().Set("Content-Type", "application/json")
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], tokenCheck)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				resp := util.ErrorResponse("invalid access", "invalid token", err)
				w.Write(resp)
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			resp := util.ErrorResponse("invalid access", "token not found", nil)
			w.Write(resp)
			return
		}
	})
}

func tokenCheck(token *jwt.Token) (interface{}, error) {
	return util.EnvConfig().SecretKey, nil
}
