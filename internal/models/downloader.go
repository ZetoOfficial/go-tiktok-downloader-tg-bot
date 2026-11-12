package models

type Media struct {
	VideoData []byte
	VideoName string
	VideoMeta VideoMeta
	Photos    map[string][]byte
	Items     []MediaInput
}
