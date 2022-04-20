package randsetu

import (
	"fmt"
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
	time.Sleep(time.Second * 500)
}
