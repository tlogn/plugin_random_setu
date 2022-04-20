package randsetu

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
	"time"
)

func TestRandSetuApi(t *testing.T) {
	fmt.Println(randSetuApi())
}

func TestDownloadImageFromID(t *testing.T) {
	fmt.Println(downloadImageFromID(97721128))
}

func TestRandDownloadImage(t *testing.T) {
	fmt.Println(randDownloadImage())
}

func TestFIFOInit(t *testing.T) {
	imgFIFO.init()
	imgName := imgFIFO.get()
	pathName, _ := filepath.Abs(path.Join(imgPath, imgName))
	pathName = "file://" + pathName
	fmt.Println(pathName)
	time.Sleep(time.Second * 5)
	fmt.Println(pathName)
	time.Sleep(time.Second * 500)
}
