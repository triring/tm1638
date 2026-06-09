package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o ScrollingText.uf2 --target=pico --size short .
*/

import (
	"fmt"
	"machine"
	"time"
	//"tm1638" // ローカルのディレクトリに置かれたtm1638のパッケージをインポートする場合
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

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)

	// 7セグでアルファベットを表示するためのセグメントパターン
	// 限られた7セグ点灯パターンの中でアルファベットに近いパターンを選んでいる。
	// ASAP
	textASAP := [8]byte{0x00, 0x00, 0x77, 0x6D, 0x77, 0x73, 0x00, 0x00}
	// Alcohol Solves All Problems.
	scrollTextData := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x77, 0x30, 0x58, 0x5C, 0x74, 0x5C, 0x30, 0x00, // Alcohol
		0x6D, 0x5C, 0x30, 0x1C, 0x79, 0x6D, 0x00, // SolveS
		0x77, 0x38, 0x38, 0x00, // All
		0x73, 0x50, 0x5C, 0x7C, 0x30, 0x79, 0x55, 0x6D, 0x80, // ProblemS.
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	// Pulse pattern
	scrollPulseData := []byte{
		0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54,
		0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54, 0x08, 0x54,
	}

	time.Sleep(3000 * time.Millisecond)

	fmt.Printf("ASAP\n")
	TM1638.Disp7SEGs(textASAP)
	time.Sleep(5000 * time.Millisecond)
	fmt.Printf("Alcohol Solves All Problems.\n")
	TM1638.ScrollingText(scrollTextData, 750)	//	表示インターバル	750ms
	TM1638.ScrollingText(scrollTextData, 500)	//	表示インターバル	500ms
	TM1638.ScrollingText(scrollTextData, 250)	//	表示インターバル	250ms
	TM1638.Disp7SEGs(textASAP)
	time.Sleep(2000 * time.Millisecond)
	fmt.Printf("Pulse pattern\n")
	for i := 0; i < 20; i++ {
		TM1638.ScrollingText(scrollPulseData, 100)
	}
}
