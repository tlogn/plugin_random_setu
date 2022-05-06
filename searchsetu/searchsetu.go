package searchsetu

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/web"
	"github.com/tidwall/gjson"
	"github.com/tlogn/plugin_random_setu/utils"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	api     = "https://api.lolicon.app/setu/v2"
	imgPath = "data/search_setu/"
)

func init() {
	os.MkdirAll(imgPath, 0777)
	compilerNum, _ := regexp.Compile(`[0-9]+`)

	engine := control.Register("setutime", &control.Options{
		DisableOnDefault: false,
		Help:             "- (R18)来点[xxx]\n",
		PublicDataFolder: "SearchSetuTime",
	})
	engine.ApplySingle(ctxext.DefaultSingle).OnRegex("^(R18)?(来点)(.+)$").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.State["regex_matched"].([]string)
			r18 := false
			if msg[1] == "R18" {
				r18 = true
			}
			keyword := msg[len(msg)-1]

			url := api + "?keyword=" + url2.QueryEscape(keyword)
			if r18 {
				url += "&r18=1"
			} else {
				url += "&r18=0"
			}

			data, err := web.GetData(url)
			if err != nil {
				ctx.SendChain(message.Text("你的xp太奇怪了！"))
				return
			}
			json := gjson.ParseBytes(data)
			url = json.Get("data.0.urls.original").Str
			if url == "" {
				ctx.SendChain(message.Text("呜呜呜，搜索失败了啦"))
				return
			}

			url = path.Base(url)
			pixivIdString := compilerNum.FindString(url)
			pixivId, _ := strconv.ParseInt(pixivIdString, 10, 64)
			imgName, err := utils.DownloadImageFromPixiv(pixivId, imgPath)
			if err != nil {
				ctx.SendChain(message.Text("呜呜呜，下载失败了啦"))
				return
			}
			pathName, _ := filepath.Abs(path.Join(imgPath, imgName))
			pathName = "file://" + pathName
			setu := ctx.SendChain(message.Image(pathName))

			time.Sleep(10 * time.Second)
			ctx.DeleteMessage(setu);
		})
}

/*
func downloadImgFromLolicon(keyword string, r18 bool) ([]string, error) {
	pool, err := loliconApiPool.NewLoliconPool(&loliconApiPool.Config{
		ApiKey:   keyword, // use your api key here
		CacheMin: 0,
		CacheMax: 1,
		Persist:  loliconApiPool.NewNilPersist(),
	})
	if err != nil {
		return nil, errors.New("pool init 寄了")
	}

	poolR18opt := loliconApiPool.R18Off
	if r18 {
		poolR18opt = loliconApiPool.R18On
	}
	img, err := pool.Get(
		loliconApiPool.NumOption(1),
		loliconApiPool.R18Option(poolR18opt),
	)
	if err != nil {
		return nil, errors.New("下载图寄了")
	}
	ret := make([]string, 0)
	for _, i := range img {
		b, err := i.Content()
		if err != nil {
			panic(err)
		}
		if len(b) == 0 {
			continue
		}
		_, cfg, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil {
			logrus.Errorf("unknown format")
			continue
		}
		filename := fmt.Sprintf("%v.%v", i.Pid, cfg)
		ioutil.WriteFile(filename, b, 0777)
		ret = append(ret, filename)
	}
	return ret, nil
} */
