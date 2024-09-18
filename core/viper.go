package core

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
	// json "github.com/bytedance/sonic"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"      // fiber
	jwt "github.com/golang-jwt/jwt/v5" // jwt
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper" // viper配置文件读取
	"go.uber.org/zap"

	"server-fiber/global"
	"server-fiber/utils"
)

// 读取配置 配置文件config.yaml
func viperInit(path ...string) (*viper.Viper, error) {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "config.yaml", "choose config file.")
		flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(utils.ConfigEnv); configEnv == "" {
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
				fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config) // server-fiber -c config.yaml
			} else {
				panic("请检查配置文件" + config + "是否存在: " + err.Error())
			}
		}
	} else {
		config = path[0]
		if isFile, err := utils.IsExistFile(config); isFile {
			fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config)
		} else {
			panic("请检查配置文件" + config + "是否存在: " + err.Error())
		}
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&global.CONFIG); err != nil {
			global.LOG.Error("config change error: ", zap.Error(err))
		}
	})
	if err := v.Unmarshal(&global.CONFIG); err != nil {
		return nil, err
	}

	publicKeyByte, err := os.ReadFile("./rsa_public_key.pem")
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return nil, err
	}
	privateKeyByte, err := os.ReadFile("./private_key.pem")
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
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
		global.RunCONFIG.FiberConfig.JSONEncoder = json.Marshal   // 自定义JSON编码器/解码器
		global.RunCONFIG.FiberConfig.JSONDecoder = json.Unmarshal // 自定义JSON编码器/解码器
		global.RunCONFIG.FiberConfig.ErrorHandler = func(ctx *fiber.Ctx, err error) error {
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
