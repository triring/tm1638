package main

/*
> go mod init main
> go mod tidy
> go get github.com/triring/tm1638
> tinygo flash --target=pico --size short -monitor .
> tinygo build -o TowerOfHanoi.uf2 --target=pico --size short .
*/

import (
	"fmt"
	"machine"
	"runtime"
	"time"
	// "tm1638" // ローカルのディレクトリに置かれたtm1638のパッケージをインポートする場合
	"github.com/triring/tm1638" // githubで公開しているパッケージをインポートする場合
)

var (
	stbPin machine.Pin
	clkPin machine.Pin
	dioPin machine.Pin
)

var (
	towerA int = 0
	towerB int = 0
	towerC int = 0
)

func main() {
	var wait  time.Duration =  time.Duration(1000)
	// ピンの初期化（ピン番号はPicoの実際の配線に合わせて変更してください）
	stbPin = machine.GP28
	clkPin = machine.GP27
	dioPin = machine.GP26

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()

	// 7セグで"HANOI"をアルファベットで表示するためのセグメントパターン
	// 限られた7セグ点灯パターンの中でアルファベットに近いパターンを選んでいる。
	var textHanoi [8]byte
	var diskFont [4]byte

	textHanoi[0] = 0x76 // H
	textHanoi[1] = 0x77 // A
	textHanoi[2] = 0x37 // N
	textHanoi[3] = 0x3F // O
	textHanoi[4] = 0x30 // I
	textHanoi[5] = 0x08
	textHanoi[6] = 0x48
	textHanoi[7] = 0x49

	diskFont[0] = 0x00
	diskFont[1] = 0x08
	diskFont[2] = 0x48
	diskFont[3] = 0x49
	time.Sleep(2 * time.Second)
	TM1638.Disp7SEGs(textHanoi)
	time.Sleep(5 * time.Second)

	fmt.Printf("Step 1\n")
	moveDisk := func(idx int, from string, to string) {
		switch from {
		case "A":
			towerA--
		case "B":
			towerB--
		case "C":
			towerC--
		}
		switch to {
		case "A":
			towerA++
		case "B":
			towerB++
		case "C":
			towerC++
		}
		fmt.Printf("%d_%d_%d\n", towerA, towerB, towerC)
		textHanoi[idx    ] = diskFont[towerA]
		textHanoi[idx + 1] = diskFont[towerB]
		textHanoi[idx + 2] = diskFont[towerC]
		TM1638.Disp7SEGs(textHanoi)
		time.Sleep(wait * time.Millisecond)
	}
	fmt.Printf("Step 2\n")

	// ハノイの塔を解く再帰関数
	// 無名関数（クロージャ）を変数に代入することで、関数の中でローカルな関数を実現する。
	// 1. 変数を先に宣言しておく
	var hanoi func(idx int, n int, from, to, via string) 
	// 2. 変数に関数を代入する
	// n: 円盤の枚数
	// from: 移動元の柱
	// to: 移動先の柱
	// via: 経由する柱
	hanoi = func(idx int, n int, from, to, via string) {
		if n == 1 {
			fmt.Printf("円盤 1 を %s から %s へ移動:", from, to)
			moveDisk(idx, from, to)
			runtime.GC()	// 強制的にガベージコレクション（GC）を実行させる
			return
		}
		// 1. n-1 枚の円盤を、from から via へ移動する
		hanoi(idx, n-1, from, via, to)

		// 2. 残った最大の円盤を、from から to へ移動する
		fmt.Printf("円盤 %d を %s から %s へ移動:", n, from, to)
		moveDisk(idx, from, to)
		// 3. via に避難させた n-1 枚の円盤を、to へ移動する
		hanoi(idx, n-1, via, to, from)
	}
	fmt.Printf("Step 3\n")

	for i := 0; i < 5; i++ {
		fmt.Printf("--------\n")
		towerA = 3
		towerB = 0
		towerC = 0
		hanoi(5, 3, "A", "C", "B")
		fmt.Printf("--------\n")
		towerA = 0
		towerB = 0
		towerC = 3
		hanoi(5, 3, "C", "A", "B")
	}
	// 表示速度を上げる。
	wait = time.Duration(250)
	towerA = 3
	towerB = 0
	towerC = 0
	hanoi(5, 3, "A", "C", "B")
	textHanoi[5] = 0x00
	textHanoi[6] = 0x00
	textHanoi[7] = 0x00
	for {
		towerA = 0
		towerB = 0
		towerC = 3
		hanoi(4, 3, "C", "A", "B")

		towerA = 0
		towerB = 0
		towerC = 3
		hanoi(2, 3, "C", "A", "B")

		towerA = 0
		towerB = 0
		towerC = 3
		hanoi(0, 3, "C", "A", "B")

		towerA = 3
		towerB = 0
		towerC = 0
		hanoi(0, 3, "A", "C", "B")

		towerA = 3
		towerB = 0
		towerC = 0
		hanoi(2, 3, "A", "C", "B")

		towerA = 3
		towerB = 0
		towerC = 0
		hanoi(4, 3, "A", "C", "B")
	}
}
