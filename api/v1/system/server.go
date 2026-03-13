package system

import systemServer "server/service/system"

// service 端的服务

var apiService = systemServer.ApiServiceApp
var authorityBtnService = systemServer.AuthorityBtnServiceApp
var authorityService = systemServer.AuthorityServiceApp
var menuService = systemServer.MenuServiceApp
var casbinService = systemServer.CasbinServiceApp
var autoCodeHistoryService = systemServer.AutoCodeHistoryServiceApp
var autoCodeService = systemServer.AutoCodeServiceApp
var dictionaryDetailService = systemServer.DictionaryDetailServiceApp
var dictionaryService = systemServer.DictionaryServiceApp
var githubService = systemServer.GithubServiceApp
var initDBService = systemServer.InitDBServiceApp
var jwtService = systemServer.JwtServiceApp
var baseMenuService = systemServer.BaseMenuServiceApp
var operationRecordService = systemServer.OperationRecordServiceApp
var systemConfigService = systemServer.SystemConfigServiceApp
var userProblem = systemServer.ProblemApp
var userService = systemServer.UserServiceApp
