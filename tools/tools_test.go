package tools

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecodeRotDataAndResize(t *testing.T) {
	dat, err := ioutil.ReadFile("sample/cha_2d_card_26173.scr")
	if err != nil {
		t.Error(err)
	}

	output := DecodeRotData(dat)
	newOutput := ResizeImage(output, 960, 1440)

	ioutil.WriteFile("sample/output.jpg", newOutput, 0644)
}

func TestDecodeByPath(t *testing.T) {
	DecodeByPath("../sample")
}

func TestGetCardList(t *testing.T) {
	scrList, err := GetCardList("http://content.cc.mobimon.com.tw/382/Prod/")
	if err != nil {
		t.Error(err)
	}
	if len(scrList) == 0 {
		t.Error("list size should not be 0")
	}
	fmt.Println(scrList)
}

func TestStartConvert(t *testing.T) {
	contentRoot := "http://content.cc.mobimon.com.tw/382/Prod/"
	StartConvert(contentRoot, true, 10, "./images")
}

func TestConvertToJpgFileByScrUrl(t *testing.T) {
	convertToJpgFileByScrUrl("http://content.cc.mobimon.com.tw/382/Prod/Resource/Card/cha_2d_card_07554.scr", "temp.jpg")
}
