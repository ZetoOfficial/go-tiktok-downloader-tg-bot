package models

type Media struct {
	VideoData []byte
	VideoName string
	Photos    map[string][]byte
}
