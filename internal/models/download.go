package models

// DownloadResponse содержит результат загрузки.
// Если скачан единичный файл (например, видео), то поле Data будет заполнено.
// Если скачан набор фотографий (слайдшоу), то Photos будет содержать файлы.
type DownloadResponse struct {
	FileName string
	Data     []byte            // содержимое файла (например, видео)
	Photos   map[string][]byte // распакованные фотографии (слайдшоу)
}

// DownloadOption представляет возможные параметры для загрузки.
type DownloadOption struct {
	Key   string
	Value bool
}
