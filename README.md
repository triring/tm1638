# tm1638
This repository publishes drivers for controlling the TM1638, created using tinygo.  
The TM1638 is an LED driver controller with a key scan interface.  
このリポジトリでは、tinygoで作成したTM1638を制御するドライバを公開しています。  
TM1638はキースキャンインターフェースを備えたLEDドライバコントローラです。  
このコントローラは、「ストローブ(STB)」「クロック(CLK)」「データ(DIO)」のわずか3本の信号線で接続するだけで、複数の7セグメントLEDとスイッチ、LEDを同時に制御できます。  
また、このコントローラ内部で自動的にダイナミック点灯（スキャン表示）とキー入力のスキャンを行うため、マイコン側のの負荷軽減してくれます。  
ArduinoやMicroPythonで、このデバイスをコントロールするデバイスドライバはたくさんあったのですが、tinygoで書かれたものが見当たらなかったので、自作してみました。  

