package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
	33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
	61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
	36, 20, 34, 44, 52,
}

func getMixinKey(orig string) string {
	var str strings.Builder
	for _, v := range mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
		}
	}
	return str.String()[:32]
}
func EncWbi(params map[string]string, imgKey string, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	params["wts"] = currTime
	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Remove unwanted characters
	for k, v := range params {
		v = strings.ReplaceAll(v, "!", "")
		v = strings.ReplaceAll(v, "'", "")
		v = strings.ReplaceAll(v, "(", "")
		v = strings.ReplaceAll(v, ")", "")
		v = strings.ReplaceAll(v, "*", "")
		params[k] = v
	}
	// Build URL parameters
	var str strings.Builder
	for _, k := range keys {
		str.WriteString(fmt.Sprintf("%s=%s&", k, params[k]))
	}
	query := strings.TrimSuffix(str.String(), "&")
	// Calculate w_rid
	hash := md5.Sum([]byte(query + mixinKey))
	params["w_rid"] = hex.EncodeToString(hash[:])
	return params
}

var cache sync.Map
var lastUpdateTime time.Time

func updateCache() {
	if time.Now().Sub(lastUpdateTime).Minutes() < 10 {
		return
	}
	imgKey, subKey := GetWbiKeys()
	cache.Store("imgKey", imgKey)
	cache.Store("subKey", subKey)
	lastUpdateTime = time.Now()
}

func GetWbiKeysCached() (string, string) {
	updateCache()
	imgKeyI, _ := cache.Load("imgKey")
	subKeyI, _ := cache.Load("subKey")
	return imgKeyI.(string), subKeyI.(string)
}

func GetWbiKeys() (string, string) {
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		fmt.Println("Error:", err)
		return "", ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", ""
	}
	json := string(body)
	imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
	subURL := gjson.Get(json, "data.wbi_img.sub_url").String()
	imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]
	return imgKey, subKey
}

// 签名
func SignURL(urlStr string) string {
	urlObj, _ := url.Parse(urlStr)
	imgKey, subKey := GetWbiKeysCached()
	//fmt.Println(imgKey, subKey)
	query := urlObj.Query()
	params := map[string]string{}
	for k, v := range query {
		params[k] = v[0]
	}
	newParams := EncWbi(params, imgKey, subKey)
	for k, v := range newParams {
		query.Set(k, v)
	}
	urlObj.RawQuery = query.Encode()
	newUrlStr := urlObj.String()
	return newUrlStr
}
