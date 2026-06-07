package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o LEDtest.uf2 --target=pico --size short .
*/

import (
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
	for {
		for i := 0; i < 8; i++ {
			TM1638.SetLED(0x7, byte(i))
			time.Sleep(1000 * time.Millisecond)
			//	TM1638.SetLED(0x00, byte(i))
		}
		time.Sleep(1000 * time.Millisecond)
		TM1638.ClearLEDs()
		time.Sleep(1000 * time.Millisecond)

		for i := 0; i < 50; i++ {
			TM1638.SetLEDs(0x55)
			time.Sleep(200 * time.Millisecond)
			TM1638.SetLEDs(0xaa)
			time.Sleep(200 * time.Millisecond)
		}
		TM1638.ClearLEDs()
		time.Sleep(1000 * time.Millisecond)
	}
}
