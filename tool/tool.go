package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"unicode"
)

func interface2String(inter interface{}) string {

	switch inter.(type) {

	case string:
		// rt.Println("string", inter.(string))
		return inter.(string)
	case int:
		fmt.Println("int", inter.(int))
		return ""
	case float64:
		fmt.Println("float64", inter.(float64))
		return ""
	}
	return ""

}

func Isnumber(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func String2Int(str5 string) int {
	int5, err := strconv.Atoi(str5)
	if err != nil {
		fmt.Println(err)
		return -1 //error
	} else {

		return int5
	}
}

func MapToJson(param map[string]map[string]string/*interface{}*/) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func InterfaceToJson(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}


func JsonToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}


func URLCode(yoururl string) string{
	return url.QueryEscape(yoururl)
}
func UnURLCode(yoururl string) string{
	decodeurl,err := url.QueryUnescape(yoururl)
	if err != nil {
		fmt.Println(err)
	}
	return decodeurl
}


func GetFileMd5(filename string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("os Open error")
        return "", err
    }
    md5 := md5.New()
    _, err = io.Copy(md5, file)
    if err != nil {
        fmt.Println("io copy error")
        return "", err
    }
    md5Str := hex.EncodeToString(md5.Sum(nil))
    return md5Str, nil
}
  
func GetStringMd5(s string) string {
    md5 := md5.New()
    md5.Write([]byte(s))
    md5Str := hex.EncodeToString(md5.Sum(nil))
    return md5Str
}