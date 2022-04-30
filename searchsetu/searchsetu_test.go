package searchsetu

import (
	"fmt"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	url2 "net/url"
	"path"
	"regexp"
	"testing"
)

func TestDownloadImgFromLolicon(t *testing.T) {
	fmt.Println(downloadImgFromLolicon("明日方舟", false))
}

func TestLoliconApi(t *testing.T) {
	data, err := web.GetData(api + "?tag=" + url2.QueryEscape(`萝莉`) + "&r18=1")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	json := gjson.ParseBytes(data)
	url := json.Get("data.0.urls.original").Str
	url = path.Base(url)
	fmt.Println(url)
	compilerNum, _ := regexp.Compile(`[0-9]+`)
	fmt.Println(compilerNum.FindString(url))

}
