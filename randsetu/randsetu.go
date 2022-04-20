// Package randsetu 随机涩图
// 从setutime改
package randsetu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FloatTech/zbputils/web"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/FloatTech/AnimeAPI/pixiv"
	sql "github.com/FloatTech/sqlite"
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

// Pools 图片缓冲池
type imgpool struct {
	db     *sql.Sqlite
	dbmu   sync.RWMutex
	path   string
	max    int
	pool   map[string][]*message.MessageSegment
	poolmu sync.Mutex
}

func (p *imgpool) List() (l []string) {
	var err error
	p.dbmu.RLock()
	defer p.dbmu.RUnlock()
	l, err = p.db.ListTables()
	if err != nil {
		l = []string{"涩图"}
	}
	return l
}

var (
	imgPath = "data/random_setu/"
	dbName  = "RandSetu.db"
	pool    = &imgpool{
		db:   &sql.Sqlite{},
		path: imgPath,
		max:  100,
		pool: make(map[string][]*message.MessageSegment),
	}
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
	fmt.Println("https://copymanga.azurewebsites.net/api/pixivel?" + str + "?page=0")
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
	fmt.Println(id, illust, err)
	for err != nil {
		illust, err = pixiv.Works(id)
	}
	u := illust.ImageUrls[0]
	n := path.Base(u)
	return n, illust.Download(0, path.Join(imgPath, n))
}

func randDownloadImage() (string, error) {
	resJson, err := randSetuApi()
	if err != nil {
		return "", err
	}
	imgName, err := downloadImageFromID(resJson.Data.Illusts[0].ID)
	if err != nil {
		return randDownloadImage()
	}
	return imgName, nil
}

func init() { // 插件主体
	os.MkdirAll(imgPath, 777)
	engine := control.Register("randsetu", &control.Options{
		DisableOnDefault: false,
		Help: "- 随机涩图\n" +
			"- >randsetu status",
		PublicDataFolder: "RandSetu",
	})
	/*
		getdb := ctxext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
			// 如果数据库不存在则下载
			pool.db.DBPath = engine.DataFolder() + dbName
			_, _ = fileutil.GetLazyData(pool.db.DBPath, false, false)
			err := pool.db.Open()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return false
			}
			for _, imgtype := range pool.List() {
				if err := pool.db.Create(imgtype, &pixiv.Illust{}); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return false
				}
			}
			return true
		})
	*/
	engine.OnFullMatch(`随机涩图`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			imgName, err := randDownloadImage()
			pathName := path.Join(imgPath, imgName)
			if err != nil {
				ctx.SendChain(message.Text("寄！！！ " + imgName + " 自己搜吧"))
				return
			}
			ctx.SendChain(message.Image(pathName))
		})
	/*
		// 查询数据库涩图数量
		engine.OnFullMatch(">setu status", getdb).SetBlock(true).
			Handle(func(ctx *zero.Ctx) {
				state := []string{"[SetuTime]"}
				pool.dbmu.RLock()
				defer pool.dbmu.RUnlock()
				for _, imgtype := range pool.List() {
					num, err := pool.db.Count(imgtype)
					if err != nil {
						num = 0
					}
					state = append(state, "\n")
					state = append(state, imgtype)
					state = append(state, ": ")
					state = append(state, fmt.Sprintf("%d", num))
				}
				ctx.SendChain(message.Text(state))
			})*/
}

/*
// size 返回缓冲池指定类型的现有大小
func (p *imgpool) size(imgtype string) int {
	return len(p.pool[imgtype])
}

func (p *imgpool) push(ctx *zero.Ctx, imgtype string, illust *pixiv.Illust) {
	u := illust.ImageUrls[0]
	n := u[strings.LastIndex(u, "/")+1 : len(u)-4]
	m, err := imagepool.GetImage(n)
	var msg message.MessageSegment
	f := fileutil.BOTPATH + "/" + illust.Path(0)
	if err != nil {
		if fileutil.IsNotExist(f) {
			// 下载图片
			if err := illust.DownloadToCache(0); err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
		}
		m.SetFile(f)
		_, _ = m.Push(ctxext.SendToSelf(ctx), ctxext.GetMessage(ctx))
		msg = message.Image("file:///" + f)
	} else {
		msg = message.Image(m.String())
		if ctxext.SendToSelf(ctx)(msg) == 0 {
			msg = msg.Add("cache", "0")
		}
	}
	p.poolmu.Lock()
	p.pool[imgtype] = append(p.pool[imgtype], &msg)
	p.poolmu.Unlock()
}

func (p *imgpool) pop(imgtype string) (msg *message.MessageSegment) {
	p.poolmu.Lock()
	defer p.poolmu.Unlock()
	if p.size(imgtype) == 0 {
		return
	}
	msg = p.pool[imgtype][0]
	p.pool[imgtype] = p.pool[imgtype][1:]
	return
}

// fill 补充池子
func (p *imgpool) fill(ctx *zero.Ctx, imgtype string) {
	times := math.Min(p.max-p.size(imgtype), 2)
	p.dbmu.RLock()
	defer p.dbmu.RUnlock()
	for i := 0; i < times; i++ {
		illust := &pixiv.Illust{}
		// 查询出一张图片
		if err := p.db.Pick(imgtype, illust); err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			continue
		}
		// 向缓冲池添加一张图片
		p.push(ctx, imgtype, illust)
		process.SleepAbout1sTo2s()
	}
}

func (p *imgpool) add(ctx *zero.Ctx, imgtype string, id int64) error {
	p.dbmu.Lock()
	defer p.dbmu.Unlock()
	if err := p.db.Create(imgtype, &pixiv.Illust{}); err != nil {
		return err
	}
	ctx.SendChain(message.Text("少女祈祷中......"))
	// 查询P站插图信息
	illust, err := pixiv.Works(id)
	if err != nil {
		return err
	}
	err = imagepool.SendImageFromPool(strconv.FormatInt(illust.Pid, 10)+"_p0", illust.Path(0), func() error {
		return illust.DownloadToCache(0)
	}, ctxext.Send(ctx), ctxext.GetMessage(ctx))
	if err != nil {
		return err
	}
	// 添加插画到对应的数据库table
	if err := p.db.Insert(imgtype, illust); err != nil {
		return err
	}
	return nil
}

func (p *imgpool) remove(imgtype string, id int64) error {
	p.dbmu.Lock()
	defer p.dbmu.Unlock()
	return p.db.Del(imgtype, fmt.Sprintf("WHERE pid=%d", id))
}*/
