package request

type ShareRequest struct {
	ExpiredTime string `json:"expiredTime"`
	File        string `json:"file"`
}

type ShareidRequest struct {
	ShareId     string `json:"shareid"`
	SharedToken string `json:"sharedToken"`
}
