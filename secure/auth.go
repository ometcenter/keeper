package secure

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	shareRedis "github.com/ometcenter/keeper/redis"
	shareStore "github.com/ometcenter/keeper/store"
	web "github.com/ometcenter/keeper/web"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyCustomClaimsLight struct {
	UserId uint `json:"userid"`
	jwt.StandardClaims
}

func GetMiddleWareLight() gin.HandlerFunc {
	return func(c *gin.Context) {

		// requestPath := c.Request.URL.Path          //current request path

		tokenHeader := c.GetHeader("TokenBearer")

		if tokenHeader == "" {
			log.Impl.Error("Отсутствует токен авторизации")
			c.AbortWithStatusJSON(403, gin.H{"error": "Отсутствует токен авторизации"})
			return
		}

		//claims := jwt.MapClaims{}  // Если хотим просто выгрузить в map и перебрать открытые данные
		//claims := jwt.StandardClaims{}
		claims := &MyCustomClaimsLight{}
		token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.SecretKeyJWT), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			log.Impl.Error("Неверно сформированный токен аутентификации")
			c.AbortWithStatusJSON(403, gin.H{"error": "Неверно сформированный токен аутентификации"})
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			log.Impl.Error("Токен недействителен")
			c.AbortWithStatusJSON(403, gin.H{"error": "Токен недействителен"})
			return
		}

		//str := fmt.Sprint(claims.UserId) //Useful for monitoring
		//fmt.Println(str)
		// ctx := context.WithValue(r.Context(), "user", tk.UserId)

		c.Next()
	}
}

//a struct to rep user account
type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	//QueryHistory []*QueryHistory `gorm:"many2many:account_query_history;"`
}

type TokenSession struct {
	//Status    bool   `json:"status"`
	ID        string `json:"id"`
	ExpiresIn int64  `json:"expiresIn"`
	ExpiresAt int64  `json:"expiresAt"`
}

type MyCustomClaims struct {
	UserId uint `json:"userid"`
	jwt.StandardClaims
}

func MiddleWareCheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Прописать проверку наличия токена в редис.

		// requestPath := c.Request.URL.Path          //current request path

		tokenHeader := c.GetHeader("TokenBearer")

		// При большем обращении нужны разные клиенты для получения токена.
		dataRedis, err := shareRedis.GetLibraryRediGo(shareRedis.PoolRedisRediGolibrary, tokenHeader, 12)

		if dataRedis == "" {
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Токен просрочен"}}
			c.AbortWithStatusJSON(401, AnswerWebV1)
			return
		}
		if err != nil {
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, err.Error()}}
			c.AbortWithStatusJSON(401, AnswerWebV1)
			return
		}

		if tokenHeader == "" {
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Отсутствует токен авторизации"}}
			c.AbortWithStatusJSON(401, AnswerWebV1)
			return
		}

		//ParseTest(tokenHeader)

		//claims := jwt.MapClaims{} // Если хотим просто выгрузить в map и перебрать открытые данные
		//claims := jwt.StandardClaims{}
		claims := &MyCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.SecretKeyJWT), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Неверно сформированный токен аутентификации"}}
			c.AbortWithStatusJSON(401, AnswerWebV1)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Неверно сформированный токен аутентификации"}}
			c.AbortWithStatusJSON(401, AnswerWebV1)
			return
		}

		// c.Set("user-id", claims["userid"])
		// log.Impl.Errorf("Запрос к ЛК от пользователя: %s Путь: %s\n ", claims["userid"], c.FullPath())

		c.Set("user-data", dataRedis)

		c.Next()
	}
}

type LoginAnswer struct {
	JWTtoken    string
	ExpiresAt   int64
	DurationSec int64
	User        web.LkUsers
}

func Login(login, password string) (LoginAnswer, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, login)

	queryText := `select
	id,
	coalesce(exp_sec, 0) as exp_sec,
	role,
	login,
	password,
	coalesce(blocked, false) as blocked,
	coalesce(user_id, '') as user_id,
	coalesce(full_name, '') as full_name,
	coalesce(email, '') as email,
	coalesce(insurance_number, '') as insurance_number,
	coalesce(notes, '') as notes,
	coalesce(status, '') as status,
	coalesce(source, '') as source,
	coalesce(person_json_byte, '{}'::jsonb) as person_json_byte,
	coalesce(additional_settings_user_json_byte, '{}'::jsonb) as additional_settings_user_json_byte
from
	public.lk_users
where
	login = $1;`

	var LkUsers web.LkUsers

	DB, err := shareStore.GetDB(config.Conf.DatabaseURL)
	if err != nil {
		return LoginAnswer{}, err
	}

	rows, err := DB.Query(queryText, argsquery...)
	if err != nil {
		fmt.Println("Auth err --- ", err)
		return LoginAnswer{}, err
	}

	var blocked bool
	var V1ActiveWorkers web.V1ActiveWorkers
	var AdditionalSettingsUser web.AdditionalSettingsUser

	for rows.Next() {
		err = rows.Scan(&LkUsers.ID, &LkUsers.ExpSec, &LkUsers.Role, &LkUsers.Login, &LkUsers.Password, &blocked, &LkUsers.UserID, &LkUsers.FullName, &LkUsers.Email,
			&LkUsers.InsuranceNumber, &LkUsers.Notes, &LkUsers.Status, &LkUsers.Source, &V1ActiveWorkers, &AdditionalSettingsUser)
		if err != nil {
			fmt.Println("Auth err --- ", err)
			return LoginAnswer{}, err
		}
	}

	LkUsers.Person = V1ActiveWorkers
	LkUsers.AdditionalSettingsUser = AdditionalSettingsUser

	if blocked {
		err := fmt.Errorf("Учетная запись заблокированна")
		fmt.Println("Auth err --- ", err)
		return LoginAnswer{}, err
	}

	defer rows.Close()

	// err = bcrypt.CompareHashAndPassword([]byte(PasswordDB ), []byte(password))
	// if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
	// 	return "", fmt.Errorf("Неверные логин или пароль. Пожалуйста, попробуйте еще раз")
	// }

	usernameHash := sha256.Sum256([]byte(strings.ToLower(login)))
	passwordHash := sha256.Sum256([]byte(strings.ToLower(password)))
	expectedUsernameHash := sha256.Sum256([]byte(strings.ToLower(LkUsers.Login)))
	expectedPasswordHash := sha256.Sum256([]byte(strings.ToLower(LkUsers.Password)))

	usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
	passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

	if usernameMatch && passwordMatch {
		// 	next.ServeHTTP(w, r)
		// 	return
	} else {
		return LoginAnswer{}, fmt.Errorf("Неверные логин или пароль. Пожалуйста, попробуйте еще раз")
	}

	var Duration time.Duration // ExpiresIn

	if LkUsers.ExpSec == 0 {
		Duration = time.Hour * 672 // 30 Дней
	} else {
		Duration = time.Second * time.Duration(LkUsers.ExpSec)
	}

	uuid := uuid.New()
	//uuid--- := strings.Replace(uuid.String(), "-", "", -1)

	//Create JWT token
	ExpiresAt := time.Now().Add(Duration).Unix() // 186 - 7 days
	//ExpiresAt := time.Now().UTC().Add(Duration).Unix() // 186 - 7 days
	claims := jwt.StandardClaims{
		ExpiresAt: ExpiresAt,
		Issuer:    "auth.keeper",
		Id:        uuid.String(),
	}
	tk := &MyCustomClaims{uint(LkUsers.ID), claims}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.Conf.SecretKeyJWT))

	// При большем обращении нужны разные клиенты для получения токена.
	err = shareRedis.SelectLibraryRediGo(shareRedis.PoolRedisRediGolibrary, 12)
	if err != nil {
		fmt.Println("Auth err --- ", err)
		return LoginAnswer{}, err
	}

	DurationSec := int64(Duration.Seconds())

	LkUsers.Password = ""
	LoginAnswerReturn := LoginAnswer{JWTtoken: tokenString, ExpiresAt: ExpiresAt, DurationSec: DurationSec, User: LkUsers}
	byteData, err := json.Marshal(LoginAnswerReturn)
	if err != nil {
		fmt.Println("Auth err --- ", err)
		return LoginAnswer{}, err
	}

	err = shareRedis.SetLibraryRediGo(shareRedis.PoolRedisRediGolibrary, tokenString, byteData, 12, DurationSec)
	if err != nil {
		fmt.Println("Auth err --- ", err)
		return LoginAnswer{}, err
	}

	return LoginAnswerReturn, nil
}

func ValidateSession(tokenHeader string) (time.Duration, error) {

	// При большем обращении нужны разные клиенты для получения токена.
	err := shareRedis.SelectLibraryGoRedis(shareRedis.RedisClientGoRedisLibrary, 12)
	if err != nil {
		return 0, err
	}

	ttl, err := shareRedis.RedisClientGoRedisLibrary.TTL(context.Background(), tokenHeader).Result()
	if err != nil {
		return 0, err
	}
	// if ttl < 0 {
	// 	if -1 == ttl.Seconds() {
	// 		fmt.Print("The key will not expire.\n")
	// 	} else if -2 == ttl.Seconds() {
	// 		fmt.Print("The key does not exist.\n")
	// 	} else {
	// 		fmt.Printf("Unexpected error %d.\n", ttl.Seconds())
	// 	}
	// }

	return ttl, nil

}

func RemoveSession(tokenHeader string) error {

	// При большем обращении нужны разные клиенты для получения токена.
	err := shareRedis.SelectLibraryGoRedis(shareRedis.RedisClientGoRedisLibrary, 12)
	if err != nil {
		return err
	}

	_, err = shareRedis.RedisClientGoRedisLibrary.Del(context.Background(), tokenHeader).Result()
	if err != nil {
		return err
	}

	return nil

}

func LoginHandlersV1(c *gin.Context) {

	account := &Account{}
	err := json.NewDecoder(c.Request.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		c.JSON(http.StatusBadRequest, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	LoginAnswer, err := Login(account.Login, account.Password)
	if err != nil {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		c.JSON(http.StatusBadRequest, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	TokenSession := TokenSession{ID: LoginAnswer.JWTtoken, ExpiresAt: LoginAnswer.ExpiresAt, ExpiresIn: LoginAnswer.DurationSec}
	// byteData, err := json.Marshal(TokenSession)
	// if err != nil {
	// 	AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	// 	c.JSON(http.StatusBadRequest, AnswerWebV1)
	// 	log.Impl.Error(err.Error())
	// 	return
	// }

	var AnswerWebV1 web.AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = TokenSession
	AnswerWebV1.Error = nil

	//c.Data(http.StatusOK, "application/json", byteData)
	c.JSON(http.StatusOK, AnswerWebV1)
}

func LoginBasicHandlersV1(c *gin.Context) {

	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	_ = hasAuth

	LoginAnswer, err := Login(user, password)
	if err != nil {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		c.JSON(http.StatusBadRequest, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	TokenSession := TokenSession{ID: LoginAnswer.JWTtoken, ExpiresAt: LoginAnswer.ExpiresAt, ExpiresIn: LoginAnswer.DurationSec}
	// byteData, err := json.Marshal(TokenSession)
	// if err != nil {
	// 	c.Abort()
	// 	c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	// 	AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	// 	c.JSON(http.StatusBadRequest, AnswerWebV1)
	// 	log.Impl.Error(err.Error())
	// 	return
	// }

	var AnswerWebV1 web.AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = TokenSession
	AnswerWebV1.Error = nil

	//c.Data(http.StatusOK, "application/json", byteData)
	c.JSON(http.StatusOK, AnswerWebV1)
}

func RemoveSessionHandlersV1(c *gin.Context) {

	tokenHeader := c.GetHeader("TokenBearer")
	err := fmt.Errorf("Отправлен пустой токен")
	if tokenHeader == "" {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		c.JSON(http.StatusBadRequest, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	err = RemoveSession(tokenHeader)
	if err != nil {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		c.JSON(http.StatusBadRequest, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	var AnswerWebV1 web.AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = "Идентификатор вашей сессии удален"
	AnswerWebV1.Error = nil

	c.JSON(http.StatusOK, AnswerWebV1)
}

func ValidateSessionHandlersV1(c *gin.Context) {

	tokenHeader := c.GetHeader("TokenBearer")
	if tokenHeader == "" {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Отсутствует токен авторизации"}}
		c.JSON(http.StatusUnauthorized, AnswerWebV1)
		log.Impl.Error("Отсутствует токен авторизации")
		return
	}

	DurationExpired, err := ValidateSession(tokenHeader)
	if err != nil {
		AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, err.Error()}}
		c.JSON(http.StatusUnauthorized, AnswerWebV1)
		log.Impl.Error(err.Error())
		return
	}

	if DurationExpired < 0 {
		if -1 == DurationExpired {
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "The key will not expire"}}
			c.JSON(http.StatusUnauthorized, AnswerWebV1)
			log.Impl.Error("The key will not expire")
			return
		} else if -2 == DurationExpired {
			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "The key does not exist"}}
			c.JSON(http.StatusUnauthorized, AnswerWebV1)
			log.Impl.Error("The key does not exist")
			return
		} else {

			AnswerWebV1 := web.AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusUnauthorized, "Unexpected error"}}
			c.JSON(http.StatusUnauthorized, AnswerWebV1)
			log.Impl.Error("Unexpected error")
			//c.JSON(http.StatusOK, gin.H{"Unexpected error": DurationExpired.Seconds()})
			return
		}
	}

	DataAnswer := struct {
		//Status    bool    `json:"status"`
		ExpiresIn float64 `json:"expiresIn"`
		ExpiresAt int64
	}{
		//Status:    true,
		ExpiresIn: DurationExpired.Seconds(),
		ExpiresAt: time.Now().Add(DurationExpired).Unix(),
	}

	var AnswerWebV1 web.AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = DataAnswer
	AnswerWebV1.Error = nil

	//c.Data(http.StatusOK, "application/json", byteData)
	c.JSON(http.StatusOK, AnswerWebV1)

	//c.JSON(http.StatusOK, DataAnswer)
}
