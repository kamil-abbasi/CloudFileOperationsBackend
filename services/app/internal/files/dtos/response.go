package dtos

type FileResponseDto struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	DirectoryId string `json:"directoryId"`
	Name        string `json:"name"`
	Size        uint64 `json:"size"`
	Location    string `json:"location"`
	Checksum    string `json:"checksum"`
}
