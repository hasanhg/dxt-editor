package uif

type ColType byte
type ObjectType byte

const (
	BYTE    ColType = 2
	INT16   ColType = 3
	UINT16  ColType = 4
	INT32   ColType = 5
	UINT32  ColType = 6
	STRING  ColType = 7
	FLOAT32 ColType = 8
	INT64   ColType = 9
	UINT64  ColType = 10

	OT_BASE        ObjectType = 0
	OT_BUTTON      ObjectType = 1
	OT_STATIC      ObjectType = 2
	OT_PROGRESS    ObjectType = 3
	OT_IMAGE       ObjectType = 4
	OT_SCROLLBAR   ObjectType = 5
	OT_STRING      ObjectType = 6
	OT_TRACKBAR    ObjectType = 7
	OT_EDIT        ObjectType = 8
	OT_AREA        ObjectType = 9
	OT_TOOLTIP     ObjectType = 10
	OT_ICON        ObjectType = 11
	OT_ICONMANAGER ObjectType = 12
	OT_ICONSLOT    ObjectType = 13
	OT_LIST        ObjectType = 14
	OT_UNK15       ObjectType = 15
	OT_UNK16       ObjectType = 16
	OT_UNK17       ObjectType = 17
	OT_FLASH       ObjectType = 18
)

type FPoint struct {
	X, Y float32
}

type FRectangle struct {
	Min, Max FPoint
}

var (
	ansiChars   = []rune{'€', '\x81', '‚', 'ƒ', '„', '…', '†', '‡', 'ˆ', '‰', 'Š', '‹', 'Œ', '\x8d', 'Ž', '\x8f', '\x90', '‘', '’', '“', '”', '•', '–', '—', '˜', '™', 'š', '›', 'œ', '\x9d', 'ž', 'Ÿ', '\xa0', '¡', '¢', '£', '¤', '¥', '¦', '§', '¨', '©', 'ª', '«', '¬', '\xad', '®', '¯', '°', '±', '²', '³', '´', 'µ', '¶', '·', '¸', '¹', 'º', '»', '¼', '½', '¾', '¿', 'À', 'Á', 'Â', 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë', 'Ì', 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á', 'â', 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í', 'î', 'ï', 'ð', 'ñ', 'ò', 'ó', 'ô', 'õ', 'ö', '÷', 'ø', 'ù', 'ú', 'û', 'ü', 'ý', 'þ', 'ÿ'}
	typeSizes   = map[ColType]uint64{BYTE: 1, INT16: 2, UINT16: 2, INT32: 4, UINT32: 4, FLOAT32: 4, INT64: 8, UINT64: 8}
	typeTitles  = map[ColType]string{BYTE: "Byte", INT16: "Int16", UINT16: "UInt16", INT32: "Int32", UINT32: "UInt32", FLOAT32: "Float32", STRING: "String", INT64: "Int64", UINT64: "UInt64"}
	rTypeTitles = make(map[string]ColType)
)

func init() {
	for typ, title := range typeTitles {
		rTypeTitles[title] = typ
	}
}

func (t *ObjectType) String() string {
	switch *t {
	case OT_BASE:
		return "Base"
	case OT_IMAGE:
		return "Image"
	case OT_AREA:
		return "Area"
	}
	return ""
}
