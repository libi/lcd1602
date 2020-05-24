package lcd1602

import "testing"

func TestLCD_PrintText(t *testing.T) {
	lcd, _ := NewLCD()
	lcd.PrintText(0, 0, "hello go!")
}
