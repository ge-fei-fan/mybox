package request

type ShareRequest struct {
	ExpiredTime string `json:"expiredTime"`
	File        string `json:"file"`
}
