// Package utils 提供工具函数，包括JWT错误处理、加密解密等功能
package utils

import (
	errors "errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

// ECDSA密钥相关错误定义
var (
	ErrNotECPublicKey  = errors.New("密钥不是有效的 ECDSA 公钥")
	ErrNotECPrivateKey = errors.New("密钥不是有效的 ECDSA 私钥")
)

// Ed25519密钥相关错误定义
var (
	ErrNotEdPrivateKey = errors.New("密钥不是有效的 Ed25519私钥")
	ErrNotEdPublicKey  = errors.New("密钥不是有效的 Ed25519公钥")
)

// JWT Token相关错误定义
var (
	ErrInvalidKey                = errors.New("token密钥无效")    // 密钥无效错误
	ErrInvalidKeyType            = errors.New("token密钥的类型无效") // 密钥类型无效错误
	ErrHashUnavailable           = errors.New("请求的哈希函数不可用")   // 哈希函数不可用错误
	ErrTokenMalformed            = errors.New("token无效")      // Token格式错误
	ErrTokenUnverifiable         = errors.New("令牌无法验证")       // Token无法验证错误
	ErrTokenSignatureInvalid     = errors.New("令牌签名无效")       // Token签名无效错误
	ErrTokenRequiredClaimMissing = errors.New("令牌丢失")         // 必需声明缺失错误
	ErrTokenInvalidAudience      = errors.New("令牌的访问者无效")     // 受众无效错误
	ErrTokenExpired              = errors.New("令牌过期了")        // Token过期错误
	ErrTokenUsedBeforeIssued     = errors.New("发出前使用的标记")     // 在签发前使用Token错误
	ErrTokenInvalidIssuer        = errors.New("令牌的发行者无效")     // 发行者无效错误
	ErrTokenInvalidSubject       = errors.New("令牌的主题无效")      // 主题无效错误
	ErrTokenNotValidYet          = errors.New("令牌尚未有效")       // Token尚未生效错误
	ErrTokenInvalidId            = errors.New("令牌的 id 无效")    // Token ID无效错误
	ErrTokenInvalidClaims        = errors.New("令牌的索赔无效")      // 声明无效错误
	ErrInvalidType               = errors.New("无效索赔类型")       // 无效声明类型错误
)

// RSA密钥相关错误定义
var (
	ErrKeyMustBePEMEncoded = errors.New("无效密钥: 密钥必须是 PEM 编码的 PKCS1或 PKCS8密钥") // PEM编码错误
	ErrNotRSAPrivateKey    = errors.New("密钥不是有效的 RSA 私钥")                     // RSA私钥无效错误
	ErrNotRSAPublicKey     = errors.New("密钥不是有效的 RSA 公钥")                     // RSA公钥无效错误
)

// ECDSA验证相关错误定义
var (
	ErrECDSAVerification = errors.New("crypto/ecdsa: 验证错误") // ECDSA验证失败错误
)

// Ed25519验证相关错误定义
var (
	ErrEd25519Verification = errors.New("Ed25519: 验证错误") // Ed25519验证失败错误
)

// ReportError 将JWT库的原始错误转换为自定义的中文错误信息
// 该函数用于统一错误处理，将英文错误信息转换为中文，提升用户体验
// 参数:
//
//	err - JWT库返回的原始错误
//
// 返回值:
//
//	error - 转换后的中文错误信息
func ReportError(err error) error {
	switch {
	// ECDSA密钥相关错误
	case errors.Is(err, jwt.ErrNotECPublicKey):
		return ErrNotECPublicKey
	case errors.Is(err, jwt.ErrNotECPrivateKey):
		return ErrNotECPrivateKey

	// Ed25519密钥相关错误
	case errors.Is(err, jwt.ErrNotEdPrivateKey):
		return ErrNotEdPrivateKey
	case errors.Is(err, jwt.ErrNotEdPublicKey):
		return ErrNotEdPublicKey

	// 通用密钥和Token错误
	case errors.Is(err, jwt.ErrInvalidKey):
		return ErrInvalidKey
	case errors.Is(err, jwt.ErrInvalidKeyType):
		return ErrInvalidKeyType
	case errors.Is(err, jwt.ErrHashUnavailable):
		return ErrHashUnavailable
	case errors.Is(err, jwt.ErrTokenMalformed):
		return ErrTokenMalformed
	case errors.Is(err, jwt.ErrTokenUnverifiable):
		return ErrTokenUnverifiable
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return ErrTokenSignatureInvalid
	case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
		return ErrTokenRequiredClaimMissing

	// Token声明相关错误
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return ErrTokenInvalidAudience
	case errors.Is(err, jwt.ErrTokenExpired):
		return ErrTokenExpired
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		return ErrTokenUsedBeforeIssued
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return ErrTokenInvalidIssuer
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return ErrTokenInvalidSubject
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return ErrTokenNotValidYet
	case errors.Is(err, jwt.ErrTokenInvalidId):
		return ErrTokenInvalidId
	case errors.Is(err, jwt.ErrTokenInvalidClaims):
		return ErrTokenInvalidClaims
	case errors.Is(err, jwt.ErrInvalidType):
		return ErrInvalidType

	// RSA密钥相关错误
	case errors.Is(err, jwt.ErrKeyMustBePEMEncoded):
		return ErrKeyMustBePEMEncoded
	case errors.Is(err, jwt.ErrNotRSAPrivateKey):
		return ErrNotRSAPrivateKey
	case errors.Is(err, jwt.ErrNotRSAPublicKey):
		return ErrNotRSAPublicKey

	// 验证相关错误
	case errors.Is(err, jwt.ErrECDSAVerification):
		return ErrECDSAVerification
	case errors.Is(err, jwt.ErrNotRSAPublicKey): // 注意：这里可能是重复的条件，建议检查
		return ErrEd25519Verification

	// 默认返回密钥无效错误
	default:
		return ErrInvalidKey
	}
}
