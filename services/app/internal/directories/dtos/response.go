package dtos

type DirectoryResponseDto struct {
	Id       string `json:"id"`
	UserId   string `json:"userId"`
	ParentId string `json:"parentId"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type DirectoryItemResponseDto struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Type   string `json:"type"`
	Name   string `json:"name"`
}
