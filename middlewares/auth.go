package middlewares

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt2"
	log "github.com/sirupsen/logrus"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
	"github.com/amir-yaghoobi/accountManager/config"
)


type LoginRequest struct {
	Username string    `json:"username"   binding:"required"`
	Password string    `json:"password"   binding:"required"`
}


func UserIdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	userIdFloat := claims["userID"].(float64)
	userId := int(userIdFloat)
	pConn, _ := db.GetPostgres()

	user := models.User{}
	pConn.Preload("Accounts").First(&user, userId)

	c.Set("USER", &user)
	return &user
}

func UserPayloadBuilder(data interface{}) jwt.MapClaims {
	if user, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"userID": user.ID,
		}
	}
	return jwt.MapClaims{}
}

func UserAuthentication(c *gin.Context) (interface{}, error) {
	var loginForm LoginRequest
	if err := c.ShouldBind(&loginForm); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	username := loginForm.Username
	password := loginForm.Password

	pConn, err := db.GetPostgres()
	if err != nil {
		return nil, err
	}

	user := models.User{}
	queryResult := pConn.Where("Username = ?", username).First(&user)
	if queryResult.Error != nil && queryResult.Error != gorm.ErrRecordNotFound {
		log.Errorf("error on fetching username: \"%s\", error: %s\n", username,
			queryResult.Error.Error())

		return nil, queryResult.Error
	} else if queryResult.Error == gorm.ErrRecordNotFound {
		log.Warnf("attempting to login with username: \"%s\", but username does not exist in database!", username)
		return nil, jwt.ErrFailedAuthentication
	}

	isMatch := user.CheckPassword(password)
	if !isMatch {
		log.Warnf("attempting to login with username: \"%s\", invalid password: \"%s\"\n",
			username, password)

		return nil, jwt.ErrFailedAuthentication
	}

	return &user, nil
}

func UserAuthorization(data interface{},c *gin.Context) bool {
	//if v, ok := data.(*User); ok && v.UserName == "admin" {
	//	return true
	//}
	//
	//return false
	return true
}

func LoginResponseHandler(context *gin.Context, status int, token string, expire time.Time) {
	context.JSON(status, gin.H{
		"status":   status,
		"token":    token,
		"expireAt": expire,
	})
}

func SetupJWT() *jwt.GinJWTMiddleware {
	cfg := config.GetConfig()
	return &jwt.GinJWTMiddleware{
		Realm:            "development",
		SigningAlgorithm: "HS256",
		Key:              []byte(cfg.SecretKey),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 48,
		IdentityHandler:  UserIdentityHandler,
		PayloadFunc:      UserPayloadBuilder,
		Authenticator:    UserAuthentication,
		Authorizator:     UserAuthorization,
		LoginResponse:    LoginResponseHandler,
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			return e.Error()
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}