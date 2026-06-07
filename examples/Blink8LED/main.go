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
	"github.com/triring/tm1638"
)

func main() {
	// ピンの初期化（ピン番号はPicoの実際の配線に合わせて変更してください）
	stbPin := machine.GP28
	clkPin := machine.GP27
	dioPin := machine.GP26

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()
	for {
			TM1638.SetLEDs(0x55)
			time.Sleep(500 * time.Millisecond)
			TM1638.SetLEDs(0xaa)
			time.Sleep(500 * time.Millisecond)
	}
}
