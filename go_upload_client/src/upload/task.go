package upload

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	_ "strings"
	_ "sync"
	"time"
	"uploadTest/src/common"
)

var (
	ApplySessionUrl string
)

type TUploadTask struct {
	ThreadNum     int64
	FileSize      int64
	TransferBytes int64
	Rate          int64
	MD5           string
	LocalPath     string
	FileName      string
	TargetPath    string
	SeesionID     string
	TaskID        string
	StateUrl      string
	UploadUrl     string
	State         string

	Timer        *time.Ticker //定时获取rate
	StartTime    time.Time    //记录开始时间
	File         *os.File     //文件句柄
	NeedAreaChan chan TUploadBlockFormData
	TransferChan chan int64
}

func init() {
	//ApplySessionUrl = "http://172.16.7.97:8081/session"
	ApplySessionUrl = "http://127.0.0.1:8080/session"
}

func (task *TUploadTask) Init() error {
	task.TransferBytes = 0
	task.ThreadNum = 0
	task.NeedAreaChan = make(chan TUploadBlockFormData, 3)
	task.TransferChan = make(chan int64, 3)
	//task.TaskID = strings.ToUpper(GetRandomSalt())

	fileInfo, err := os.Stat(task.LocalPath)
	if err != nil {
		log.Println(err)
		return err
	}
	task.FileSize = fileInfo.Size()
	task.FileName = fileInfo.Name()

	calcMD5 := md5.New()
	io.Copy(calcMD5, task.File)
	task.MD5 = fmt.Sprintf("%x", calcMD5.Sum(nil))
	return nil
}
func (task *TUploadTask) StartTask() {
	errInit := task.Init()
	if errInit != nil {
		return
	}
	var reqData TApplySessionRequest
	reqData.FileSize = task.FileSize
	reqData.Md5 = task.MD5
	reqData.SavePath = task.TargetPath
	reqData.TaskId = task.TaskID
	reqData.FileName = task.FileName

	err := task.ApplyTaskSession(reqData)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("thread num : %d \n", task.ThreadNum)

	var errOpen error
	task.File, errOpen = os.OpenFile(task.LocalPath, os.O_RDWR, 0666)
	if errOpen != nil {
		log.Print(errOpen)
		return
	}
	task.StartTime = time.Now()

	task.PostUploadState(common.State_Uploading)
	task.InitPostThread(task.ThreadNum)
	go task.StartRateTimer()
	go task.GetUploadBlock()
}

func (task *TUploadTask) ApplyTaskSession(reqData TApplySessionRequest) error {
	log.Println("ApplyTaskSession start")
	strReq, _ := json.Marshal(reqData)
	result, err := PostJsonToServer(task.TaskID, ApplySessionUrl, strReq)
	if err != nil {
		log.Printf("ApplyTaskSession error:")
		return err
	}
	var retData TApplySessionIdReply
	json.Unmarshal(result, &retData)
	if retData.RetCode != common.Error_None {
		log.Printf("ApplyTaskSession --- server return error:")
		return errors.New(string(result))
	}
	task.ThreadNum, _ = strconv.ParseInt(retData.MaxThread, 10, 64)
	task.SeesionID = retData.SessionId
	task.StateUrl = retData.StateUrl + task.MD5
	log.Println(task.StateUrl)
	log.Println(task.UploadUrl)
	task.UploadUrl = retData.UploadUrl + task.SeesionID + "_" + task.MD5
	log.Println("ApplyTaskSession finished")
	return nil
}

func (task *TUploadTask) GetUploadBlock() {
	for {
	ForLoop:
		result, err := GetJsonFromServer(task.TaskID, task.UploadUrl)
		if err != nil {
			log.Println("GetUploadBlock : get pos error")
			log.Println(err)
			return
		}
		var getData TGetBlockReply
		json.Unmarshal(result, &getData)
		if getData.RetCode == common.Task_Has_Been_Finished {
			log.Print(common.Block_State_Finish)
			task.State = common.State_Finished
			task.Quit()
			return
		} else if getData.RetCode != common.Error_None {
			time.Sleep(200000000)
			//log.Println(getData)
			goto ForLoop
		}
		var needData TUploadBlockFormData
		needData.StartOffset = getData.StartOffset
		needData.EndOffset = getData.EndOffset
		needData.Token = getData.Token

		task.NeedAreaChan <- needData
		//log.Println("get pos finished")
		//time.Sleep(500000000)
		//task.RequestChan <- 1
	}
}
func (task *TUploadTask) PostUploadState(state string) {

	//解析请求参数
	var reqData TUpdateStateRequest
	reqData.SessionId = task.SeesionID
	reqData.State = state
	strReq, _ := json.Marshal(reqData)
	//stateUrl := strings.Replace(task.UploadUrl, "/upload", "/state", -1)
	//log.Println(task.StateUrl)
	result, err := PostJsonToServer(task.TaskID, task.StateUrl, strReq)
	if err != nil {
		log.Printf("PostUploadState --- post error : ")
		log.Println(err)
		return
	}

	var getData TGetBlockReply
	json.Unmarshal(result, &getData)
	if getData.RetCode != common.Error_None {
		log.Printf("PostUploadState --- server return error : ")
		log.Println(getData.RetMsg)
		return
	}
	task.State = state
	//log.Println(resp.Status)
	//log.Println(resp.Header)
}
func (task *TUploadTask) PostUploadBlock(task_id int64) {

	for {
		//log.Printf("goroutine waiting : task id -- %d \n", task_id)
		//wait pos
		uploadReqData := <-task.NeedAreaChan
		//start post
		startPos, endPos := uploadReqData.StartOffset, uploadReqData.EndOffset+1
		if endPos <= startPos {
			log.Print("PostUploadBlock : startPos > endPos  ", startPos, endPos)
			continue
		}

		blockSize := endPos - startPos
		secReader := io.NewSectionReader(task.File, startPos, blockSize)

		//log.Printf("start post data -- %d %d %d --- %s \n", startPos, endPos, task_id, uploadReqData.Token)

		var readBuffer []byte
		var readLen int64
		var minSize int64
		readLen = 0
		minSize = 16 * 1024 * 10
		buf := make([]byte, minSize)

		for i := int64(0); i < blockSize; i += minSize {
			if i != 0 {
				secReader.Seek(i, 0)
			}
			rLen, err11 := secReader.Read(buf)
			//rLen, err11 := task.ReadFile(startPos+i, buf, task_id)
			readBuffer = append(readBuffer, buf[:]...)
			readLen += int64(rLen)
			if err11 == io.EOF {
				break
			} else if err11 != nil {
				log.Printf("read error : ")
				log.Println(err11)
				return
			}
			//fmt.Println(readLen, blockSize, uploadReqData.Token)
		}
		if readLen != blockSize {
			log.Printf("read error : readLen != blockSize", readLen, blockSize, task_id)
			continue
		}
		result, err := PostDataToServer(task_id, task.UploadUrl, readBuffer, uploadReqData)
		if err != nil {
			log.Printf("Post data error : ", uploadReqData.Token)
			log.Println(err)
			continue
		}

		var retResult common.TResult
		json.Unmarshal(result, &retResult)
		if retResult.RetCode != common.Error_None {
			log.Printf("PostUploadBlock --- retResult error : ", uploadReqData.Token)
			log.Println(retResult.RetMsg)
			continue
		}
		task.TransferBytes += blockSize
		task.TransferChan <- task.TransferBytes
		//log.Printf("end post data -- %d %d %d --- %s \n", startPos, endPos, task_id, uploadReqData.Token)
		//<-task.RequestChan
	}
}

func (task *TUploadTask) InitPostThread(num int64) {

	for i := int64(0); i < num; i++ {
		go task.PostUploadBlock(i)
	}
}

func (task *TUploadTask) Quit() {
	task.Timer.Stop()
	task.File.Close()
}

func (task *TUploadTask) StartRateTimer() {
	var lastTransfer int64
	var trnasferBytes int64
	lastTransfer = int64(0)
	task.Timer = time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-task.Timer.C:
			task.Rate = trnasferBytes - lastTransfer
			lastTransfer = trnasferBytes
			//log.Print("rec rate : ", task.Rate)
			if task.TransferBytes == task.FileSize {
				task.Timer.Stop()
				//log.Print("rec finished : ")
				return
			}
		case transferChan := <-task.TransferChan:
			//log.Print("rec bytes : ", transferChan)
			trnasferBytes = transferChan
		}
	}

}

//func (task *TUploadTask) ReadFile(pos int64, buf []byte, id int64) (n int, err error) {
//	task.Mutex.Lock()
//	task.File.Seek(pos, 0)
//	rLen, err := task.File.Read(buf)
//	task.Mutex.Unlock()
//	return rLen, err
//}
