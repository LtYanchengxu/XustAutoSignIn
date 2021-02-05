package serializer

type ResultUtil struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	IsSuccess bool   `json:"isSuccess"`
}
