package ppu

type planeDataBuffer struct {
	taileId       byte
	attributeData byte
	tileDataLow   byte
	tileDataHigh  byte
}

func (buff *planeDataBuffer) setTileId(data byte) {
	buff.taileId = data
}

func (buff *planeDataBuffer) setAttributeData(data byte) {
	buff.attributeData = data
}

func (buff *planeDataBuffer) setTileDataLow(data byte) {
	buff.tileDataLow = data
}

func (buff *planeDataBuffer) setTileDataHigh(data byte) {
	buff.tileDataHigh = data
}
