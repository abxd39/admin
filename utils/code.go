package utils

import (
	"fmt"

	google "github.com/mojocn/base64Captcha"
)

func demoCodeCaptchaCreate() {
	//config struct for digits
	//数字验证码配置
	var configD = google.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	//config struct for audio
	//声音验证码配置
	// var configA = google.ConfigAudio{
	// 	CaptchaLen: 6,
	// 	Language:   "zh",
	// }
	//config struct for Character
	//字符,公式,验证码配置
	// var configC = google.ConfigCharacter{
	// 	Height: 60,
	// 	Width:  240,
	// 	//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
	// 	Mode:               google.CaptchaModeNumber,
	// 	ComplexOfNoiseText: google.CaptchaComplexLower,
	// 	ComplexOfNoiseDot:  google.CaptchaComplexLower,
	// 	IsShowHollowLine:   false,
	// 	IsShowNoiseDot:     false,
	// 	IsShowNoiseText:    false,
	// 	IsShowSlimeLine:    false,
	// 	IsShowSineLine:     false,
	// 	CaptchaLen:         6,
	// }
	//创建声音验证码
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	//idKeyA, capA := google.GenerateCaptcha("", configA)
	//以base64编码
	//base64stringA := google.CaptchaWriteToBase64Encoding(capA)
	//创建字符公式验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	//idKeyC, capC := google.GenerateCaptcha("", configC)
	//以base64编码
	//base64stringC := google.CaptchaWriteToBase64Encoding(capC)
	//创建数字验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyD, capD := google.GenerateCaptcha("", configD)
	//以base64编码
	base64string := google.CaptchaWriteToBase64Encoding(capD)

	// fmt.Println(idKeyA, base64stringA, "\n")
	// fmt.Println(idKeyC, base64stringC, "\n")
	fmt.Println(idKeyD, base64string)

}
