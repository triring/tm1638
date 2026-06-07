package tm1638

import (
	"errors"
	// "fmt"
)

// 7セグで擬似的に表示する英数字データの定義
// 文字コードから、7segの表示データへの変換テーブル
var CharTo7Seg map[rune]byte = map[rune]byte{
	//Display: Hex  // Bin
	'0': 0x3F, // 0b00111111
	'1': 0x06, // 0b00000110
	'2': 0x5B, // 0b01011011
	'3': 0x4F, // 0b01001111
	'4': 0x66, // 0b01100110
	'5': 0x6D, // 0b01101101
	'6': 0x7D, // 0b01111101
	'7': 0x07, // 0b00000111
	'8': 0x7F, // 0b01111111
	'9': 0x6F, // 0b01101111
	'A': 0x77, // 0b01110111
	'B': 0x7C, // 0b01111100
	'C': 0x39, // 0b00111001
	'D': 0x5E, // 0b01011110
	'E': 0x79, // 0b01111001
	'F': 0x71, // 0b01110001
	'G': 0x3D, // 0b00111101
	'H': 0x76, // 0b01110110
	'I': 0x06, // 0b00000110
	'J': 0x1E, // 0b00011110
	'K': 0x76, // 0b01110110
	'L': 0x38, // 0b00111000
	'M': 0x55, // 0b01010101
	'N': 0x54, // 0b01010100
	'O': 0x3F, // 0b00111111
	'P': 0x73, // 0b01110011
	'Q': 0x67, // 0b01100111
	'R': 0x50, // 0b01010000
	'S': 0x6D, // 0b01101101
	'T': 0x78, // 0b01111000
	'U': 0x3E, // 0b00111110
	'V': 0x1C, // 0b00011100
	'W': 0x2A, // 0b00101010
	'X': 0x76, // 0b01110110
	'Y': 0x6E, // 0b01101110
	'Z': 0x5B, // 0b01011011
	'a': 0x5F, // 0b01011111
	'b': 0x7C, // 0b01111100
	'c': 0x58, // 0b01011000
	'd': 0x5E, // 0b01011110
	'e': 0x79, // 0b01111001
	'f': 0x71, // 0b01110001
	'g': 0x3D, // 0b00111101
	'h': 0x74, // 0b01110100
	'i': 0x04, // 0b00000100
	'j': 0x1E, // 0b00011110
	'k': 0x76, // 0b01110110
	'l': 0x30, // 0b00110000
	'm': 0x55, // 0b01010101
	'n': 0x54, // 0b01010100
	'o': 0x5C, // 0b01011100
	'p': 0x73, // 0b01110011
	'q': 0x67, // 0b01100111
	'r': 0x50, // 0b01010000
	's': 0x6D, // 0b01101101
	't': 0x78, // 0b01111000
	'u': 0x3E, // 0b00111110
	'v': 0x1C, // 0b00011100
	'w': 0x2A, // 0b00101010
	'x': 0x76, // 0b01110110
	'y': 0x6E, // 0b01101110
	'z': 0x5B, // 0b01011011
	' ': 0x00, // 0b00000000   blank
	'-': 0x40, // 0b01000000
	'*': 0x63, // 0b01100011
	'=': 0x48, // 0b01001000
	'.': 0x80, // 0b10000000
}

// Converts a single character to 7-segment data.
//
// 一文字のデータを7セグメント用のデータに変換します。
// r: 一文字のデータ,rune型です。byte型ではないことに注意して下さい。
func (d *Device) Get7Seg(r rune) (byte, error) {
	val, exists := CharTo7Seg[r]
	if !exists {
		return 0, errors.New("character not found")
	}
	return val, nil
}

// Converts a string to 8-byte data for 7-segment displays.
//
// 文字列を7セグメント用の8byteデータに変換します。
func (d *Device) StrTo7Seg(text string) [8]byte {
	txtBuf := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var count int = 0
	for i, r := range text {
		if i > 7 {
			break
		}
		// i はバイト位置のインデックス、r が rune 型の文字
		// fmt.Printf("インデックス: %d, rune値: %U, 文字: %c\n", i, r, r)
		b, _ := d.Get7Seg(r)
		txtBuf[i] = b
		count++
	}
	//	fmt.Println(txtBuf)
	return txtBuf
}

// Converts an integer to 8-byte data for 7-segment displays.
//
// 整数を7セグメント用の8byteデータに変換します。
// 負の数や桁溢れ等には対応していません。
func (d *Device) IntTo7Seg(n int) [8]byte {
	txtBuf := [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//	var txtBuf []byte = []byte{	0, 0, 0, 0, 0, 0, 0, 0 }
	var ch rune
	for i := 7; i >= 0; i-- {
		switch n % 10 {
		case 0:
			ch = '0'
		case 1:
			ch = '1'
		case 2:
			ch = '2'
		case 3:
			ch = '3'
		case 4:
			ch = '4'
		case 5:
			ch = '5'
		case 6:
			ch = '6'
		case 7:
			ch = '7'
		case 8:
			ch = '8'
		case 9:
			ch = '9'
		default:
			ch = ' '
		}
		seg, _ := d.Get7Seg(ch)
		txtBuf[i] = seg
		n = n / 10
		if n == 0 {
			break
		}
	}
	//	fmt.Println(txtBuf)
	return txtBuf
}
