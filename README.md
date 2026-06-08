# tm1638

This repository publishes drivers for controlling the TM1638, created using tinygo.  
The TM1638 is an LED driver controller with a key scan interface.  
このリポジトリでは、tinygoで作成したTM1638を制御するドライバを公開しています。  
TM1638はキースキャンインターフェースを備えたLEDドライバコントローラです。  
このコントローラは、「ストローブ(STB)」「クロック(CLK)」「データ(DIO)」のわずか3本の信号線で接続するだけで、複数の7セグメントLEDとスイッチ、LEDを同時に制御できます。  
また、このコントローラ内部で自動的にダイナミック点灯（スキャン表示）とキー入力のスキャンを行うため、マイコン側の負荷を大幅に軽減してくれます。  
ArduinoやMicroPythonで、このデバイスをコントロールするドライバやライブラリはたくさんあったのですが、tinygoで書かれたものが見当たらなかったので、このパッケージを自作してみました。  


## ハードウェア

以下のamazonで購入したTM1638の評価ボードを使用しました。2個で1000円弱でした。  
検証用に使用したマイコンボードは、Raspberry Pi Picoで、これらを、以下のように接続しました。  

| TM1638 | Raspberry Pi Pico |
|:-------|:------------------|
| VCC    | 3.3V              |
| GND    | GND               |
| STB    | GPIO28            |
| CLK    | GPIO27            |
| DIO    | GPIO26            |

この接続は、I2CやSPIではなく、独自のプロトコルで通信しており、汎用入出力端子に割り当てればOKです。
使用するマイコンボードの空き端子に合わせて、配線して下さい。

## ソフトウェア

開発には、以下のバージョンの go と tinygo を使用しました。

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
この中にテスト用コードがあります。

1. ターゲットボードとtm1638評価ボードを3本の信号線、電源、GND線で接続して下さい。
2. PCとターゲットボードをUSBケーブルで接続して下さい。
3. 実行したいコードのあるディレクトリ内に移動して下さい。
4. 以下のコマンドで、コンパイル&実行ファイルの転送を行って下さい。  
(-targetオプションは、使用するマイコンボードに合わせて修正して下さい。)

```bash
tinygo flash -target=pico -size=short -monitor .
```

## 解説

このドライバを使うと、以下のようなコードでTM1638を制御し、簡単に8連Lチカができます。

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
	// 入出力ピンの設定
    // ピン番号は使用するマイコンボードとtm1638の配線に合わせて変更してください）
	stbPin := machine.GP28
	clkPin := machine.GP27
	dioPin := machine.GP26

	TM1638 := tm1638.New(stbPin, clkPin, dioPin)
	TM1638.Setup()  // tm1638の初期化
	for {
			TM1638.SetLEDs(0x55)
			time.Sleep(500 * time.Millisecond)
			TM1638.SetLEDs(0xaa)
			time.Sleep(500 * time.Millisecond)
	}
}
```
## メソッドの使い方

準備中
### 初期化等の基本制御

* func New(stbPin machine.Pin, clkPin machine.Pin, dioPin machine.Pin) Device
    * tm1638 ドライバを設定します。
        - stbPin:strobePinの設定(Chip select)
        - clkPin:clockPinの設定(Clock input)
        - dioPin:dataPinの設定(DataI/O)

* func (d *Device) Setup() 
    * tm1638の初期化を行います。Newの後に必ず実行して下さい。

* 表示の明るさを設定します。設定範囲は、0-7までです。
    * func (d *Device) SetBrightness(val byte)

* func (d *Device) reset()
    * tm1638 をリセットします。

* func (d *Device) sendCMD(cmd byte)
    * tm1638 に1byteの制御命令を送ります。

* func (d *Device) sendDATA(data byte, position int)
    * tm1638に 1文字のデータを送ります.

* func (d *Device) sendDATAs(data [8]byte)
    * tm1638に 8byteのデータを送ります。

* func (d *Device) shift_out(val byte)
    * データ送信 (shift_out相当)
    * tm1638に1byteのデータを送信する。
        - val: 送信データ

* func (d *Device) shift_in() (val byte)
    * データ受信 (shift_in相当)
    * tm1638から1byteのデータを受信する。
        - val: 受信データ

### 7セグLEDの制御

* func (d *Device) Disp7SEG(segments byte, position int)
    * 指定した桁にセグメントデータを送信します。
        - segments: 1つの7セグメントLEDの点灯パターン
        - position: 点灯する7セグメントLEDの位置

* func (d *Device) Disp7SEGs(segments [8]byte)
    * 8桁分のセグメントデータを一括送信する
        - segments: 8バイトの配列 (各バイトが1つの7セグメントLEDに対応)

* func (d *Device) ScrollingText(scrollTextData []byte, interval time.Duration)
    * 文字列のスクロール表示
        - scrollTextData: 表示する文字列データ
        - interval: 文字の表示間隔(単位は、ms)

### 7セグ表示データの生成

* func (d *Device) Get7Seg(r rune) (byte, error)
    * 一文字のデータを7セグメント用のデータに変換します。
        - r: 一文字のデータ,rune型です。byte型ではないことに注意して下さい。

* func (d *Device) StrTo7Seg(text string) [8]byte
    * 文字列を7セグメント用の8byteデータに変換します。

* func (d *Device) IntTo7Seg(n int) [8]byte
    * 整数を7セグメント用の8byteデータに変換します。
    * 負の数や桁溢れ等には対応していません。

### 8連LEDの制御

* func (d *Device) SetLED(value byte, position byte)
    * 指定する1つのLEDの値を設定します。
        - value: 設定値
        - position: LEDの番号

* func (d *Device) SetLEDs(value byte)
    * LEDの点灯パターンを設定します。
        - value: 1byteのデータで構成される点灯パターン

* func (d *Device) ClearLEDs()
    * 全LEDを消灯する。

### 8連キーからの入力

* func (d *Device) ScanKeys() byte
    * 8個のキーのうち、押されているキーの状態を1byteのデータとして返します。LSBはSW1です。

## サンプルコードについて

[Blink8LED/main.go](Blink8LED/main.go)  

    8連LEDを点滅させます。

[DispNum/main.go](DispNum/main.go)  

    7セグに数字の表示します。

[HelloWorld/main.go](HelloWorld/main.go)  

    7セグに文字列"Hello"を疑似アルファベットで表示します。

[Keytest/main.go](Keytest/main.go)  

    8個のキーの入力状態を読取ります。押されているキーに対応するLEDを消灯します。

[LEDtest/main.go](LEDtest/main.go)  

    8連LEDをいろいろなパターンで点滅させます。

[ScrollingText/main.go](ScrollingText/main.go)  

    7セグに文字列を疑似アルファベットでスクロール表示します。

## このパッケージのドキュメント

以下のリンクをご覧ください。
[package tm1638のドキュメント](https://pkg.go.dev/github.com/triring/tm1638)
