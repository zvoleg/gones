package ppu

import "fmt"

type objectAttributeEntity struct {
	yCoordinate byte
	tileId      byte
	attributes  byte
	xCoordinate byte
}

func (o *objectAttributeEntity) VerticalFlip() bool {
	return o.attributes&0x80 != 0
}

func (o *objectAttributeEntity) HorizontalFlip() bool {
	return o.attributes&0x40 != 0
}

func (o *objectAttributeEntity) BehindBackground() bool {
	return o.attributes&0x20 != 0
}

func (o *objectAttributeEntity) Palette() uint16 {
	return uint16(o.attributes&0x03) + 4 // sprite palette address starts after background palette
}

// only for 8x16 sprite size
func (o *objectAttributeEntity) SpriteTableAddress() uint16 {
	if o.tileId&0x01 == 0 {
		return 0x0000
	} else {
		return 0x1000
	}
}

func (o *objectAttributeEntity) ToString() string {
	return fmt.Sprintf("x: %03d y: %03d | tileId: %02X | attributes: %02X", o.xCoordinate, o.yCoordinate, o.tileId, o.attributes)
}
