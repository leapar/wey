package main

import (
	"github.com/kataras/iris"
	"net/http"
	//"io"
	//"os"
	//"io/ioutil"
	"net/http/cookiejar"
	"github.com/facebookgo/httpcontrol"
	"time"
	"github.com/jimuyida/glog"

	//"os"
	//"io"
	//"strings"
	"fmt"
	//"github.com/aws/aws-sdk-go/service/ses"
	//"bytes"
	"strings"
	"path/filepath"
	"os"
	"io"
	//"net/url"
	"net/url"
	"io/ioutil"
	"strconv"
)

const (
	//WEB_HOST = "https://cdn.weshape3d.com/gjbz001/4129"
/*
	WEB_HOST = "https://www.artvrpro.com"
	WEB_HOST2 = "https://images.artvrpro.com"
*/
/*
	WEB_HOST = "https://cdn.weshape3d.com"
	WEB_HOST2 = "https://cdn.weshape3d.com"
*/


	WEB_HOST = "https://ssl-player.720static.com/"
	WEB_HOST2 = "https://ssl-panoimg%s.720static.com"
)

var COOKIES = ""

type tsdbrelayHTTPTransport struct {
	Transport http.RoundTripper
}

func (this *tsdbrelayHTTPTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36")
	//r.Header.Add("Cookie",COOKIES)
	return this.Transport.RoundTrip(r)
}


func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		//beego.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		glog.Error(err)
	}
	//proxyUrl, err := url.Parse("http://192.168.20.162:1080")
	client := &http.Client{
		Jar: jar,
		Transport: &tsdbrelayHTTPTransport{
			&httpcontrol.Transport{
				RequestTimeout:      time.Minute,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: 500,
			//	Proxy:http.ProxyURL(proxyUrl),
			},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	http.DefaultClient = client
}

func DoRequest(url string,request* http.Request) (*http.Response, error){

	//byts,err := ioutil.ReadAll(request.Body)

	//body := bytes.NewBuffer(byts)
	if request.Method == "OPTIONS" {
		request.Method = "POST"
	}
	req,err := http.NewRequest(request.Method,url,request.Body)
	if err != nil {
		return nil,err
	}
	req.Header.Set("Content-Type",request.Header.Get("Content-Type"))

//	req.Header.Set("x-tool-name","diy")


	req.Header.Set("Cookie",COOKIES)
	//req.Header.Set("Referer","https://720yun.com/t/79vku97y5r7?scene_id=38662838")
	//req.Header.Set("Cookie",/*request.Header.Get("Cookie")*/"ktrackerid=8eda961c-34f7-479b-977c-261a75a10ddc; h5diy_ht3=true; h5diy_inform_stacking=true; h5diy_ht4=true; qhdi=717c670552b011e8953b03285bd77c06; gr_user_id=a403c712-cb2d-4bfd-b805-f43d1ed11763; _ga=GA1.2.702865735.1525777841; lastEnteredTool=h5diy; ktrackerid=b1e4686a-70ff-4a40-b024-be597249e568; miniViewWidth=230; miniViewHeight=190; modifyDoorOrWindowTip=Thu%20May%2010%202018; wall.areaCopy.tip=true; isFirstTimeVisit=false; usersource=www.baidu.com; KSESSIONID=846afcf95cda11e8a3a165f1336ec207; _gid=GA1.2.1687682202.1526895422; kjl_usercityid=258; landingpageurl=http%3A%2F%2Fwww.kujiale.com%2F; Hm_lvt_bd8fd4c378d7721976f466053bd4a855=1525932936,1526363860,1526895422,1526954247; DIYSERVERS=1; gr_cs1_20d6d3a3-5784-4939-b125-f03f61047184=userId%3A3FO4JN92W6LA; designer_called_level_api_3FO4JN92W6LA=true; gr_cs1_c3cacbe9-d49e-4323-aa3a-e490b36a040a=userId%3A%E8%AE%BF%E5%AE%A2; gr_cs1_34481ccc-fbbf-41e6-a0f0-c58f61ac7526=userId%3A3FO4JN92W6LA; gr_session_id_a4a13a22eb51522b=34481ccc-fbbf-41e6-a0f0-c58f61ac7526_true; Hm_lpvt_bd8fd4c378d7721976f466053bd4a855=1526971510; _gat=1; JSESSIONID=jkya61p5z9s018mzj74tspic2")
	//"machine_cookie=d83edf9f-83ab-4370-ace8-95039847cc26; 1013-84eb-06w2yb4k46s2=3394a2d2bc4fb1393652e445b7eece1b1614687238171eb783d62b03b8a7393aaa803353039c93b1; user_id=452; 1013-b51f-w82ykf1b22s8=c4537db8-e5ed-4d98-a4b5-e44d4e46595f; _pk_ses..d225=*"
	resp, err := http.DefaultClient.Do(req)
	return resp,err
}

func main() {
	CACHE_DIR   := GetCurrentDirectory() +"/cache/"
	os.MkdirAll(CACHE_DIR,777)



	f,err := os.Open(GetCurrentDirectory() +"/cookie.txt")
	if err != nil {
		fmt.Println("no cookie")
		return;
	}
	defer f.Close()
	bytes,err := ioutil.ReadAll(f)
	COOKIES = string(bytes)


	app := iris.New()
	app.UseGlobal(func(ctx iris.Context) {
		//	ctx.Header("Access-Control-Allow-Headers", "X-Requested-With")
		//	ctx.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		//	ctx.Header("Access-Control-Allow-Origin", "http://127.0.0.1:8887")

		ctx.Next()
	})

	app.OnErrorCode(404, func(ctx iris.Context) {


		fileDir := CACHE_DIR
		orginUrl := ctx.Request().RequestURI
		orginHost := WEB_HOST



		info,err := url.Parse(orginUrl)
		if err != nil {
			return
		}


		orginUrl = info.Path +"?"+ info.RawQuery
		dataFile := fileDir + info.Path
		_,err = os.Stat(dataFile)


		if ctx.Request().Method == "OPTIONS" {
			fmt.Println(ctx.Request().Method)
			ctx.StatusCode(http.StatusOK)
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Headers", "*");
			ctx.Header("Access-Control-Allow-Methods", "*")
			return
		}

		if err == nil {
			fmt.Println("serve:"+dataFile)
			ctx.StatusCode(http.StatusOK)
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Headers", "*");
			ctx.Header("Access-Control-Allow-Methods", "GET")
			ctx.ServeFile(dataFile, false)
			return
		} else {

			resp, err := DoRequest(orginHost + orginUrl, ctx.Request())
			//resp,err := http.DefaultClient.Get(orginHost + orginUrl)
			if err != nil {
				fmt.Println("err1:",err)
				ctx.StatusCode(http.StatusNotFound)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {

				tmps := strings.Split(orginUrl,"/")
				fmt.Println(tmps)

				resp, err = DoRequest(fmt.Sprintf(WEB_HOST2,tmps[1]) + orginUrl[len(tmps[1]) + 1:], ctx.Request())
				//resp,err := http.DefaultClient.Get(orginHost + orginUrl)
				if err != nil {
					fmt.Println("err1:",err)
					ctx.StatusCode(http.StatusNotFound)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Println("StatusCode:", resp.StatusCode, WEB_HOST2+orginUrl)
					fmt.Println("err2:", err)
					ctx.StatusCode(http.StatusNotFound)
					return
				}
			}

			os.MkdirAll(filepath.Dir(dataFile),os.ModePerm)
			f,err := os.Create(dataFile)
			if err != nil {
				fmt.Println("err3:",err)
				ctx.StatusCode(http.StatusNotFound)
				return
			}


			tmp := resp.Header.Get("Content-Length")
			tsize := -1
			if len(tmp) > 0 {
				sz,err := strconv.Atoi(tmp)
				if err == nil {
					tsize = sz
				}
			}
			sz,err := io.Copy(f,resp.Body)
			if err != nil || (tsize > 0 && sz != int64(tsize)) {
				fmt.Println("err4:",err)
				f.Close()
				os.RemoveAll(dataFile)
				ctx.StatusCode(http.StatusNotFound)
				return
			}
			f.Close()
			fmt.Println("cache:"+dataFile)
			ctx.StatusCode(http.StatusOK)

			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Headers", "*");
			ctx.Header("Access-Control-Allow-Methods", "GET")
			ctx.ServeFile(dataFile,false)
		}
		return

	})


	app.Run(iris.Addr(":9991"))
}
