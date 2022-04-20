package searchsetu

import (
	"fmt"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	api      = "https://api.lolicon.app/setu/v2"
	capacity = 10
)

func init() {
	engine := control.Register("setutime", &control.Options{
		DisableOnDefault: false,
		Help:             "- (R18)来点[xxx]\n",
		PublicDataFolder: "SearchSetuTime",
	})
	engine.ApplySingle(ctxext.DefaultSingle).OnRegex("^(R18)?(来点)(\\s\\S)+").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			imgtype := ctx.State["regex_matched"].([]string)
			fmt.Println(imgtype)
		})
}
