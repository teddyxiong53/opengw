package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// 定义一个jwt对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

// 初始化jwt对象
func NewJWT() *JWT {
	return &JWT{
		[]byte("bgbiao.top"),
	}
}

// 自定义有效载荷(这里采用自定义的Name和Email作为有效载荷的一部分)
type CustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// 构造用户表
type User struct {
	Id          int32  `gorm:"AUTO_INCREMENT"`
	Name        string `json:"username"`
	Pwd         string `json:"password"`
	Phone       int64  `gorm:"DEFAULT:0"`
	Email       string `gorm:"type:varchar(20);unique_index;"`
	CreatedAt   *time.Time
	UpdateTAt   *time.Time
}

// LoginReq请求参数
type LoginReq struct {
	Name string `json:"username"`
	Pwd  string `json:"password"`
}

// 登陆结果
type LoginResult struct {
	Token string `json:"token"`
	// 用户模型
	Name string `json:"name"`
	//model.User
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "bgbiao.top" // 签名信息应该设置成动态从库中获取

	loginResult LoginResult
)

// 调用jwt-go库生成token
// 指定编码的算法为jwt.SigningMethodHS256
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// 返回一个token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}


// token解码
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不可用")
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token过期")
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("无效的token")
			} else {
				return nil, fmt.Errorf("token不可用")
			}

		}
	}

	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token无效")

}

// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"Code"		: "-1",
				"Message"	: "请求未携带token，无权限访问",
				"Data"		: "",
			})
			c.Abort()
			return
		}

		//log.Print("get token: ", token)

		// 初始化一个JWT对象实例，并根据结构体方法来解析token
		j := NewJWT()
		// 解析token中包含的相关信息(有效载荷)
		claims, err := j.ParserToken(token)

		if err != nil {
			// token过期
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"Code"		: "-1",
					"Message"	: "token授权已过期，请重新申请授权",
					"Data"		: "",
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"Code"		: "-1",
				"Message"	: err.Error(),
				"Data"		: "",
			})
			c.Abort()
			return
		}

		// 将解析后的有效载荷claims重新写入gin.Context引用对象中
		c.Set("claims", claims)

	}
}

// LoginCheck验证
func LoginCheck(login LoginReq) (bool, User, error) {
	userData := User{}
	userExist := false

	//var user = User{
	//	Name:"admin",
	//	Pwd:"admin",
	//}

	if login.Name == "admin" && login.Pwd == "admin" {
		userExist = true
		userData.Name = "admin"
		userData.Email = "admin"
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("登陆信息有误")
	}
	return userExist, userData, nil
}

// token生成器
// md 为上面定义好的middleware中间件
func generateToken(c *gin.Context, user User) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := NewJWT()

	// 构造用户claims信息(负荷)
	claims := CustomClaims{
		user.Name,
		user.Email,
		jwt.StandardClaims{
			//NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			NotBefore: int64(time.Now().Unix()), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 签名过期时间
			Issuer:    "bgbiao.top",                    // 签名颁发者
		},
	}

	// 根据claims生成token对象
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Code"		: "-1",
			"Message"	: err.Error(),
			"Data"		: "",
		})
	}

	log.Println(token)
	// 封装一个响应数据,返回用户名和token
	data := LoginResult{
		Name:  user.Name,
		Token: token,
	}

	loginResult = data

	c.JSON(http.StatusOK, gin.H{
		"Data"		:  data,
		"Message"	: "login sucess",
		"Code"		: "0",
	})
	return

}
