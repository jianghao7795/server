package system

import (
	"server-fiber/config"
	global "server-fiber/model"
	"server-fiber/model/system"
	"server-fiber/utils"

	"go.uber.org/zap"
)

//@author: wuhao
//@function: GetSystemConfig
//@description: 读取配置文件
//@return: err error, conf config.Server

func (systemConfigService *SystemConfigService) GetSystemConfig() (conf config.Server, err error) {
	return global.CONFIG, nil
}

// @description   set system config,
//@author: wuhao
//@function: SetSystemConfig
//@description: 设置配置文件
//@param: system model.System
//@return: err error

func (systemConfigService *SystemConfigService) SetSystemConfig(system system.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.VIPER.Set(k, v)
	}
	err = global.VIPER.WriteConfig()
	return err
}

//@author: wuhao
//@function: GetServerInfo
//@description: 获取服务器信息
//@return: server *utils.Server, err error

func (systemConfigService *SystemConfigService) GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil {
		global.LOG.Error("func utils.InitCPU() Failed", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Ram, err = utils.InitRAM(); err != nil {
		global.LOG.Error("func utils.InitRAM() Failed", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		global.LOG.Error("func utils.InitDisk() Failed", zap.String("err", err.Error()))
		return &s, err
	}

	return &s, nil
}
