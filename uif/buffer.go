package uif

import (
	"image"
	"log"
	"strconv"
	"strings"

	"dxt-editor/utils"

	"github.com/google/uuid"
)

type Object struct {
	ID         string
	Name       string
	Rect       image.Rectangle
	MovRect    image.Rectangle
	Style      uint32
	Reserved   uint32
	Tooltip    string
	SoundOpen  string
	SoundClose string
	Tail       []byte
	Type       ObjectType
	Children   []*Object
	Visible    bool

	// Image properties
	Texture        string
	Crop           FRectangle
	AnimationFrame float32

	// Area properties
	// AreaType uint32
}

type UIFFile struct {
	Root *Object
}

type Buffer struct {
	data   []byte
	offset uint64
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{data: data, offset: 0}
}

func (b *Buffer) ParseUIF() *UIFFile {
	uifFile := &UIFFile{}
	uifFile.Root = b.readObj(uifFile, OT_BASE)
	return uifFile
}

func (b *Buffer) readObj(uifFile *UIFFile, otype ObjectType) *Object {
	obj := &Object{Visible: true}

	nameSize := b.Read(UINT32).(uint32)
	obj.Name = utils.BytesToString(b.ReadN(uint64(nameSize)))
	obj.Type = otype

	manifest := b.Read(UINT32).(uint32)
	childCount := manifest & 0xFF
	tail := manifest >> 16

	for i := 0; i < int(childCount); i++ {
		childType := ObjectType(b.Read(UINT32).(uint32))
		obj.Children = append(obj.Children, b.readObj(uifFile, childType))
	}

	idSize := b.Read(UINT32).(uint32)
	obj.ID = utils.BytesToString(b.ReadN(uint64(idSize)))
	if obj.ID == "" {
		obj.ID = strings.Split(uuid.New().String(), "-")[0]
	}

	obj.Rect.Min.X = int(b.Read(UINT32).(uint32))
	obj.Rect.Min.Y = int(b.Read(UINT32).(uint32))
	obj.Rect.Max.X = int(b.Read(UINT32).(uint32))
	obj.Rect.Max.Y = int(b.Read(UINT32).(uint32))

	obj.MovRect.Min.X = int(b.Read(UINT32).(uint32))
	obj.MovRect.Min.Y = int(b.Read(UINT32).(uint32))
	obj.MovRect.Max.X = int(b.Read(UINT32).(uint32))
	obj.MovRect.Max.Y = int(b.Read(UINT32).(uint32))

	obj.Style = b.Read(UINT32).(uint32)
	obj.Reserved = b.Read(UINT32).(uint32)

	size := b.Read(UINT32).(uint32)
	obj.Tooltip = utils.BytesToString(b.ReadN(uint64(size)))

	size = b.Read(UINT32).(uint32)
	obj.SoundOpen = utils.BytesToString(b.ReadN(uint64(size)))

	size = b.Read(UINT32).(uint32)
	obj.SoundClose = utils.BytesToString(b.ReadN(uint64(size)))

	obj.Tail = b.ReadN(uint64(tail))

	switch obj.Type {
	case OT_BASE:
		break
	case OT_IMAGE:
		size = b.Read(UINT32).(uint32)
		obj.Texture = utils.BytesToString(b.ReadN(uint64(size)))

		obj.Crop.Min.X = b.Read(FLOAT32).(float32)
		obj.Crop.Min.Y = b.Read(FLOAT32).(float32)
		obj.Crop.Max.X = b.Read(FLOAT32).(float32)
		obj.Crop.Max.Y = b.Read(FLOAT32).(float32)

		obj.AnimationFrame = b.Read(FLOAT32).(float32)
	case OT_AREA:
		//obj.AreaType = b.Read(UINT32).(uint32)
	}

	return obj
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
