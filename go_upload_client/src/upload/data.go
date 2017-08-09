package upload

type TGetBlockReply struct {
	RetCode      string `json:"retCode"`
	RetMsg       string `json:"retMsg"`
	TotalSize    int64  `json:"total_size"`
	TransferSize int64  `json:"transfer_size"`
	StartOffset  int64  `json:"start_offset"`
	EndOffset    int64  `json:"end_offset"`
	Token        string `json:"token"`
}

type TApplySessionRequest struct {
	FileSize int64  `json:"fileSize"`
	Md5      string `json:"md5"`
	SavePath string `json:"savePath"`
	FileName string `json:"fileName"`
	TaskId   string `json:"taskID"`
	Token    string `json:"token"`
}

type TApplySessionIdReply struct {
	RetCode   string `json:"retCode"`
	RetMsg    string `json:"retMsg"`
	SessionId string `json:"session_id"`
	UploadUrl string `json:"upload_url"`
	StateUrl  string `json:"state_url"`
	MaxThread string `json:"max_upload_thread"`
}

type TUpdateStateRequest struct {
	SessionId string `json: "sessionId"`
	State     string `json: "state"`
}

type TUploadBlockFormData struct {
	PartMd5     string `json:"part_md5"`
	Token       string `json:"token"`
	StartOffset int64  `json:"start_byte"`
	EndOffset   int64  `json:"end_byte"`
}
