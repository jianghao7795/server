package utils

import (
	"crypto/rsa"
	global "server/model"
	"server/model/system/request"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJWT() *JWT {
	return &JWT{
		PrivateKey: global.RunCONFIG.JWT.PrivateKey,
		PublicKey:  global.RunCONFIG.JWT.PublicKey,
	}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: global.CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		// StandardClaims: jwt.StandardClaims{
		// 	NotBefore: time.Now().Unix(),                                 // 签名生效时间
		// 	ExpiresAt: time.Now().Unix() + global.CONFIG.JWT.ExpiresTime, // 过期时间 7天  配置文件
		// 	Issuer:    global.CONFIG.JWT.Issuer,                          // 签名的发行者
		// },
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.CONFIG.JWT.ExpiresTime) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                               // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                               // 生效时间
			Issuer:    global.CONFIG.JWT.Issuer,                                                                     // 签名的发行者
		},
	}
	return claims
}

// 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims) // 使用RS256算法
	return token.SignedString(j.PrivateKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
// func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
// 	v, err, _ := global.ConcurrencyControl.Do("JWT:"+oldToken, func() (any, error) {
// 		return j.CreateToken(claims)
// 	})
// 	return v.(string), err
// }

// 解析 token 校验token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i any, e error) {
		return j.PublicKey, nil
	})
	if err != nil {
		return nil, ReportError(err)
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ReportError(err)

	} else {
		return nil, ReportError(err)
	}
}
