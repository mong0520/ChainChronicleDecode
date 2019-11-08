package tools

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

var wg sync.WaitGroup
var debug = false

// DecodeByPath decodes all scr by given `path` and save the resized jpg file in the same folder
func DecodeByPath(path string) {
	var scrSlice = []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		if strings.HasSuffix(f.Name(), "scr") {
			scrSlice = append(scrSlice, filepath.Join(path, f.Name()))
		}
	}

	for _, scr := range scrSlice {
		outputFullName := strings.ReplaceAll(scr, "scr", "jpg")
		fmt.Println(outputFullName)
		dat, err := ioutil.ReadFile(scr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		output := DecodeRotData(dat)
		newOutput := ResizeImage(output, 960, 1440)

		ioutil.WriteFile(outputFullName, newOutput, 0644)
	}

	fmt.Println(scrSlice)
}

// DecodeRotData decodes *.scr data into image bytes
func DecodeRotData(from []byte) []byte {
	num := int(from[0])
	tmpVal := 0
	if num == 0 {
		tmpVal = 4
	} else {
		tmpVal = num
	}

	num2 := len(from) - (4 + (4 - tmpVal))
	buffer := make([]byte, num2)
	v18 := int(from[0])
	v4 := num2 / 4

	if int(from[0]) == 0 {
		v4--
		v18 = 4
	}
	v6 := int(from[1])
	v7 := int(from[2])
	v9 := 0
	v10 := 0
	for v8 := 0; v8 <= v4; v8++ {
		v13 := make([]byte, 4)
		v13[0] = from[v9+7]
		v13[1] = from[v9+6]
		v13[2] = from[v9+5]
		v13[3] = from[v9+4]

		// to be test
		v14 := binary.LittleEndian.Uint32(v13)

		// int ror = 32 - (v8 + v6) % v7;
		ror := uint32(32 - (v8+v6)%v7)
		v14 = (v14 >> ror) | (v14 << (32 - ror))
		v15 := getBytes(v14)
		v16 := 4

		if v8 == v4 {
			v16 = v18
		}
		for i := v16 - 1; i >= 0; i-- {
			buffer[v10] = v15[i]
			v10++
		}
		v9 += 4
	}
	return buffer
}

func ResizeImage(data []byte, width uint, height uint) []byte {
	image, _, _ := image.Decode(bytes.NewReader(data))
	// check err
	newImage := resize.Resize(width, height, image, resize.Lanczos3)

	// Encode uses a Writer, use a Buffer if you need the raw []byte
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := jpeg.Encode(w, newImage, nil)
	if err != nil {
		fmt.Println(err)
	}
	return b.Bytes()
}

func simpleGet(url string) (data []byte, err error) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "content.cc.mobimon.com.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	if res.StatusCode != 200 {
		return data, errors.New("not a valid url")
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func StartConvert(contentRoot string, d bool, concurrence int, outputFolder string) {
	debug = d

	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	scrList, _ := GetCardList(contentRoot)
	if debug {
		fmt.Println(scrList)
	}
	msg := make(chan string)

	// produce msg
	go produceScr(msg, scrList)

	wg.Add(concurrence)
	for idx := 0; idx < concurrence; idx++ {
		go consumeScr(msg, idx, outputFolder)
	}
	wg.Wait()
	fmt.Println("Comlete!")
}

func produceScr(msg chan string, scrList []string) {
	defer close(msg)

	for _, scr := range scrList {
		msg <- scr
	}
}

func consumeScr(msg chan string, workerID int, folderToSave string) {
	for scr := range msg {
		baseName := strings.ReplaceAll(path.Base(scr), "scr", "jpg")
		filePath := filepath.Join(folderToSave, baseName)
		if debug {
			fmt.Printf("Worker[%d] start to cunsume: [%s] and saved to [%s]\n", workerID, scr, filePath)
		} else {
			fmt.Printf("Saving [%s]...\n", filePath)
		}
		convertToJpgFileByScrUrl(scr, filePath)
	}
	wg.Done()
}

func convertToJpgFileByScrUrl(scrUrl string, fPath string) {
	data, err := simpleGet(scrUrl)
	if err != nil {
		if debug {
			fmt.Println(err)
		}
	} else {
		jppData := DecodeRotData(data)
		finalJpgData := ResizeImage(jppData, 960, 1440)
		ioutil.WriteFile(fPath, finalJpgData, 0644)
	}
}

func GetCardList(contentRoot string) (scrList []string, err error) {
	// url := "http://content.cc.mobimon.com.tw/382/Prod/Bdl54_iOS/files.json?cnt=16e48d9c78a&timestamp=1573180257282"
	if strings.HasSuffix(contentRoot, "/") {
		contentRoot = contentRoot[0 : len(contentRoot)-1]
	}
	url := contentRoot + "/Bdl54_iOS/files.json"

	// fmt.Println(res)
	// fmt.Println(string(body))
	body, err := simpleGet(url)

	rawData := map[string]interface{}{}
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		return scrList, err
	}
	files := rawData["files"].(map[string]interface{})
	cardsOrgSize := files["Card_OrgSize"].(map[string]interface{})
	for k := range cardsOrgSize {
		scrName := strings.Replace(k, "bdl", "scr", -1)
		scrURL := contentRoot + "/Resource/Card/" + scrName
		scrList = append(scrList, scrURL)
	}

	return scrList, nil
}

func getBytes(x uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}
