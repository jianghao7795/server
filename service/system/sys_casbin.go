package system

import (
	"errors"
	"server/model/system/request"
	"sync"

	global "server/model"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "gorm.io/driver/mysql"
	"go.uber.org/zap"
)

//
//@function: UpdateCasbin
//@description: 更新casbin权限
//@param: authorityId string, casbinInfos []request.CasbinInfo
//@return: error

func (casbinService *CasbinService) UpdateCasbin(authorityId string, casbinInfos []request.CasbinInfo) error {
	casbinService.ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, v := range casbinInfos {
		rules = append(rules, []string{authorityId, v.Path, v.Method})
	}
	e := casbinService.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

//
//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]any{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

//
//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []request.CasbinInfo

func (casbinService *CasbinService) GetPolicyPathByAuthorityId(authorityId string) (pathMaps []request.CasbinInfo) {
	e := casbinService.Casbin()
	list, _ := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//
//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool

func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := casbinService.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//
//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (casbinService *CasbinService) Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				// gorm-adapter 在策略为空或 ptype 为空时会 panic: slice bounds out of range
				// 回退为仅模型、无 DB 策略的 Enforcer，避免服务崩溃
				global.LOG.Warn("casbin 从数据库加载策略时发生 panic，使用空策略 Enforcer。请检查 casbin_rule 表：ptype 字段不能为空，且无异常数据", zap.Any("panic", r))
				var err error
				syncedEnforcer, err = casbin.NewSyncedEnforcer(global.CONFIG.Casbin.ModelPath)
				if err != nil {
					global.LOG.Error("casbin 回退 Enforcer 创建失败", zap.Error(err))
				}
			}
		}()
		a, _ := gormadapter.NewAdapterByDB(global.DB)
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(global.CONFIG.Casbin.ModelPath, a)
	})
	if syncedEnforcer != nil {
		_ = syncedEnforcer.LoadPolicy()
	}
	return syncedEnforcer
}
