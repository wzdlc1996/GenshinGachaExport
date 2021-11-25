package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/xuri/excelize/v2"
)

const gachaTypeAPI = "https://hk4e-api.mihoyo.com/event/gacha_info/api/getConfigList"
const gachaListAPI = "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"

var logFileList = map[string]string{
	"常驻祈愿":   "./GachaLog常驻祈愿.json",
	"新手祈愿":   "./GachaLog新手祈愿.json",
	"武器活动祈愿": "./GachaLog武器活动祈愿.json",
	"角色活动祈愿": "./GachaLog角色活动祈愿.json",
}

func main() {
	path, _ := os.UserHomeDir()
	path = path + "/AppData/LocalLow/miHoYo/原神/"
	logUrl := readLogUrl(path)
	saveLogAsJSON(logUrl)
	data := loadLocalLog()
	makeExcelFile(data)
}

func loadLocalLog() map[string][]map[string]string {
	res := make(map[string][]map[string]string)
	for name, fp := range logFileList {
		jsonbody, _ := os.ReadFile(fp)
		v := new([]interface{})
		err := json.Unmarshal(jsonbody, v)
		if err != nil {
			res[name] = nil
			continue
		}
		temp := make([]map[string]string, len(*v))
		for i := range *v {
			vv := (*v)[i].(map[string]interface{})
			ttemp := make(map[string]string)
			for k, v := range vv {
				ttemp[k] = v.(string)
			}
			temp[i] = ttemp
		}
		res[name] = temp
	}
	return res
}

func saveLogAsJSON(logUrl string) {
	types := getGachaTypes(logUrl)
	for i := range types {
		entry := types[i]
		fmt.Println("Start Fetching Pool ", entry["name"])
		z := getFullGachaLog(logUrl, entry["key"])
		saveGachaLog(z, "./GachaLog"+entry["name"]+".json")
	}
	fmt.Println("End Fetching")
}

func getQuery(inurl string) url.Values {
	urlObj, err := url.Parse(inurl)
	if err != nil {
		log.Fatal(err)
	}
	return urlObj.Query()
}

func getGachaTypes(logUrl string) []map[string]string {
	typeAddr := gachaTypeAPI + "?" + getQuery(logUrl).Encode()
	resp, err := http.Get(typeAddr)
	if err != nil {
		log.Fatal(err)
	}
	data := new(map[string](map[string]([]map[string]string)))
	json.NewDecoder(resp.Body).Decode(data)
	return (*data)["data"]["gacha_type_list"]
}

func getGachaLog(logUrl string, key string, page string, end_id string) []interface{} {
	qry := getQuery(logUrl)
	qry.Add("gacha_type", key)
	qry.Add("page", page)
	qry.Add("size", "20")
	qry.Add("end_id", end_id)

	gachaLogUrl := gachaListAPI + "?" + qry.Encode()
	resp, err := http.Get(gachaLogUrl)
	if err != nil {
		log.Fatal(err)
	}

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	d := data.(map[string]interface{})["data"]
	if d == nil {
		log.Fatal("Fetching failed, check the time duration")
		return nil
	}
	return d.(map[string]interface{})["list"].([]interface{})
}

func getFullGachaLog(logUrl, key string) []map[string]string {
	page := 1
	var data []map[string]string
	end_id := "0"
	lenBlock := 1
	for lenBlock > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("\tFetching Page ", page)
		blk := getGachaLog(logUrl, key, fmt.Sprint(page), end_id)
		lenBlock = len(blk)
		if lenBlock > 0 {
			if blk[len(blk)-1] == nil {
				fmt.Println("!")
			}
			end_id = blk[len(blk)-1].(map[string]interface{})["id"].(string)

			page += 1
		} else {
			end_id = "0"
			break
		}
		for x := range blk {

			r := blk[x].(map[string]interface{})
			temp := make(map[string]string)
			for k, v := range r {
				temp[k] = v.(string)
			}
			data = append(data, temp)
		}
	}
	return data
}

func saveGachaLog(data []map[string]string, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if len(data) == 0 {
		return
	}

	jsonData, _ := json.MarshalIndent(data, "", "\t")
	f.Write(jsonData)

	fmt.Println("Data stored in ", filename)
}

func makeExcelFile(dataStruct map[string][]map[string]string) {
	f := excelize.NewFile()
	colKeyMap := map[string]string{"A": "time", "B": "name", "C": "item_type", "D": "rank_type"}
	for name, data := range dataStruct {
		sh := f.NewSheet(name)
		idx := 0
		pdx := 0
		lin := len(data) + 1
		// Set column title
		f.SetCellValue(name, "A1", "时间")
		f.SetCellValue(name, "B1", "名称")
		f.SetCellValue(name, "C1", "类别")
		f.SetCellValue(name, "D1", "星级")
		f.SetCellValue(name, "E1", "总次数")
		f.SetCellValue(name, "F1", "保底内")

		for i := len(data) - 1; i >= 0; i-- {
			for col, key := range colKeyMap {
				f.SetCellValue(name, col+fmt.Sprint(lin), data[i][key])
			}
			f.SetCellValue(name, "E"+fmt.Sprint(lin), idx)
			f.SetCellValue(name, "F"+fmt.Sprint(lin), pdx)
			idx++
			pdx++
			if data[i]["rank_type"] == "3" {
				pdx = 0
			}
			lin--
		}
		f.SetActiveSheet(sh)
	}
	if err := f.SaveAs("data.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func readLogUrl(dir string) string {
	genshinLog, _ := os.ReadFile(dir + "output_log.txt")
	z := string(genshinLog)
	//fmt.Println(z)
	regpatt, _ := regexp.Compile(".*(https://webstatic.mihoyo.com/.*)\n")
	res := regpatt.FindAllStringSubmatch(z, -1)
	if res == nil {
		return ""
	}
	return res[len(res)-1][1]
	//fmt.Println(res[len(res)-1][1])
	//os.Stdout.Write(genshinLog)
}