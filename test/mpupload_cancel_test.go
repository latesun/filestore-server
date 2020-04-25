package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	cfg "filestore-server/config"
	fsUtil "filestore-server/util"

	jsonit "github.com/json-iterator/go"
)

// 登录
func signin(username, password string) (token string, err error) {
	resp, err := http.PostForm(
		"http://localhost:8080/user/signin",
		url.Values{
			"username": {username},
			"password": {password},
		})
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(body) == "FAILED" {
		return "", err
	}
	token = jsonit.Get(body, "data").Get("Token").ToString()
	return token, nil
}

// 文件按指定大小分块，分段上传
func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		fmt.Printf("完成传输块index: %d\n", <-ch)
	}

	return nil
}

func TestMpUploadCancel(t *testing.T) {
	// 用户名/密码
	username := "test2020"
	password := "test2020"
	uploadEntry := cfg.UploadLBHost
	// 需要上传的文件名及文件hash
	fpath := "/data/pkg/armory.zip"
	fhash, err := fsUtil.ComputeSha1ByShell(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize, err := fsUtil.ComputeFileSizeByShell(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 0: 登录，获取token
	token, err := signin(username, password)
	if err != nil {
		fmt.Println(err)
		return
	} else if token == "" {
		fmt.Println("登录失败，请检查用户名、密码")
		return
	}

	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		fmt.Sprintf("%s/file/mpupload/init", uploadEntry),
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {fhash},
			"filesize": {strconv.Itoa(filesize)},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	tURL := fmt.Sprintf("%s/file/mpupload/uppart?username=%s&token=%s&uploadid=%s",
		uploadEntry, username, token, uploadID)
	multipartUpload(fpath, tURL, chunkSize)

	// 4. 取消分块上传接口
	resp, err = http.PostForm(
		fmt.Sprintf("%s/file/mpupload/cancel", uploadEntry),
		url.Values{
			"username": {username},
			"token":    {token},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	// 5. 打印分块上传结果
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("cancel upload result: %s\n", string(body))
}
