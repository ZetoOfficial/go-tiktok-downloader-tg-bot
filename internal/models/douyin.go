package models

// DownloadResponse содержит результат загрузки.
// Если скачан единичный файл (например, видео), то поле Data будет заполнено.
// Если скачан zip-архив с фотографиями, то Photos будет содержать распакованные файлы.
type DownloadResponse struct {
	FileName string
	Data     []byte            // содержимое файла (например, видео)
	Photos   map[string][]byte // для zip-архива – распакованные фотографии
}

type APIError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Router  string            `json:"router,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
}

// DownloadOption представляет возможные параметры для загрузки.
type DownloadOption struct {
	Key   string
	Value bool
}

var (
	WithPrefix       = DownloadOption{"prefix", true}
	WithoutPrefix    = DownloadOption{"prefix", false}
	WithWatermark    = DownloadOption{"with_watermark", true}
	WithoutWatermark = DownloadOption{"with_watermark", false}
)
