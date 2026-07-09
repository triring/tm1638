package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o DispNum.uf2 --target=pico --size short .
*/

import (
	"fmt"
	"machine"
	"time"
	// "tm1638" // ローカルのディレクトリに置かれたtm1638のパッケージをインポートする場合
	"github.com/triring/tm1638" // githubで公開しているパッケージをインポートする場合
)

var (
	stbPin machine.Pin
	clkPin machine.Pin
	dioPin machine.Pin
)

func main() {
	// ピンの初期化（ピン番号はPicoの実際の配線に合わせて変更してください）
	stbPin = machine.GP28
	clkPin = machine.GP27
	dioPin = machine.GP26

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()

	disp7Seg := func(text string) {
		TM1638.Disp7SEGs(TM1638.StrTo7Seg(text))
		fmt.Println(text)
	}
	for {
		disp7Seg("HELLO")
		time.Sleep(2 * time.Second)
		disp7Seg(" 01234  ")
		time.Sleep(2 * time.Second)
		disp7Seg(" 56789  ")
		time.Sleep(2 * time.Second)
		disp7Seg("ABCDEFGH")
		time.Sleep(5 * time.Second)
		disp7Seg("IJKLMNOP")
		time.Sleep(5 * time.Second)
		disp7Seg("QRSTUVW ")
		time.Sleep(5 * time.Second)
		disp7Seg(" XYZ*-=!")
		time.Sleep(5 * time.Second)
		disp7Seg("abcdefgh")
		time.Sleep(5 * time.Second)
		disp7Seg("ijklmnop")
		time.Sleep(5 * time.Second)
		disp7Seg("qrstuvw ")
		time.Sleep(5 * time.Second)
		disp7Seg(" xyz.*-=!")
		time.Sleep(5 * time.Second)
	}
}
