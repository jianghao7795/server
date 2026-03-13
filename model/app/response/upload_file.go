/*
 * @Author: jianghao
 * @Date: 2022-10-17 16:58:13
 * @LastEditors: jianghao
 * @LastEditTime: 2022-10-17 17:01:29
 */

package response

import "server/model/app"

type ResponseUploadFile struct {
	File app.FileUploadAndDownload `json:"file" form:"file"`
}
