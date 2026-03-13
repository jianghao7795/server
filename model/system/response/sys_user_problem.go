package response

import global "server/model"

type SysUserProblemSetting struct {
	global.MODEL
	SysUserId uint   `json:"sys_user_id"`
	Problem   string `json:"problem"`
	// Answer    string `json:"answer"`
}

type SysUserProblemSettingResponse struct {
	ProblemsSetting []SysUserProblemSetting `json:"problems_setting"`
}

type SysUserProblem struct {
	SysUserId uint   `json:"sys_user_id"`
	Problem   string `json:"problem"`
}

type SysUserProblemResponse struct {
	Problems []SysUserProblem `json:"problems"`
}
