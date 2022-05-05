package dingzhen

import (
	"errors"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	api     = "https://www.yiyandingzhen.top/livesearch.php"
	url     = "https://www.yiyandingzhen.top/"
	imgPath = "data/dingzhen"
)

func init() { // 插件主体
	os.MkdirAll(imgPath, 0777)
	engine := control.Register("dingzhen", &control.Options{
		DisableOnDefault: false,
		Help:             "- 一眼丁真\n",
		PublicDataFolder: "Dingzhen",
	})

	engine.OnFullMatch(`一眼丁真`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			imgSrc, err := GetImgSrc()
			if err != nil || imgSrc == "" {
				ctx.SendChain(message.Text("鉴定为搜索源寄了"))
				return
			}
			pathName, err := GetImg(imgSrc)
			if err != nil {
				ctx.SendChain(message.Text("鉴定为下载寄了"))
				return
			}
			pathName = "file://" + pathName
			ctx.SendChain(message.Image(pathName))
		})
}

func GetImgSrc() (string, error) {
	respBody, err := web.GetData(api)
	if err != nil {
		return "", err
	}
	imgRegex, _ := regexp.Compile(`"pic\\/(\w*?).(jpeg|jpg|png)"`)
	imgSrcs := imgRegex.FindAll(respBody, -1)
	if len(imgSrcs) == 0 {
		return "", errors.New("imgSrcs nil ")
	}
	imgSrc := string(imgSrcs[rand.Intn(len(imgSrcs))])
	imgSrc = strings.ReplaceAll(imgSrc, "\\", "")
	imgSrc = strings.ReplaceAll(imgSrc, "\"", "")
	return imgSrc, nil
}

func GetImg(imgSrc string) (string, error) {

	data, err := web.GetData(url + imgSrc)
	if err != nil {
		return "", err
	}
	imgSrcBase := path.Base(imgSrc)
	err = os.WriteFile(path.Join(imgPath, imgSrcBase), data, 0777)
	if err != nil {
		return "", err
	}
	retPath, _ := filepath.Abs(path.Join(imgPath, imgSrcBase))
	return retPath, nil
}
