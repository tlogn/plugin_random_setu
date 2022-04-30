package utils

import (
	"errors"
	"github.com/FloatTech/AnimeAPI/pixiv"
	"path"
)

func DownloadImageFromPixiv(id int64, dir string) (string, error) {
	illust, err := pixiv.Works(id)
	if err != nil || len(illust.ImageUrls) == 0 || illust.ImageUrls[0] == "" {
		return "", errors.New("downloadImageFromId error")
	}
	u := illust.ImageUrls[0]
	n := path.Base(u)
	return n, illust.Download(0, path.Join(dir, n))
}
