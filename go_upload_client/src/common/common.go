package common

const (
	Error_None                     = "0000"
	Error_Illeage_Request          = "0001"
	Error_Task_Not_Found           = "0002"
	Error_Verify_Md5_Failed        = "0003"
	Error_Token_Illeage            = "0004"
	Error_Task_State_Not_Equeal    = "0005"
	Error_Request_Data_Not_Correct = "0006"
	Error_Task_Id_Exist            = "0007"
	Error_Upload_Timeout           = "0008"
	Error_Unknow                   = "9999"
)

const (
	Msg_None                     = "success"
	Msg_Illeage_Request          = "request Illeage"
	Msg_Task_Not_Found           = "task not exsit"
	Msg_Verify_Md5_Failed        = "verify md5 failed"
	Msg_Token_Illeage            = "token Illeage"
	Msg_Task_State_Not_Queal     = "task state not equal uploding"
	Msg_Request_Data_Not_Correct = "request data not correct"
	Msg_Task_Id_Exist            = "task Id is already exsit"
	Msg_Upload_Timeout           = "block upload time out"
	Msg_Unknow                   = "unknow error"
)

const (
	Task_Has_Been_Finished  = "1001"
	No_Avaible_Block_Exist  = "1002"
	Running_Task_Over_Limit = "1003"
)
const (
	Block_State_Busy   = "server is busy"
	Block_State_Finish = "task is finished"
	Block_State_Ready  = "server is ready"
)

const (
	State_Quene     = "quene"
	State_Pause     = "pause"
	State_Uploading = "uploading"
	State_Cancle    = "cancle"
	State_Finished  = "finished"
)

type TResult struct {
	RetCode string `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}
