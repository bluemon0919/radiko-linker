package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bluemon0919/go-timeext"

	"github.com/gin-gonic/gin"
)

// WebAPI is webAPI information
type WebAPI struct {
}

// NewWebAPI create new WebAPI
func NewWebAPI() *WebAPI {
	return &WebAPI{}
}

// JumpAPI はリンクを受け取ってradikoにジャンプするAPI
// station : "TBS"などのradikoが受け取る放送局の略称を文字列で指定する
// dayOfWeek : time.Weekday型の数値を文字列で指定する
// startTime : 放送開始時間を"27:04"形式(30時間制)の文字列で指定する
func (w *WebAPI) JumpAPI(c *gin.Context) {
	station := c.Param("station")
	dayOfWeek := c.Param("dayOfWeek")
	startTime := c.Param("startTime")
	fmt.Println(station, dayOfWeek, startTime)

	weekday, err := strconv.Atoi(dayOfWeek)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	url := PreviousTimeURL(station, time.Weekday(weekday), startTime)
	if url != "" {
		fmt.Println("url=", url)
		// url = "http://radiko.jp/share/?sid=TBS&t=20210125150000"
		http.Redirect(c.Writer, c.Request, url, 302)
		// リダイレクトのStatusCodeを301にするとブラウザがページをキャッシュしてしまう
	} else {
		c.JSON(http.StatusInternalServerError, url)
	}
}

// Layout30 is timeext.Timeext format
const Layout30 = "2006.01.02-27:04"

// Layout24 is time.Time format
const Layout24 = "2006.01.02-15:04"

// PreviousTimeURL は前回放送のURLを返します
// station : "TBS"などのradikoが受け取る放送局の略称を文字列で指定する
// dayOfWeek : time.Weekday型の数値を指定する
// startTime : 放送開始時間を"27:04"形式(30時間制)の文字列で指定する
func PreviousTimeURL(station string, dayOfWeek time.Weekday, startTime string) string {
	tExt, err := PreviousTime(dayOfWeek, startTime)
	if err != nil {
		return ""
	}
	return GetRadikoURL(station, time.Time(tExt))
}

// GetRadikoURL ラジオのタイムシフトの番組URLを取得する
// station : "TBS"などのradikoが受け取る放送局の略称を文字列で指定する
// startTime : 放送開始時間をtime.Time形式で指定する
func GetRadikoURL(stationID string, start time.Time) string {
	const defaultEndpoint = "http://radiko.jp"
	location, _ := time.LoadLocation("Asia/Tokyo")
	localTime := start.In(location)
	endpoint := "share/?sid=" + stationID + "&t=" + localTime.Format("20060102150405")
	return defaultEndpoint + "/" + endpoint
}

const RadioLayout = "27:04"

// PreviousTime 前回の放送時間を取得する
// startTime,dayOfWeekは30時間表記(6:00-29:59)で指定する
// Sunday 25:00 = Monday 1:00
// returnは30時間表現のtimeext.Timeextで返す
func PreviousTime(dayOfWeek time.Weekday, startTime string) (timeext.TimeExt, error) {
	// ロケーション情報を取得
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return timeext.TimeExt{}, err
	}

	// 指定ロケーション(Asia/Tokyo)の現在時刻を取得
	now := time.Now().In(loc)

	// 指定ロケーション(Asia/Tokyo)の該当時刻を作成
	s, err := timeext.Parse(RadioLayout, startTime)
	if err != nil {
		return timeext.TimeExt{}, err
	}
	start := time.Time(s)
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), start.Hour(), start.Minute(), 0, 0, loc)

	// weekdayの分だけ日付を進める
	wd := dayOfWeek - nextStart.Weekday()
	if wd < 0 {
		wd += 7
	}

	// 24:00以降の時刻の場合は調整する
	// weekdayを1日進める
	isext := nextStart.Hour() >= 0 && nextStart.Hour() < 6
	if isext {
		wd++
		if wd >= 7 {
			wd -= 7
		}
	}
	nextStart = nextStart.AddDate(0, 0, int(wd))

	// 現在の時間を基準にして、前回の指定時刻を取得する
	if now.Before(nextStart) { // now >= old --- ! now < old
		nextStart = nextStart.AddDate(0, 0, -7)
	}
	next30 := timeext.TimeExt(nextStart)
	return next30, nil
}
