package login

import (
	util "duke/init/src/helpers"
	"duke/init/src/login/database"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Config struct {
	Database               *mongo.Database
	CollectionName         string
	Aud                    string
	Iss                    string
	ForgotPasswordCallback func(emailId string, token string)
}

var dbConfig database.LoginDBConfig

func (c *Config) Init() {
	dbConfig.CollectionName = c.CollectionName
	dbConfig.Database = c.Database
	dbConfig.Iss = c.Iss
	dbConfig.Aud = c.Aud
	dbConfig.Init()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/forgotPassword", c.forgotPasswordHandler)
	http.HandleFunc("/resetPassword", resetPasswordHandler)
}

var loginHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp := util.ErrorResponse("invalid param", "login not valid", err)
		_, _ = w.Write(resp)
		return
	}
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	util.LogInfo(username, password)
	userInfo, err := dbConfig.FindUser(username)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp := util.ErrorResponse("invalid access", "login not valid", err)
		_, _ = w.Write(resp)
		return
	}

	hashPassword := []byte(userInfo["password"].(string))

	if !isPasswordValid(hashPassword, []byte(password)) {
		w.WriteHeader(http.StatusUnauthorized)
		resp := util.ErrorResponse("invalid access", "not matching", nil)
		_, _ = w.Write(resp)
		return
	}

	objId := userInfo["_id"].(primitive.ObjectID)
	validToken, err := GetJWT(username, objId)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp := util.ErrorResponse("please try after sometime", "Failed to generate token", err)
		_, _ = w.Write(resp)
		return
	}

	data := `{
			"id":"` + objId.Hex() + `",
			"username":"` + userInfo["username"].(string) + `",
			"token":"` + validToken + `"
			}`
	resp := util.SuccessResponse(data)
	_, _ = w.Write(resp)
}

var signUpHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()
	if err != nil {
		resp := util.ErrorResponse("invalid param", "param is not valid", err)
		_, _ = w.Write(resp)
		return
	}
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	emailId := req.Form.Get("emailId")

	if len(username) < 4 && len(password) < 8 && len(emailId) < 5 {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("incomplete data", "param is not valid", nil)
		_, _ = w.Write(resp)
		return
	}
	user := database.User{}
	user.Username = username
	user.EmailId = emailId
	user.Password = getHash([]byte(password))
	objId, err := dbConfig.CreateUser(user)

	if err != nil {
		util.LogError("", err)
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("please try after sometime", "unable to create user in db", err)
		_, _ = w.Write(resp)
		return
	}

	util.LogInfo(user.Password)
	user.Password = getHash([]byte(user.Password))
	validToken, err := GetJWT(username, objId)
	fmt.Println(validToken)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("please try after sometime", "Failed to generate token", err)
		_, _ = w.Write(resp)
		return
	}
	resp := util.SuccessResponse(`{"token":"` + string(validToken) + `"}`)
	_, _ = w.Write(resp)
	return
}

func (c *Config) forgotPasswordHandler (w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()
	if err != nil {
		resp := util.ErrorResponse("invalid param", "param incorrect", err)
		_, _ = w.Write(resp)
		return
	}
	emailId := req.Form.Get("emailId")
	util.LogInfo(emailId,"====>")
	if dbConfig.IsUserValid(emailId) {
		token,err := resetToken(emailId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := util.ErrorResponse("please try after sometime", "Failed to generate token", err)
			_, _ = w.Write(resp)
		}

		c.ForgotPasswordCallback(emailId,token)
		resp := util.SuccessResponse("null")
		_, _ = w.Write(resp)
		return
		return
	}else{
		resp := util.ErrorResponse("invalid emailId","emailId not found",nil)
		_, _ = w.Write(resp)
		return
	}

}

var resetPasswordHandler = func(w http.ResponseWriter, req *http.Request) {
	fmt.Print("in 1")
	w.Header().Set("Content-Type", "application/json")
		err := req.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := util.ErrorResponse("invalid param", "login not valid", err)
			_, _ = w.Write(resp)
			return
		}
		tokenString := req.Form.Get("token")
		password := req.Form.Get("password")
		if len(password) < 8{
			w.WriteHeader(http.StatusUnauthorized)
			resp := util.ErrorResponse("invalid password", "not a strong password", err)
			_, _ = w.Write(resp)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(util.EnvConfig().SecretKey), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := util.ErrorResponse("invalid token", "token not valid", err)
			_, _ = w.Write(resp)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		emailId := claims["emailId"].(string)

		if dbConfig.IsUserValid(emailId){
			passwordHash := getHash([]byte(password))
			if dbConfig.UpdatePassword(emailId,passwordHash){
				resp := util.SuccessResponse(`{"response":"password successfully"}`)
				_, _ = w.Write(resp)
				return
			}
		}else{
			w.WriteHeader(http.StatusUnauthorized)
			resp := util.ErrorResponse("invalid emailId", "not a valid emailId", err)
			_, _ = w.Write(resp)
			return
		}
}

///-------------------///

func GetJWT(username string, userId primitive.ObjectID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["userId"] = userId
	claims["aud"] = dbConfig.Aud
	claims["iss"] = dbConfig.Iss
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(util.EnvConfig().SecretKey)

	if err != nil {
		util.LogError("Something Went Wrong:", err)
		return "", err
	}

	return tokenString, nil
}

func resetToken(emailId string) (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["emailId"] = emailId
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, err := token.SignedString(util.EnvConfig().SecretKey)
		if err != nil {
			util.LogError("Something Went Wrong:", err)
			return "", err
		}
		return tokenString, nil
	}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		util.LogError("could not creat hash", err)
	}
	return string(hash)
}

func isPasswordValid(hash []byte, password []byte) bool {
	if bcrypt.CompareHashAndPassword(hash, password) != nil {
		return false
	}
	return true
}
