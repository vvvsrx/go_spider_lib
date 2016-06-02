package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	//"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	// net包
	//"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	//"regexp"
	"bytes"
	//"strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	BookongProduct.Register()
}

var BookongProduct = &Spider{
	Name:        "booking中文",
	Description: "booking数据抓取 [http://www.booking.com/destination.html]",
	// Pausetime: 300,
	//Keyin:        KEYIN,
	//Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "http://www.booking.com/destination.zh-cn.html", Rule: "目的地首页"})
			//ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"目的地首页": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".flatList a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {

							ctx.AddQueue(&request.Request{
								Url: "http://www.booking.com" + url, 
								Rule: "国家",
								Temp: map[string]interface{}{
									"title": s.Text(),
								},
							})
						}
					})
				},
			},

			"国家": {
				ItemFields: []string{
					"CountryName",
					"CountryShortName",
				},
				ParseFunc: func(ctx *Context) {
					var title string
					ctx.GetTemp("title", &title)

					currentUrl := ctx.GetUrl()
					lastUrl := strings.FieldsFunc(currentUrl, isSlash)
					countryShortName := strings.FieldsFunc(lastUrl[4], isSlash2)

					ctx.Output(map[int]interface{}{
						0: title,
						1: countryShortName[0],
					})

					query := ctx.GetDom()
					query.Find(".general a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							if isCityLink := strings.Contains(url,"/city/")  ; isCityLink {
								ctx.AddQueue(&request.Request{
									Url: "http://www.booking.com" + url, 
									Rule: "城市",
									Temp: map[string]interface{}{
										"title": s.Text(),
									},
								})
							}
						}
					})
				},
			},

			"城市": {
				ItemFields: []string{
					"CityName",
					"CityShortName",
				},
				ParseFunc: func(ctx *Context) {
					var title string
					ctx.GetTemp("title", &title)

					currentUrl := ctx.GetUrl()
					lastUrl := strings.FieldsFunc(currentUrl, isSlash)
					cityShortName := strings.FieldsFunc(lastUrl[5], isSlash2)

					ctx.Output(map[int]interface{}{
						0: title,
						1: cityShortName[0],
					})



					query := ctx.GetDom()

					query.Find(".general a").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							if isHotelLink := strings.Contains(url,"/hotel/")  ; isHotelLink {
								ctx.AddQueue(&request.Request{Url: "http://www.booking.com" + url, Rule: "酒店"})
							}
						}
					})
				},
			},

			"酒店": {
				ItemFields: []string{
					"Name",
					"Summary",
					"ShortName",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					title := query.Find("#hp_hotel_name").Text()
					var summaryBuffer bytes.Buffer
					query.Find("#summary p").Each(func(i int, s *goquery.Selection) {
						summaryBuffer.WriteString(s.Text()+"|||")
					})



					currentUrl := ctx.GetUrl()
					lastUrl := strings.FieldsFunc(currentUrl, isSlash)
					hotelShortName := strings.FieldsFunc(lastUrl[4], isSlash2)

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: title,
						1: summaryBuffer.String(),
						2: hotelShortName[0],
					})
				},
			},
		},
	},
}


func isSlash(r rune) bool {
	return r == '\\' || r == '/'
}
func isSlash2(r rune) bool {
	return r == '.'
}