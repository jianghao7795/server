package core

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	global "server/model"
	"server/utils"
	"time"

	// json "github.com/bytedance/sonic"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v3"      // fiber
	jwt "github.com/golang-jwt/jwt/v5" // jwt
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper" // viper配置文件读取
	"go.uber.org/zap"
)

var (
	publicKeyString = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD5aaI0TGfPWVrlX2m6l97j5gIf
IljZwDbOfJrycOPt+YQyPZFGbGoGYqC9MUFh9NSBHQySY4XrvWbpvuRH62S4IIyu
2uZ5HwfcbrwLe4vRYZEHNX6lMIvvsTwm+Iw96QhtnDRp5tcT+BCSJ2R/UPCy3sij
Uis1nNRTyiANM9xwlQIDAQAB
-----END PUBLIC KEY-----`)
	privateKeyString = []byte(`-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAPlpojRMZ89ZWuVf
abqX3uPmAh8iWNnANs58mvJw4+35hDI9kUZsagZioL0xQWH01IEdDJJjheu9Zum+
5EfrZLggjK7a5nkfB9xuvAt7i9FhkQc1fqUwi++xPCb4jD3pCG2cNGnm1xP4EJIn
ZH9Q8LLeyKNSKzWc1FPKIA0z3HCVAgMBAAECgYEAlfGTuq1Um6u0ocFC5xL2xTw6
ek1D7WVJjcoLvUVxQliHyUyKUjRLaZtR12wNSr7OIvE75aY/lQJBJLXGuBYkbf4O
d44b3fpnz0u3UWYnIz6EiZrnxeubtcBYwh0oLQYFJrHpPxrs27wYQD7Iou1U5xWC
sfyCL6qbtjnkbj9um+ECQQD+Sr1vKeBbdetRDqP+5313LNwTiVpBD9VaJ7ilpC8g
+CW6hFKmDYJ9CNu7Tq2Ky6mrI8h9fqG4KIB1i6AUegVrAkEA+xaA6DSx/D4JG44h
lzM1/ZK90q4gjNEUszKQLRo0N/Y1VLgkHe9rD7fX9JQD1/mNXtPZBEKq/jknIyE9
BCrh/wJBAI1Did7SPZ+xbysVXgdCJBHrasqzl97i49v0iwABqGjBUXw+/AbhJGc/
X4m7WG7tWvEGrN8CVOaQIWSrTx9w9ccCQGtnquipceFWoAq7d7731waLOMvgVGgo
H9SvsmC4rtGpbhK5Wim7+m4U1Dn4/tPqGgEQWeqqof5xcD5kv3cC6P0CQQCx6VXB
sqMBYlujFPwiaf6KzJY+dg8gNCZEZoytbd8TL7RolhcpQ//wTvsJY47PN5hPV6KD
nL9o4PchskjTFRVR
-----END PRIVATE KEY-----`)
)

// 读取配置 配置文件conf/config.yaml
func viperInit() (*viper.Viper, error) {
	var config string
	flag.StringVar(&config, "c", "./conf", "choose config file.")
	flag.Parse()
	if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
		configEnv := os.Getenv(utils.ConfigEnv)
		if configEnv == "" {
			config = utils.ConfigFile
			if isFile, err := utils.IsExistFile(config); isFile {
				fmt.Printf("您正在使用config的默认值,config的路径为%v\n", utils.ConfigFile)
			} else {
				panic("请检查配置文件" + config + "是否存在: " + err.Error())
			}
		} else {
			config = configEnv
			if isFile, err := utils.IsExistFile(config); isFile {
				fmt.Printf("您正在使用CONFIG环境变量,config的路径为%v\n", config)
			} else {
				panic("请检查配置文件" + config + "是否存在: " + err.Error())
			}
		}
	} else {
		if isFile, err := utils.IsExistFile(config); isFile {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config) // server -c config.yaml
		} else {
			panic("请检查配置文件" + config + "是否存在: " + err.Error())
		}
	}

	v := viper.New()
	// v.SetConfigFile(config)
	v.SetConfigName("config") // 指定文件名
	// v.SetConfigName("base")
	v.AddConfigPath(config) // 配置目录
	v.SetConfigType("yaml") // 配置文件类型

	// // 处理找不到配置文件的情况
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig() // 监听变化
	// config 变动之后的回调
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&global.CONFIG); err != nil {
			global.LOG.Error("config change error: ", zap.Error(err))
		}
	})
	if err := v.Unmarshal(&global.CONFIG); err != nil {
		return nil, err
	}

	// , err := os.ReadFile("./rsa_public_key.pem")
	// if err != nil {
	// 	return nil, err
	// }
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyString)
	// if err != nil {
	// 	return nil, err
	// }
	// , err := os.ReadFile("./private_key.pem")
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyString)
	if err != nil {
		return nil, err
	}
	// jwt
	global.RunCONFIG.JWT.PrivateKey = privateKey
	global.RunCONFIG.JWT.PublicKey = publicKey
	// root 适配性
	// 根据root位置去找到对应迁移位置,保证root路径有效
	global.CONFIG.AutoCode.Root, err = filepath.Abs("..") // filepath.Abs 是相对路径 变为绝对路径
	if err != nil {
		panic(err)
	}
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(global.CONFIG.JWT.ExpiresTime)),
	)
	{
		// global.RunCONFIG.FiberConfig.JSONEncoder = json.Marshal   // 自定义JSON编码器/解码器
		// global.RunCONFIG.FiberConfig.JSONDecoder = json.Unmarshal // 自定义JSON编码器/解码器
		global.RunCONFIG.FiberConfig.ErrorHandler = func(ctx fiber.Ctx, err error) error {
			// 状态代码默认为500
			code := fiber.StatusInternalServerError
			var message string = "服务器错误"
			// 如果是fiber.*Error，则检索自定义状态代码。
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			return ctx.Status(code).JSON(fiber.Map{"msg": message})
		}
	}
	return v, nil
}
