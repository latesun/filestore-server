package test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var (
	username  = "admin"
	token     = "54eefa7dbd5bcf852c52fecd816f2a315c61832c"
	targetURL = "http://localhost:28080/file/fastupload"
	filehash  = "no_such_file_hash"
	filename  = "just_for_test"
)

func TestUpload(t *testing.T) {
	resp, err := http.PostForm(targetURL, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filename": {filename},
	})

	t.Logf("error: %+v\n", err)
	t.Logf("resp: %+v\n", resp)
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		t.Logf("parseBodyErr: %+v\n", err)
		if err == nil {
			t.Logf("parseBody: %+v\n", string(body))
		}
	}
}
