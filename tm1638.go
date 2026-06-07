package tm1638

import (
	"machine"
	"time"
)

// Device wraps an I2C connection to a BMP180 device.
type Device struct {
	stbPin     machine.Pin
	clkPin     machine.Pin
	dioPin     machine.Pin
	brightness byte
}

// 表示アドレス設定コマンド (通常 0xC0 から開始)
const (
	cmdActivate     = 0x8f // activate
	cmdDataSettings = 0x40 // データ設定コマンド (自動インクリメント)
	cmdDataRead     = 0x42 // データ設定コマンド (1byte読み出し。)
	cmdSetLed       = 0x44 // データ設定コマンド (LEDの点灯情報)
	cmdAddressStart = 0xC0 // アドレス設定コマンド (00H番地)
	//	cmdDisplayConfig = 0x88 // ディスプレイ制御 (表示ON、輝度設定)
	cmdDisplayBright = 0x80 // ディスプレイ制御 (輝度設定)
	cmdDisplayOnOff  = 0x08 // ディスプレイ制御 (表示ON_OFF)
	// 通信タイミングの定数（1us = 1000ナノ秒）
	// 将来、より高速なMPUに変更した場合はここを増やして調整します。
	tmDelay = 2 * time.Microsecond
)

// New creates a new package tm1638 device.
//
// tm1638 ドライバを設定します。
// stbPin:strobePinの設定(Chip select)
// clkPin:clockPinの設定(Clock input)
// dioPin:dataPinの設定(DataI/O)
func New(stbPin machine.Pin, clkPin machine.Pin, dioPin machine.Pin) Device {
	var dev Device
	// Configure sets up the pins.
	stbPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	clkPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dioPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dev = Device{stbPin: stbPin, clkPin: clkPin, dioPin: dioPin}
	dev.brightness = 0x07 // 初期値として、明るさを最大輝度に設定する。
	return dev
}

// Setting up the tm1638.
// 
// tm1638の初期化を行います。Newの後にか鳴らす実行して下さい。
func (d *Device) Setup() {
	d.sendCMD(cmdActivate) // activate	0x8f
	d.reset()
}

// Set the display brightness 0-7.
// 
// 表示のかkるさを設定します。設定範囲は、0-7までです。
func (d *Device) SetBrightness(val byte) {
	// brightness 0 = 1/16th pulse width
	// brightness 7 = 14/16th pulse width
	if 0 > val {
		return
	}
	if val > 7 {
		return
	}
	d.brightness = val
	// d.sendCMD(cmdDisplayBright | cmdDisplayOnOff | d.brightness) // 0x88(ON) | 0x07(最大輝度)
}

// Resetting the tm1638.
// 
// tm1638 をリセットします。
func (d *Device) reset() {
	d.sendCMD(cmdDataSettings) // set auto increment mode
	d.stbPin.Low()
	d.shift_out(cmdAddressStart) // set starting address to 0
	for i := 0; i < 16; i++ {
		d.shift_out(0x00)
	}
	d.stbPin.High()
}

// Send 1-byte control command to tm1638.
//
// tm1638 に1byteの制御命令を送ります。
func (d *Device) sendCMD(cmd byte) {
	// 1. データ設定コマンド送信 (書き込みモード、アドレス自動加算)
	d.stbPin.Low()
	time.Sleep(tmDelay)
	d.shift_out(cmd)
	d.stbPin.High()
	time.Sleep(tmDelay)
}

// Send a 1-character data to tm1638.
//
// tm1638に 1文字のデータを送ります.
func (d *Device) sendDATA(data byte, position int) {
	d.stbPin.Low()
	time.Sleep(tmDelay)
	// 開始アドレス 0xC0 を送信
	d.shift_out(cmdAddressStart + byte(position)*2)
	// TM1638は1つの桁に対して「セグメントデータ」と「グリッドデータ」で2バイト分
	// 使う構成のため、8桁表示の場合は 16回 shift_out する必要があります。
	// (多くのモジュールでは偶数番地にセグメント、奇数番地は未使用かLED等)
	d.shift_out(data) // セグメントデータ
	d.stbPin.High()
	time.Sleep(tmDelay)
}

// Send an 8-byte data to tm1638.
//
// tm1638に 8byteのデータを送ります。
func (d *Device) sendDATAs(data [8]byte) {
	d.stbPin.Low()
	time.Sleep(tmDelay)
	// 開始アドレス 0xC0 を送信
	d.shift_out(cmdAddressStart)
	// TM1638は1つの桁に対して「セグメントデータ」と「グリッドデータ」で2バイト分
	// 使う構成のため、8桁表示の場合は 16回 shift_out する必要があります。
	// (多くのモジュールでは偶数番地にセグメント、奇数番地は未使用かLED等)
	for i := 0; i < 8; i++ {
		d.shift_out(data[i]) // セグメントデータ
		d.shift_out(0x00)
		//	d.shift_out(0xFF)    // 空データ(またはモジュール上のLED制御)
	}
	d.stbPin.High()
	time.Sleep(tmDelay)
}

// ScanKeys Return a byte representing which keys are pressed. LSB is SW1
// 
// 8個のキーのうち、押されているキーの状態を1byteのデータとして返します。LSBはSW1です。
func (d *Device) ScanKeys() byte {
	var keys byte = 0
	var i byte
	d.stbPin.Low()
	time.Sleep(tmDelay)
	d.shift_out(cmdDataRead)
	time.Sleep(tmDelay)
	d.dioPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	for i = 0; i < 4; i++ {
		keys |= (d.shift_in() << i)
	}
	d.dioPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.stbPin.High()
	return keys
}

// Set the value of a single LED.
//
// 指定する1つのLEDの値を設定します。
// value: 設定値
// position: LEDの番号
func (d *Device) SetLED(value byte, position byte) {
	d.dioPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d.shift_out(cmdSetLed)
	d.stbPin.Low()
	d.shift_out(0xC1 + (position << 1))
	d.shift_out(value)
	d.stbPin.High()
}

// Set the LED lighting pattern.
//
// LEDの点灯パターンを設定します。
// value: 1byteのデータで構成される点灯パターン
func (d *Device) SetLEDs(value byte) {
	var position byte
	var mask byte
	for position = 0; position < 8; position++ {
		mask = 0x1 << position
		if 0 == (value & mask) {
			d.SetLED(0, position)
		} else {
			d.SetLED(1, position)
		}
	}
}

// Clear all. Write zeros to each address.
//
// 全LEDを消灯する。
func (d *Device) ClearLEDs() {
	for i := 0; i < 8; i++ {
		d.SetLED(0, byte(i))
	}
}

// Sends 8-digit segment data in a single batch.
//
// 8桁分のセグメントデータを一括送信する
// segments: 8バイトの配列 (各バイトが1つの7セグメントLEDに対応)
func (d *Device) Disp7SEGs(segments [8]byte) {
	// 1. データ設定コマンド送信 (書き込みモード、アドレス自動加算)
	d.sendCMD(cmdDataSettings)
	// 2. 表示アドレスの指定とデータ送信
	d.sendDATAs(segments)
	// 3. 表示の有効化と輝度設定
	d.sendCMD(cmdDisplayBright | cmdDisplayOnOff | d.brightness) // 0x88(ON) | 0x07(最大輝度)
	// d.sendCMD(cmdDisplayBright | cmdDisplayOnOff | 0x07) // 0x88(ON) | 0x07(最大輝度)
	// cmdDisplayBright = 0x80 // ディスプレイ制御 (輝度設定)
	// cmdDisplayOnOff  = 0x08 // ディスプレイ制御 (表示ON_OFF)
}

// Sends segment data to a specified digit.
//
// 指定した桁にセグメントデータを送信します。
// segments: 1つの7セグメントLEDの点灯パターン
// position: 点灯する7セグメントLEDの位置
func (d *Device) Disp7SEG(segments byte, position int) {
	// 1. データ設定コマンド送信 (書き込みモード、アドレス自動加算)
	d.sendCMD(cmdDataSettings)
	// 2. 表示アドレスの指定とデータ送信
	d.sendDATA(segments, position)
	// 3. 表示の有効化と輝度設定
	d.sendCMD(cmdDisplayBright | cmdDisplayOnOff | d.brightness) // 0x88(ON) | 0x07(最大輝度)
	// d.sendCMD(cmdDisplayBright | cmdDisplayOnOff | 0x07) // 0x88(ON) | 0x07(最大輝度)
	// cmdDisplayBright = 0x80 // ディスプレイ制御 (輝度設定)
	// cmdDisplayOnOff  = 0x08 // ディスプレイ制御 (表示ON_OFF)
}

// Scroll text display
//
// 文字列のスクロール表示
// scrollTextData: 表示する文字列データ
// interval: 文字の表示間隔(単位は、ms)
func (d *Device) ScrollingText(scrollTextData []byte, interval time.Duration) {
	scrollLength := len(scrollTextData)
	var index int = 0
	var dispData [8]byte
	for {
		for i := 0; i < 8; i++ {
			dispData[i] = scrollTextData[index+i]
		//	fmt.Printf("%c ", dispData[i])
		}
	//	fmt.Printf("\n")
		d.Disp7SEGs(dispData)
		time.Sleep(interval * time.Millisecond)
		index++
		if index > (scrollLength - 8) {
			break
		}
	}
}

// Sends 1 byte of data to tm1638.
//
// データ送信 (shift_out相当)
// tm1638に1byteのデータを送信する。
// val: 送信データ
func (d *Device) shift_out(val byte) {
	d.dioPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := 0; i < 8; i++ {
		// LSB First
		bit := (val >> i) & 1
		d.dioPin.Set(bit == 1)
		time.Sleep(tmDelay)

		d.clkPin.High()
		time.Sleep(tmDelay)

		d.clkPin.Low()
		time.Sleep(tmDelay)
	}
}

// Receives 1 byte of data from tm1638.
//
// データ受信 (shift_in相当)
// tm1638から1byteのデータを受信する。
// val: 受信データ
func (d *Device) shift_in() (val byte) {
	var value byte
	d.dioPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	time.Sleep(tmDelay)

	for i := 0; i < 8; i++ {
		d.clkPin.Low()
		time.Sleep(tmDelay)

		d.clkPin.High()
		// クロック立ち上がりで値を読み取る
		if d.dioPin.Get() {
			value |= (1 << i)
		}
		time.Sleep(tmDelay)
	}
	d.clkPin.Low()
	return value
}
