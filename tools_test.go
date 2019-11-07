package tools

import (
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
	DecodeByPath("sample")
}
