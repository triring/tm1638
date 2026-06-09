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

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)

	// ---------------------------------------------------------// 使用例// ---------------------------------------------------------func runDisplayExample() {
	// 7セグでアルファベットをを表示するためのセグメントパターン
	// 限られた7セグ点灯パターンの中でアルファベットに近いパターンを選んでいる。
	textFlashing := [8]byte{
		/*F*/ /*L*/ /*A*/ /*S*/ /*H*/ /*I*/ /*L*/ /*G*/
		0x71, 0x38, 0x77, 0x6D, 0x76, 0x06, 0x37, 0x3D,
	}
	var b byte = 0x01
	// 明るさは,0から7までの範囲で設定可能。
	for {
		// 徐々に表示を明るくする。
		for b = 0; b <= 7; b++ {
			TM1638.SetBrightness(b)
			TM1638.Disp7SEGs(textFlashing) 
			fmt.Println("Brightness\t",b)
			time.Sleep(1000 * time.Millisecond)
		}
		// 次第に表示を暗くする。
		for b = 7; b > 0; b-- {
			TM1638.SetBrightness(b)
			TM1638.Disp7SEGs(textFlashing)
			fmt.Println("Brightness\t",b)
			time.Sleep(1000 * time.Millisecond)
		}		
	}
}
