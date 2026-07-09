package main

/*
> go mod init main
> go mod tidy
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o KnightRider.uf2 --target=pico --size short .

アメリカの特撮テレビドラマ『ナイトライダー』（Knight Rider）に
「ナイト2000（Knight Industries Two Thousand）」（愛称：キット / K.I.T.T.）
という車両が登場します。
この車両のフロントで点滅する赤い光「スキャナー（ナイトフラッシャー）」を再現してみました。
*/

import (
	"machine"
	"time"
	// "tm1638" // ローカルのディレクトリに置かれたtm1638のパッケージをインポートする場合
	"github.com/triring/tm1638" // githubで公開しているパッケージをインポートする場合
)

var (
	stbPin     machine.Pin
	clkPin     machine.Pin
	dioPin     machine.Pin
	Afterimage [4]byte = [4]byte{1, 2, 3, 4}
)

func main() {
	// ピンの初期化（ピン番号は使用しているマイコンボードの実際の配線に合わせて変更してください）
	stbPin = machine.GP28
	clkPin = machine.GP27
	dioPin = machine.GP26
	duration := 2

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()

	Flickering := func(index byte, duration int) {
		TM1638.SetLED(0x7, index)
		TM1638.SetLED(0x7, Afterimage[0])
		TM1638.SetLED(0x7, Afterimage[1])
		TM1638.SetLED(0x7, Afterimage[2])
		TM1638.SetLED(0x7, Afterimage[3])
		time.Sleep(time.Duration(duration) * time.Millisecond)
		TM1638.SetLED(0x0, Afterimage[3])
		time.Sleep(time.Duration(duration*2) * time.Millisecond)
		TM1638.SetLED(0x0, Afterimage[2])
		time.Sleep(time.Duration(duration*4) * time.Millisecond)
		TM1638.SetLED(0x0, Afterimage[1])
		time.Sleep(time.Duration(duration*8) * time.Millisecond)
		TM1638.SetLED(0x0, Afterimage[0])
		time.Sleep(time.Duration(duration*16) * time.Millisecond)
		TM1638.SetLED(0x0, index)
		// 残像の情報をシフトしていく。
		Afterimage[3] = Afterimage[2]
		Afterimage[2] = Afterimage[1]
		Afterimage[1] = Afterimage[0]
		Afterimage[0] = index
	}
	// 残像を残しながら、LEDの点滅が左右に移動していく。
	for {
		for i := 0; i < 8; i++ {
			Flickering(byte(i), duration)
		}
		for i := 7; i >= 0; i-- {
			Flickering(byte(i), duration)
		}
	}
}
