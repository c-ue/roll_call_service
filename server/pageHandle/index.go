package pageHandle

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"html/template"
	"roll_call_service/server/logger"
	"runtime"
	"strconv"
	"time"
)

type THREE_DATYS_DATA struct {
	today              DAY_DATA
	tomorrow           DAY_DATA
	day_after_tomorrow DAY_DATA
}

type DAY_DATA struct {
	all_people    int
	events        []EVENT_DATA
	remain_people int
}

type EVENT_DATA struct {
	event_name  string
	people_num  int
	people_list []string
}

type TABLE_EVENT struct {
	event_name      string
	people_td_num   int
	people_td_list  []string
	people_tm_num   int
	people_tm_list  []string
	people_atm_num  int
	people_atm_list []string
}

type TABLE_DATA struct {
	row_events        []TABLE_EVENT
	people_td_total   int
	people_tm_total   int
	people_atm_total  int
	people_td_remain  int
	people_tm_remain  int
	people_atm_remain int
}

func Index(ctx *fasthttp.RequestCtx) {
	var ConnID = strconv.FormatUint(ctx.ConnID(), 10)
	var log *zap.Logger = logger.Console()

	// -------------------------------------------------------
	// 处理 HTTP 响应数据
	// -------------------------------------------------------
	// HTTP header 构造
	ctx.Response.Header.SetStatusCode(200)
	ctx.Response.Header.SetConnectionClose() // 关闭本次连接, 这就是短连接 HTTP
	ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/html; charset=utf8"))
	ctx.Response.Header.SetBytesKV([]byte("TransactionID"), []byte(ConnID))

	// -------------------------------------------------------
	// 处理逻辑开始
	// -------------------------------------------------------
	templateFileName := "template/index.html"
	t, err := template.New("index.html").ParseFiles(templateFileName)
	if err != nil {
		log.Debug("---------------- Template File Not Found [" + templateFileName + "]-------------")
		return
	}

	//debug data
	//類別 	總計 今天 總計 明天 總計 後天
	//任管 1員 董悅言 1員 陳居億 1員 陳居億
	//人數
	//應到 1員 1員 1員

	//test day
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	Taipei := time.FixedZone("Taipei", secondsEastOfUTC)
	td := time.Date(2020, 3, 6, 10, 0, 0, 0, Taipei)
	tm := time.Date(2020, 3, 7, 10, 0, 0, 0, Taipei)
	atm := time.Date(2020, 3, 8, 10, 0, 0, 0, Taipei)
	var today DAY_DATA = renew(td)
	var tomorrow DAY_DATA = renew(tm)
	var day_after_tomorrow DAY_DATA = renew(atm)
	var rawData THREE_DATYS_DATA
	rawData.today = today
	rawData.tomorrow = tomorrow
	rawData.day_after_tomorrow = day_after_tomorrow
	var table_data TABLE_DATA = ThreeDaysDateToTableData(rawData)
	if err := t.Execute(ctx, table_data); err != nil {
		_, file, _, _ := runtime.Caller(1)
		log.Debug("---------------- Template Produce Error [" + file + ";" + err.Error() + "]-------------")
		return
	}
}

func ThreeDaysDateToTableData(s THREE_DATYS_DATA) TABLE_DATA {
	var table_data TABLE_DATA
	table_data.people_td_total = s.today.all_people
	table_data.people_td_remain = s.today.remain_people
	table_data.people_atm_total = s.day_after_tomorrow.all_people
	table_data.people_atm_remain = s.day_after_tomorrow.remain_people
	table_data.people_tm_total = s.tomorrow.all_people
	table_data.people_tm_remain = s.tomorrow.remain_people
	for i := 0; i < len(s.today.events); i++ {
		var row_event TABLE_EVENT
		row_event.event_name = s.today.events[i].event_name
		row_event.people_td_num = s.today.events[i].people_num
		row_event.people_td_list = s.today.events[i].people_list
		for j := 0; j < len(s.tomorrow.events); j++ {
			if s.tomorrow.events[j].event_name == row_event.event_name {
				row_event.people_tm_list = s.tomorrow.events[j].people_list
				row_event.people_tm_num = s.tomorrow.events[j].people_num
				s.tomorrow.events = append(s.tomorrow.events[:j], s.tomorrow.events[j+1:]...)
			}
		}
		for k := 0; k < len(s.day_after_tomorrow.events); k++ {
			if s.day_after_tomorrow.events[k].event_name == row_event.event_name {
				row_event.people_atm_list = s.day_after_tomorrow.events[k].people_list
				row_event.people_atm_num = s.day_after_tomorrow.events[k].people_num
				s.day_after_tomorrow.events = append(s.day_after_tomorrow.events[:k], s.day_after_tomorrow.events[k+1:]...)
			}
		}
		table_data.row_events = append(table_data.row_events, row_event)
	}
	return table_data
}

//unit test
func renew(date time.Time) DAY_DATA {
	var event1 EVENT_DATA
	event1.event_name = "任管"
	event1.people_num = 1
	event1.people_list = append(event1.people_list, "董悅言")
	var event2 EVENT_DATA
	event2.event_name = "休假"
	event2.people_num = 2
	event2.people_list = append(event2.people_list, "董悅言", "陳居億")
	var event3 EVENT_DATA
	event3.event_name = "補休"
	event3.people_num = 3
	event3.people_list = append(event3.people_list, "董悅言", "陳居億", "陳齊修")

	var day DAY_DATA
	if date.Day() == 6 {
		day.all_people = 8
		day.events = append(day.events, event1)
		day.remain_people = day.all_people - day.events[0].people_num
	} else if date.Day() == 7 {
		day.all_people = 8
		day.events = append(day.events, event1)
		day.events = append(day.events, event2)
		day.remain_people = day.all_people - day.events[0].people_num - day.events[1].people_num
	} else if date.Day() == 8 {
		day.all_people = 8
		day.events = append(day.events, event1)
		day.events = append(day.events, event2)
		day.events = append(day.events, event3)
		day.remain_people = day.all_people - day.events[0].people_num - day.events[1].people_num - day.events[2].people_num
	}
	return day
}
