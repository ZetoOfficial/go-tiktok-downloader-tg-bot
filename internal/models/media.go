package models

type MediaInput struct {
	Type     string
	FileName string
	Data     []byte
	Caption  string
	Video    VideoMeta // размеры и длительность для Type == "video"
}

// VideoMeta описывает видео для Telegram. Без ширины и высоты клиенты
// Telegram показывают ролик квадратным плеером. Нулевые значения не передаются.
type VideoMeta struct {
	Width    int
	Height   int
	Duration int // секунды
}
