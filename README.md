# tm1638

This repository publishes drivers for controlling the TM1638, created using tinygo.  
The TM1638 is an LED driver controller with a key scan interface.  
このリポジトリでは、tinygoで作成したTM1638を制御するドライバを公開しています。  
TM1638はキースキャンインターフェースを備えたLEDドライバコントローラです。  
このコントローラは、「ストローブ(STB)」「クロック(CLK)」「データ(DIO)」のわずか3本の信号線で接続するだけで、複数の7セグメントLEDとスイッチ、LEDを同時に制御できます。  
また、このコントローラ内部で自動的にダイナミック点灯（スキャン表示）とキー入力のスキャンを行うため、マイコン側のの負荷軽減してくれます。  
ArduinoやMicroPythonで、このデバイスをコントロールするデバイスドライバはたくさんあったのですが、tinygoで書かれたものが見当たらなかったので、自作してみました。  


## ハードウェア

以下のamazonで購入したTM1638の評価ボードを使用しました。2個で1000円弱でした。
今回、検証用に使用したマイコンボードは、Raspberry Pi Picoです。
配線は、以下のように接続しました。この接続は、I2CでもSPIでもないようなので、汎用入出力端子であればOKです。
使用するマイコンボードの空き端子に合わせて、配線して下さい。

| TM1638 | Raspberry Pi Pico |
|:-------|:------------------|
| VCC    | 3.3V              |
| GND    | GND               |
| STB    | GPIO28            |
| CLK    | GPIO27            |
| DIO    | GPIO26            |

## ソフトウェア

開発には、以下のバージョンのgo と tinygo を使用しました。

    > go version
    go version go1.26.4 windows/amd64
    > tinygo version
    tinygo version 0.41.1 windows/amd64 (using go version go1.26.4 and LLVM version 20.1.1)

## 使用方法

以下のコマンドで、このリポジトリの内容をローカルにコピーして下さい。
```bash
git clone https://github.com/triring/tm1638.git
```

コピーされたtm1638ディレクトリ内のファイル構成
```bash

D:.
|   .gitignore
|   go.mod
|   LICENSE
|   README.md
|   registers.go
|   tm1638.go
|
+---examples
    +---DispNum
    |       main.go
    |
    +---HelloWorld
    |       main.go
    |
    +---Keytest
    |       main.go
    |
    +---LEDtest
    |       main.go
    |
    \---ScrollingText
            main.go
```

コピーしたディレクトリ内に、examplesディレクトリがあります。

1. ターゲットボードとtm1638評価ボードを3本の信号線、電源、GND線で接続して下さい。
2. PCとターゲットボードを接続して下さい。
3. 実行したいコードのあるディレクトリ内に移動して下さい。
4. 以下のコマンドで、コンパイル&実行ファイルの転送を行って下さい。  
(-targetオプションは、使用するマイコンボードに合わせて修正して下さい。)

```bash
tinygo flash -target=pico -size=short -monitor .
```

## 解説

このドライバを使うと、以下のようなコードでTM1638を制御し、簡単にLチカができます。

1. "github.com/triring/tm1638"をインポートする。
2. tm1638の入出力ピンを設定する。
3. tm1638を初期化する。
4. 必要なメソッドを呼び出す

```go
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
```
## メソッドの使い方


### 7セグLEDの制御

### 7セグ表示データの生成

### 8連LEDの制御


### 8連キーからの入力

## このパッケージのドキュメント

[package tm1638のドキュメント](https://pkg.go.dev/github.com/triring/tm1638)
