package lcd1602

import "time"
import "github.com/ftl/i2c"

const LCD_ADDR = 0x0
const LCD_BASE_ADDR = 0x27

type LCD struct {
	bus     *i2c.I2C
	lightOn bool
}

func NewLCD() (*LCD, error) {
	i2c, err := i2c.Open(LCD_BASE_ADDR, 1)
	if err != nil {
		return nil, err
	}
	return &LCD{
		i2c, true,
	}, nil

}

func (i *LCD) writeWord(addr, data uint8) {
	temp := data
	if i.lightOn {
		temp |= 0x08
	} else {
		temp &= 0xF7
	}
	i.bus.WriteReg(addr, temp)
}

func (i *LCD) sendCommand(comm uint8) {
	// Send bit7-4 firstly
	buf := comm & 0xF0
	buf |= 0x04 // RS = 0, RW = 0, EN = 1
	i.writeWord(LCD_ADDR, buf)
	time.Sleep(time.Millisecond * 2)

	buf &= 0xFB // Make EN = 0
	i.writeWord(LCD_ADDR, buf)

	// Send bit3-0 secondly
	buf = (comm & 0x0F) << 4
	buf |= 0x04 // RS = 0, RW = 0, EN = 1
	i.writeWord(LCD_ADDR, buf)
	time.Sleep(time.Millisecond * 2)
	buf &= 0xFB // Make EN = 0
	i.writeWord(LCD_ADDR, buf)
}

func (i *LCD) sendData(data uint8) {
	// Send bit7-4 firstly
	buf := data & 0xF0
	buf |= 0x05 // RS = 1, RW = 0, EN = 1
	i.writeWord(LCD_ADDR, buf)
	time.Sleep(time.Millisecond * 2)
	buf &= 0xFB // Make EN = 0
	i.writeWord(LCD_ADDR, buf)

	// Send bit3-0 secondly
	buf = (data & 0x0F) << 4
	buf |= 0x05 // RS = 1, RW = 0, EN = 1
	i.writeWord(LCD_ADDR, buf)
	time.Sleep(time.Millisecond * 2)
	buf &= 0xFB // Make EN = 0
	i.writeWord(LCD_ADDR, buf)
}

func (i *LCD) TurnLight(on bool) {
	i.lightOn = on
	if i.lightOn {
		i.bus.WriteReg(LCD_ADDR, 0x08)
	} else {
		i.bus.WriteReg(LCD_ADDR, 0x00)
	}
}

func (i *LCD) Clear() {
	i.sendCommand(0x01)
}
func (i *LCD) PrintText(x, y uint8, text string) {
	if y > 1 {
		y = 1
	}
	if x > 15 {
		x = 15
	}

	addr := 0x80 + 0x40*y + x
	i.sendCommand(addr)

	for _, char := range []byte(text) {
		i.sendData(byte(char))
	}
}
