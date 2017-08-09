package serves

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"uploadTest/src/common"
	"uploadTest/src/upload"
	"uploadTest/src/view"

	//"github.com/khlipeng/websocket"
)

func (this *TRouters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/upload" {
		DisplayWebForm(w, r)
	} else if r.URL.Path == "/webSocket" {
		//fmt.Println("uploadFilesWebSocket")
		//		routers.handleSoc.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/jsRequest") {
		//fmt.Println("uploadjsRequest")
		uploadByJsRequest(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static") {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	} else {
		log.Print(r.URL.Path)
		badRequest(w, r)
	}
}

func DisplayWebForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := views.GetHtml(2)
		if err != nil {
			log.Println(err)
			return
		}
		token := GetRandomSalt()
		log.Println(token)
		err1 := t.Execute(w, token)
		if err1 != nil {
			log.Print(err1)
		}
	}
}

func uploadByJsRequest(w http.ResponseWriter, r *http.Request) {
	//解析请求数据
	r.ParseForm()
	if r.Method == "GET" {
		token := strings.Replace(r.URL.Path, "/jsRequest", "", -1)
		reply := ReplyToClient(token)
		//log.Print("send to client: " + string(reply))
		fmt.Fprintf(w, string(reply))

	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		//log.Print("receive from client: " + string(result))
		_, err := ReceiveFromClient(string(result))
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

func ReceiveFromClient(result string) (string, error) {
	var fileTask tokenTaskRequest
	result = strings.Replace(result, "\\", "\\\\", -1)
	json.Unmarshal([]byte(result), &fileTask)
	log.Print("recv from client parse json : ", fileTask)
	if fileTask.Token != "" {
		if taskMap, ok := tokenMap[fileTask.Token]; !ok {
			taskMap = make(map[string]*upload.TUploadTask)
			uploadFiles(fileTask.FilePath, taskMap)
			tokenMap[fileTask.Token] = taskMap
			return fileTask.Token, nil
		} else {
			log.Print("the token map is exit!")
			if len(tokenMap[fileTask.Token]) == 0 {
				return fileTask.Token, nil
			} else {
				return fileTask.Token, errors.New("the last tasks is not complete!")
			}
		}
	}
	return fileTask.Token, nil
}

func ReplyToClient(token string) string {
	if token != "" {
		var index = 100
		var nowTime = time.Now()
		var dataArray []tokenTaskReturn
		for taskToken, existTask := range tokenMap[token] {
			var data tokenTaskReturn
			data.FileName = existTask.FileName
			data.FileSize = GetFormatSize(existTask.FileSize)
			data.State = existTask.State
			data.TransferSize = GetFormatSize(existTask.TransferBytes)
			data.Rate = GetFormatSize(existTask.Rate) + "/s"
			//data.UsedTime = strconv.FormatFloat(nowTime.Sub(existTask.StartTime).Seconds(), 'f', 1, 64)
			data.UsedTime = GetTimeString(nowTime.Sub(existTask.StartTime))
			//data.ID = existTask.TaskID
			data.ID = index
			data.Progress = fmt.Sprintf("%d", existTask.TransferBytes*100/existTask.FileSize)
			if existTask.State == common.State_Finished {
				data.Progress = "100"
				data.Rate = "0B/s"
				delete(tokenMap[token], taskToken)
			}
			dataArray = append(dataArray, data)
			index++
		}
		if len(tokenMap[token]) == 0 {
			delete(tokenMap, token)
		}
		reply, _ := json.Marshal(dataArray)
		return string(reply)
	} else {
		log.Print("the token is empty")
		return "null"
	}
}

/*
func SocketLoop(ws *websocket.Conn) {
	var strRequest string
	if err := websocket.Message.Receive(ws, &strRequest); err != nil {
		log.Println("Can't receive")
		log.Println(err)
		return
	}
	token, _ := ReceiveFromClient(strRequest)
	SendToClientLoop(token, ws)
}

func SendToClientLoop(token string, ws *websocket.Conn) {
loop:
	for {
		reply := ReplyToClient(token)
		//log.Print("send to client: " + string(reply))
		if err := websocket.Message.Send(ws, string(reply)); err != nil {
			log.Println(err)
			break loop
		}
		time.Sleep(500000000)
	}
}
*/
func badRequest(w http.ResponseWriter, r *http.Request) {
	var result common.TResult
	result.RetCode = common.Error_Illeage_Request
	result.RetMsg = common.Msg_Illeage_Request

	data, err := json.Marshal(result)
	if err == nil {
		fmt.Fprint(w, string(data))
	} else {
		log.Print(err)
		fmt.Fprint(w, err)
	}
}
