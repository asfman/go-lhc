package main

import (
  "sort"
  "regexp"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"github.com/gin-gonic/gin"
)

type Lhc struct {
	Date string `json:"date"`
	Sno  string `json:"sno"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/lhc", func(c *gin.Context) {
		url := "https://bet.hkjc.com/contentserver/jcbw/cmc/last30draw.json"
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		resp, respErr := client.Do(req)
		if respErr != nil {
			c.String(http.StatusInternalServerError, "respErr %s", respErr.Error())
			return
		}
		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			result, _ := ioutil.ReadAll(resp.Body)
			var lhc []Lhc
			error := json.Unmarshal(bytes.TrimPrefix(result, []byte("\xef\xbb\xbf")), &lhc)
			fmt.Println(lhc)
			if error != nil {
				fmt.Println(error)
				c.String(http.StatusOK, error.Error())
				return
			}
			showAll := c.Query("lhc")
			if len(showAll) == 0 {
				fmt.Println("showAll: " + showAll)
				splitDate := strings.Split(lhc[0].Date, ("/"))
				sno := lhc[0].Sno
				c.String(http.StatusOK, "【特码：" + format(sno) + ", 日期: " + splitDate[2] + "-" + splitDate[1] + "-" + splitDate[0] + "】")
			} else {
        var fString = stat(lhc) + "\n\n"
				if showAll == "clean" {
					for index, item := range lhc {
						if index != len(lhc)-1 {
							fString += format(item.Sno) + ", "
						} else {
							fString += format(item.Sno)
						}
					}
				} else {
					for index, item := range lhc {
						splitDate := strings.Split(item.Date, ("/"))
						if index != len(lhc)-1 {
							fString += splitDate[2] + "-" + splitDate[1] + "-" + splitDate[0] + ": " + format(item.Sno) + "\n"
						} else {
							fString += splitDate[2] + "-" + splitDate[1] + "-" + splitDate[0] + ": " + format(item.Sno)
						}
					}
				}
				c.String(http.StatusOK, fString)
			}
			return
		}
		c.String(http.StatusInternalServerError, "status: %s", resp.Status)
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8888")
}

func stat(lhc []Lhc) string {
  var st map[string]int
  st = make(map[string]int)
  var nstat map[string]int
  nstat = make(map[string]int)
  for _, item := range lhc {
    no, _ := strconv.Atoi(item.Sno)
    st[xs[no%12]] += 1
    nstat[strconv.Itoa(no%10)] += 1
  }
  var ret = "最近" + strconv.Itoa(len(lhc)) + "次开奖，"
  ret += sortMap(st)
  ret += "\n\n"
  ret += sortMap(nstat)
  return ret
}

func sortMap(m map[string]int) string {
   var keys []string
   for k := range m {
      keys = append(keys, k)
   }
   sort.Strings(keys)
   var ret string = ""
   for _, k := range keys {
     ret += k +  "出现" + strconv.Itoa(m[k]) + "次, "
   }
   r, _ := regexp.Compile(", $")
   ret = r.ReplaceAllString(ret, "")
   return ret
}

func format(sno string) string {
	no, _ := strconv.Atoi(sno)
	return sno + "(" + xs[no%12] + ")"
}


var xs map[int]string = make(map[int]string, 12)
func init() {
  xs[1] = "猪"
	xs[2] = "狗"
	xs[3] = "鸡"
	xs[4] = "猴"
	xs[5] = "羊"
	xs[6] = "马"
	xs[7] = "蛇"
	xs[8] = "龙"
	xs[9] = "兔"
	xs[10] = "虎"
	xs[11] = "牛"
	xs[0] = "鼠"
}
