package router

import (
	util "duke/init/src/helpers"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func Init() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/forgotPassword", forgotPasswordHandler)
	http.HandleFunc("/resetPassword", resetPasswordHandler)
	http.Handle("/profile", isAuthorized(profileHandler))
}

var loginHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"login"}`))
}

var signUpHandler = func(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	emailId := req.Form.Get("emailId")

	if len(username) > 4 && len(password) > 8 && len(emailId) > 5 {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("incomplet data recived", "param is not valid", nil)
		w.Write(resp)
		return
	}
	fmt.Println(username, password, emailId)
	w.Header().Set("Content-Type", "application/json")
	//TODO: - replace bellow 1 with userId from db
	validToken, err := GetJWT(username, 1)
	fmt.Println(validToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("please try after sometime", "Failed to generate token", err)
		w.Write(resp)
	}

	w.Write([]byte(`{"token":` + string(validToken) + `}`))
	return
}

var forgotPasswordHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"signup"}`))
}

var resetPasswordHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"signup"}`))
}

var profileHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"profileHandler"}`))
}

func GetJWT(username string, userId int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["userId"] = userId
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(util.EnvConfig().SecretKey)

	if err != nil {
		util.LogError("Something Went Wrong:", err)
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		w.Header().Set("Content-Type", "application/json")
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					w.WriteHeader(http.StatusUnauthorized)
					resp := util.ErrorResponse("invalid access", "Invalid tokend", nil)
					w.Write(resp)
				}
				if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
					w.WriteHeader(http.StatusUnauthorized)
					resp := util.ErrorResponse("invalid access", "Expired token", nil)
					w.Write(resp)
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					w.WriteHeader(http.StatusUnauthorized)
					resp := util.ErrorResponse("invalid access", "invalid token", nil)
					w.Write(resp)
				}
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					w.WriteHeader(http.StatusUnauthorized)
					resp := util.ErrorResponse("invalid access", "invalid token", nil)
					w.Write(resp)
				}

				return util.EnvConfig().SecretKey, nil
			})
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
