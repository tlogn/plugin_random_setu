// Package randsetu 随机涩图
// 从setutime改
package randsetu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FloatTech/zbputils/web"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/FloatTech/AnimeAPI/pixiv"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type resultjson struct {
	Data struct {
		Illusts []struct {
			ID          int64  `json:"id"`
			Title       string `json:"title"`
			AltTitle    string `json:"altTitle"`
			Description string `json:"description"`
			Sanity      int    `json:"sanity"`
		} `json:"illusts"`
	} `json:"data"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ImgFIFO struct {
	queue []string
	mu    sync.Mutex
}

var (
	imgPath = "data/random_setu/"
	imgFIFO = ImgFIFO{
		queue: make([]string, 0),
		mu:    sync.Mutex{},
	}
)

const (
	updateInterval = time.Second * 20
	imgFIFOLimit   = 100
)

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// soutuapi 随机请求api
func randSetuApi() (r resultjson, err error) {
	str := RandStr(rand.Intn(5) + 1)
	data, err := web.GetData("https://copymanga.azurewebsites.net/api/pixivel?" + str + "?page=0")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &r)
	if err == nil && r.Error {
		err = errors.New(r.Message)
	}
	return
}

func downloadImageFromID(id int64) (string, error) {
	illust, err := pixiv.Works(id)
	if err != nil || len(illust.ImageUrls) == 0 || illust.ImageUrls[0] == "" {
		return "", err
	}
	/*
		for err != nil || illust.ImageUrls[0] == "" {
			illust, err = pixiv.Works(id)
		}*/
	u := illust.ImageUrls[0]
	n := path.Base(u)
	return n, illust.Download(0, path.Join(imgPath, n))
}

func randDownloadImage() (string, error) {
	resJson, err := randSetuApi()
	if err != nil {
		return "", err
	}
	if len(resJson.Data.Illusts) == 0 {
		return "", errors.New("resJson.Data is nil")
	}
	imgName, err := downloadImageFromID(resJson.Data.Illusts[0].ID)
	if err != nil {
		//return randDownloadImage()
		return "", err
	}
	return imgName, nil
}

func init() { // 插件主体
	os.MkdirAll(imgPath, 0777)
	imgFIFO.init()
	imgFIFO.run()
	engine := control.Register("randsetu", &control.Options{
		DisableOnDefault: false,
		Help: "- 随机涩图\n" +
			"- >randsetu status",
		PublicDataFolder: "Random_setu",
	})

	engine.OnFullMatch(`随机涩图`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			imgName := imgFIFO.get()
			pathName, _ := filepath.Abs(path.Join(imgPath, imgName))
			pathName = "file://" + pathName
			fmt.Println(pathName)
			if imgName == "" {
				ctx.SendChain(message.Text("别急别急，在下载了，哼哼哼啊啊啊啊啊啊啊啊"))
				return
			}
			ctx.SendChain(message.Image(pathName))
		})
}

func (q *ImgFIFO) init() {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 读取所有图片到fifo里
	dir, _ := ioutil.ReadDir(imgPath)
	for _, file := range dir {
		q.queue = append(q.queue, file.Name())
	}
}

func (q *ImgFIFO) run() {
	go func() {
		tick := time.Tick(updateInterval)
		for range tick {
			if len(q.queue) < imgFIFOLimit {
				q.insert()
			} else {
				q.pop()
				q.insert()
			}
		}
	}()
}

func (q *ImgFIFO) get() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return ""
	}

	return q.queue[rand.Intn(len(q.queue))]
}

func (q *ImgFIFO) insert() {
	// 下载的时候不用加锁，只有对queue进行操作再加锁
	imgName, err := randDownloadImage()
	for err != nil || imgName == "" {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, imgName)
}

func (q *ImgFIFO) pop() {
	if len(q.queue) == 0 {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	popName := q.queue[0]
	q.queue = q.queue[1:]
	os.Remove(path.Join(imgPath, popName))
}
