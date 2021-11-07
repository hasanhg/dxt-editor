package dxt

import (
	"image"
	"log"
	"strconv"
	"strings"

	"dxt-editor/utils"
)

type DXTFile struct {
	Compression string
	Header      string
	Image       *image.RGBA
}

type Buffer struct {
	data   []byte
	offset uint64
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{data: data, offset: 0}
}

func (b *Buffer) ParseDXT() *DXTFile {
	dxtFile := &DXTFile{}
	headerSize := b.Read(UINT32).(uint32)
	header := utils.BytesToString(b.ReadN(uint64(headerSize)))
	dxtFile.Header = header

	b.ReadN(3)   // NTF
	b.Read(BYTE) // LENGTH
	width := b.Read(UINT16).(uint16)
	x := b.Read(UINT16).(uint16)
	height := b.Read(UINT16).(uint16)
	y := b.Read(UINT16).(uint16)
	dxtFile.Compression = utils.BytesToString(b.ReadN(4))
	b.Read(UINT32) // SOMETHING

	dxtFile.Image = image.NewRGBA(image.Rect(int(x), int(y), int(width), int(height)))

	switch strings.ToUpper(dxtFile.Compression) {
	case "DXT1":
		parseDXT1(b, dxtFile.Image)
	case "DXT2":
		parseDXT2(b, dxtFile.Image)
	case "DXT3":
		parseDXT3(b, dxtFile.Image)
	case "DXT4":
		parseDXT4(b, dxtFile.Image)
	case "DXT5":
		parseDXT5(b, dxtFile.Image)
	}

	return dxtFile
}

func (b *Buffer) Read(t ColType) interface{} {

	data := b.ReadN(typeSizes[t])
	switch t {
	case BYTE:
		return data[0]

	case INT16:
		return int16(utils.BytesToInt(data, true))

	case UINT16:
		return uint16(utils.BytesToInt(data, true))

	case INT32:
		return int32(utils.BytesToInt(data, true))

	case UINT32:
		return uint32(utils.BytesToInt(data, true))

	case INT64:
		return int64(utils.BytesToInt(data, true))

	case UINT64:
		return uint64(utils.BytesToInt(data, true))

	case FLOAT32:
		return float32(utils.BytesToFloat(data, true))

	default:
		return nil
	}
}

func (b *Buffer) ReadN(n uint64) []byte {

	if b.offset == uint64(len(b.data)) {
		return []byte{}
	}

	if uint64(len(b.data)) >= b.offset+n {
		b.offset += n
		return b.data[b.offset-n : b.offset]
	}

	data := b.data[b.offset:]
	b.offset = uint64(len(b.data))
	return data
}

func (b *Buffer) Write(value string, t ColType) {

	if t == STRING {
		data := value
		length := uint64(len(data))
		b.data = append(b.data, utils.IntToBytes(length, 4, true)...)
		b.offset += 4

		b.data = append(b.data, []byte(data)...)
		b.offset += length
		return
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}

	switch t {
	case BYTE:
		data := byte(val)
		b.data = append(b.data, data)
		b.offset++

	case INT16:
		data := utils.IntToBytes(uint64(int16(val)), 2, true)
		b.data = append(b.data, data...)
		b.offset += 2

	case UINT16:
		data := utils.IntToBytes(uint64(uint16(val)), 2, true)
		b.data = append(b.data, data...)
		b.offset += 2

	case INT32:
		data := utils.IntToBytes(uint64(int32(val)), 4, true)
		b.data = append(b.data, data...)
		b.offset += 4

	case UINT32:
		data := utils.IntToBytes(uint64(uint32(val)), 4, true)
		b.data = append(b.data, data...)
		b.offset += 4

	case FLOAT32:
		data := utils.FloatToBytes(float64(val), 4, true)
		b.data = append(b.data, data...)
		b.offset += 4
	}
}

func (b *Buffer) Overwrite(data []byte, i int) {
	_data := make([]byte, len(data))
	copy(_data, data)
	(*b).data = append(((*b).data)[:i], append(_data, ((*b).data)[i+len(data):]...)...)
}

func (b *Buffer) GetBytes() []byte {
	return b.data
}

func (b *Buffer) GetOffset() uint64 {
	return b.offset
}
