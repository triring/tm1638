package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o CalcPi.uf2 --target=pico --size short .
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

	var numerator int = 355
	var denominator int = 113
	var pi float32

	pi = float32(numerator) / float32(denominator)
	fmt.Printf("%f", pi)
	for {
		TM1638.Disp7SEGs(TM1638.IntTo7Seg(numerator))
		fmt.Printf("%d", numerator)		
		time.Sleep(2 * time.Second)
		TM1638.Disp7SEGs(TM1638.StrTo7Seg("     div"))
		fmt.Printf("%s", "   div  ")		
		time.Sleep(2 * time.Second)
		TM1638.Disp7SEGs(TM1638.IntTo7Seg(denominator))
		fmt.Printf("%d", denominator)		
		time.Sleep(2 * time.Second)
		TM1638.Disp7SEGs(TM1638.StrTo7Seg("   EQUAL"))
		fmt.Printf("%s", " EQUAL  ")		
		time.Sleep(2 * time.Second)
		pi_str := fmt.Sprintf("%.6f", pi)
		TM1638.Disp7SEGs(TM1638.StrTo7Seg(pi_str))
		fmt.Printf("%s\n", pi_str)
		time.Sleep(8 * time.Second)
	}
}
