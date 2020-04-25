package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	cfg "filestore-server/config"
	fsUtil "filestore-server/util"

	jsonit "github.com/json-iterator/go"
)

func TestMpUpload(t *testing.T) {
	// 用户名/密码
	username := "test2020"
	password := "test2020"
	uploadEntry := cfg.UploadLBHost
	// 需要上传的文件名及文件hash
	fpath := "/data/pkg/armory.zip"
	fhash, err := fsUtil.ComputeSha1ByShell(fpath)
	if err != nil {
		t.Log(err)
		return
	}

	filesize, err := fsUtil.ComputeFileSizeByShell(fpath)
	if err != nil {
		t.Log(err)
		return
	}

	// 0: 登录，获取token
	token, err := signin(username, password)
	if err != nil {
		t.Log(err)
		return
	} else if token == "" {
		t.Log("登录失败，请检查用户名、密码")
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
		t.Log(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err.Error())
		os.Exit(-1)
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	t.Logf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	tURL := fmt.Sprintf("%s/file/mpupload/uppart?username=%s&token=%s&uploadid=%s",
		uploadEntry, username, token, uploadID)
	if err = multipartUpload(fpath, tURL, chunkSize); err != nil {
		t.Log(err.Error())
		os.Exit(-1)
	}

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		fmt.Sprintf("%s/file/mpupload/complete", uploadEntry),
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {fhash},
			"filesize": {strconv.Itoa(filesize)},
			"filename": {filepath.Base(fpath)},
			"uploadid": {uploadID},
		})

	if err != nil {
		t.Log(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	// 5. 打印分块上传结果
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err.Error())
		os.Exit(-1)
	}
	t.Logf("complete result: %s\n", string(body))
}
