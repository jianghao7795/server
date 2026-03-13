package system

import (
	global "server/model"
	"server/model/common/response"
	systemRes "server/model/system/response"

	"github.com/gofiber/fiber/v3"
	"github.com/mojocn/base64Captcha"
	"github.com/wenlng/go-captcha/captcha"
	"go.uber.org/zap"
	"golang.org/x/image/font"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
// var store = captcha.NewDefaultRedisStore()
var store = base64Captcha.DefaultMemStore

type BaseApi struct{}

// Captcha
// @Tags Base
// @Summary 生成验证码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=systemRes.SysCaptchaResponse,msg=string} "生成验证码,返回包括随机数id,base64,验证码长度"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /base/captcha [get]
func (b *BaseApi) Captcha(c fiber.Ctx) error {
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.CONFIG.Captcha.ImgHeight, global.CONFIG.Captcha.ImgWidth, global.CONFIG.Captcha.KeyLong, global.CONFIG.Captcha.MaxSkew, global.CONFIG.Captcha.DotCount)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	newCaptcha := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, _, err := newCaptcha.Generate(); err != nil {
		global.LOG.Error("验证码获取失败!", zap.Error(err))
		return response.FailWithMessage("验证码获取失败", c)
	} else {
		return response.OkWithDetailed(systemRes.SysCaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: global.CONFIG.Captcha.KeyLong,
		}, "验证码获取成功", c)
	}
}

// CaptchaImg
// @Tags Base
// @Summary 生成图片验证码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=systemRes.SysCaptchaImgResponse,msg=string} "生成验证码,返回包括随机数id,base64,验证码长度"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /base/captcha/img [get]
func (b *BaseApi) CaptchaImg(c fiber.Ctx) error {
	capt := captcha.GetCaptcha()
	dots, imageBase64, thumbImageBase64, key, err := capt.Generate()
	// 设置图片属性
	// ====================================================
	// Method: SetImageSize(size Size);
	// Desc: 设置验证码主图的尺寸
	// ====================================================
	capt.SetImageSize(captcha.Size{Width: 300, Height: 300})

	// ====================================================
	// Method: SetImageQuality(val int);
	// Desc: 设置验证码主图清晰度，压缩级别范围 QualityCompressLevel1 - 5，QualityCompressNone无压缩，默认为最低压缩级别
	// ====================================================
	capt.SetImageQuality(captcha.QualityCompressNone)

	// ====================================================
	// Method: SetFontHinting(val font.Hinting);
	// Desc: 设置字体Hinting值 (HintingNone,HintingVertical,HintingFull)
	// ====================================================
	capt.SetFontHinting(font.HintingFull)

	// ====================================================
	// Method: SetTextRangLen(val captcha.RangeVal);
	// Desc: 设置验证码文本显示的总数随机范围
	// ====================================================
	capt.SetTextRangLen(captcha.RangeVal{Min: 6, Max: 7})

	// ====================================================
	// Method: SetRangFontSize(val captcha.RangeVal);
	// Desc: 设置验证码文本的随机大小
	// ====================================================
	capt.SetRangFontSize(captcha.RangeVal{Min: 32, Max: 42})

	// ====================================================
	// Method: SetTextRangFontColors(colors []string);
	// Desc: 设置验证码文本的随机十六进制颜色
	// ====================================================
	capt.SetTextRangFontColors([]string{
		"#1d3f84",
		"#3a6a1e",
	})

	// ====================================================
	// Method: SetImageFontAlpha(val float64);
	// Desc:设置验证码字体的透明度
	// ====================================================
	capt.SetImageFontAlpha(0.5)

	// ====================================================
	// Method: SetTextShadow(val bool);
	// Desc:设置字体阴影
	// ====================================================
	capt.SetTextShadow(true)

	// ====================================================
	// Method: SetTextShadowColor(val string);
	// Desc:设置字体阴影颜色
	// ====================================================
	capt.SetTextShadowColor("#101010")

	// ====================================================
	// Method: SetTextShadowPoint(val captcha.Point);
	// Desc:设置字体阴影偏移位置
	// ====================================================
	capt.SetTextShadowPoint(captcha.Point{X: 1, Y: 1})

	// ====================================================
	// Method: SetTextRangAnglePos(pos []captcha.RangeVal);
	// Desc:设置验证码文本的旋转角度
	// ====================================================
	capt.SetTextRangAnglePos([]captcha.RangeVal{
		{Min: 1, Max: 15},
		{Min: 15, Max: 30},
		{Min: 30, Max: 45},
		{Min: 315, Max: 330},
		{Min: 330, Max: 345},
		{Min: 345, Max: 359},
	})

	// ====================================================
	// Method: SetImageFontDistort(val int);
	// Desc:设置验证码字体的扭曲程度
	// ====================================================
	capt.SetImageFontDistort(captcha.DistortLevel2)

	if err != nil {
		global.LOG.Error("验证码获取失败!", zap.Error(err))
		return response.FailWithMessage("验证码获取失败", c)
	} else {
		return response.OkWithDetailed(systemRes.SysCaptchaImgResponse{
			CaptchaKey: key,
			Dots:       dots,
			ThemeImage: thumbImageBase64,
			Images:     imageBase64,
		}, "验证码获取成功", c)
	}
}
