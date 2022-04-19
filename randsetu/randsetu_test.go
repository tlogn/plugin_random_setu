package randsetu

import (
	"fmt"
	"testing"
)

func TestRandSetuApi(t *testing.T) {
	fmt.Println(randSetuApi())
}

func TestDownloadImageFromID(t *testing.T) {
	fmt.Println(downloadImageFromID(97721128))
}

func TestRandDownloadImage(t *testing.T) {
	//fmt.Println(randDownloadImage())
	//fmt.Println(imagepool.GetImage("44921453_p0.jpg"))
}
