package upload

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

func PostDataToServer(task_id int64, url string, buf []byte, reqData TUploadBlockFormData) ([]byte, error) {
	//	startPos, endPos := reqData.StartOffset, reqData.EndOffset+1
	//	blockSize := endPos - startPos

	//make body buffer
	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)

	//make file header
	fileHead := make(textproto.MIMEHeader)
	fileHead.Set("Content-Type", "multipart/form-data")
	fileHead.Set("Content-Disposition", "form-data; name=\"file\"; filename=\"unknow\"")
	fileWriter, _ := bodyWriter.CreatePart(fileHead)

	//write post data and cal MD5
	var strCalMd5 string
	if false {
		calcMD5 := md5.New()
		multiWriters := io.MultiWriter(calcMD5, fileWriter)
		multiWriters.Write(buf[:])
		strCalMd5 = fmt.Sprintf("%x", calcMD5.Sum(nil))
	} else {
		md5 := md5.Sum(buf)
		strCalMd5 = hex.EncodeToString(md5[:])
		fileWriter.Write(buf[:])
	}
	reqData.PartMd5 = strCalMd5

	//start make field data
	bodyWriter.WriteField("part_md5", reqData.PartMd5)
	bodyWriter.WriteField("start_byte", strconv.FormatInt(reqData.StartOffset, 10))
	bodyWriter.WriteField("end_byte", strconv.FormatInt(reqData.EndOffset, 10))
	bodyWriter.WriteField("token", reqData.Token)

	//close writer
	bodyWriter.Close()

	//make request
	req, _ := http.NewRequest("POST", url, bodyBuf)

	//request config
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	//start request
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Request error : ", reqData.Token)
		log.Println(err)
		return nil, err
	}
	result, err1 := ioutil.ReadAll(resp.Body)
	return result, err1
}

func GetJsonFromServer(taskId string, url string) ([]byte, error) {
	rep, err := http.Get(url)
	if err != nil {
		log.Printf("GetJsonFromServer : Get error --- ", taskId)
		log.Println(url)
		return nil, err
	}
	return ioutil.ReadAll(rep.Body)
}

func PostJsonToServer(taskId string, url string, strJson []byte) ([]byte, error) {
	rep, err := http.Post(url, "application/json", bytes.NewBuffer(strJson))
	if err != nil {
		log.Printf("PostJsonToServer : Post error --- ", taskId)
		log.Println(url)
		return nil, err
	}
	return ioutil.ReadAll(rep.Body)
}
