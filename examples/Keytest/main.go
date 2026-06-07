package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o HelloWorld.uf2 --target=pico --size short .
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

	stbPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	clkPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dioPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	var key_state byte
	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()
	for {
		key_state = TM1638.ScanKeys() // キーの状態を読み取る。
		fmt.Printf("%02x, %08b\n", key_state, key_state)
		// ビット反転を行い、キーが押されたら、LEDを消灯する。
		TM1638.SetLEDs(^key_state)
		time.Sleep(100 * time.Millisecond)
	}
}
