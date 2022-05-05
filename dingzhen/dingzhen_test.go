package dingzhen

import (
	"fmt"
	"github.com/FloatTech/zbputils/web"
	"regexp"
	"strings"
	"testing"
)

func TestApiGetData(t *testing.T) {

	respBody, _ := web.GetData(api)
	imgRegex, _ := regexp.Compile(`"pic\\/(\w*?).(jpeg|jpg|png)"`)
	data := imgRegex.FindAll(respBody, -1)
	for _, img := range data {
		imgSrc := string(img)
		imgSrc = strings.ReplaceAll(imgSrc, "\\", "")
	}

	//fmt.Println(render.Render(dingzhenRegex.FindString(body)))
}

func TestGetImgSrc(t *testing.T) {
	fmt.Println(GetImgSrc())
}

func TestGetImg(t *testing.T) {
	imgSrc, _ := GetImgSrc()
	fmt.Println(GetImg(imgSrc))
}
