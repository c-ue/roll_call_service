package pageHandle

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"html/template"
	"roll_call_service/server/config"
	"roll_call_service/server/member/DataStruct"
	CALAPI "roll_call_service/server/member/calendar_api"
	"runtime"
	"strconv"
	"time"
)

type (
	TABLE_DATA struct {
		RowEvents       []TABLE_EVENT
		PeopleTdTotal   int
		PeopleTmTotal   int
		PeopleAtmTotal  int
		PeopleTdRemain  int
		PeopleTmRemain  int
		PeopleAtmRemain int
	}
	TABLE_EVENT struct {
		EventName     string
		PeopleTdNum   int
		PeopleTdList  []string
		PeopleTmNum   int
		PeopleTmList  []string
		PeopleAtmNum  int
		PeopleAtmList []string
		times         []int
	}
)

type MEMBEREVENT struct {
	Name  string
	Event [][]DataStruct.EVENT
}

func Index(ctx *fasthttp.RequestCtx, serverConf config.Config, log *zap.Logger) {
	var ConnID = strconv.FormatUint(ctx.ConnID(), 10)

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
	templateFileName := "template/index.tmpl"
	t := template.Must(template.ParseFiles(templateFileName))

	//debug data
	//類別 	總計 今天 總計 明天 總計 後天
	//任管 1員 董悅言 1員 陳居億 1員 陳居億
	//人數
	//應到 1員 1員 1員

	//test day
	//secondsEastOfUTC := int((8 * time.Hour).Seconds())
	//Taipei := time.FixedZone("Taipei", secondsEastOfUTC)
	//td := time.Date(2020, 3, 6, 10, 0, 0, 0, Taipei)
	//tm := time.Date(2020, 3, 7, 10, 0, 0, 0, Taipei)
	//atm := time.Date(2020, 3, 8, 10, 0, 0, 0, Taipei)
	//var today DAY_DATA = renew(td)
	//var tomorrow DAY_DATA = renew(tm)
	//var day_after_tomorrow DAY_DATA = renew(atm)
	//var rawData THREE_DATYS_DATA
	//rawData.today = today
	//rawData.tomorrow = tomorrow
	//rawData.day_after_tomorrow = day_after_tomorrow
	//var table_data TABLE_DATA = ThreeDaysDateToTableData(rawData)
	members := serverConf.MEMBERS
	var MemberEvent []MEMBEREVENT
	for i := 0; i < len(members); i++ {
		MemberEvent = append(MemberEvent, MEMBEREVENT{
			Name:  members[i].NAME,
			Event: GetMemberEventList(members[i].CalendarId, log),
		})
	}
	var table TABLE_DATA
	table = MemberEvent2TableData(MemberEvent, serverConf, log)

	if err := t.Execute(ctx, table); err != nil {
		_, file, _, _ := runtime.Caller(1)
		log.Debug("---------------- Template Produce Error [" + file + ";" + err.Error() + "]-------------")
		return
	}
}

//type (
//	TABLE_DATA struct {
//		RowEvents       []TABLE_EVENT
//		PeopleTdTotal   int
//		PeopleTmTotal   int
//		PeopleAtmTotal  int
//		PeopleTdRemain  int
//		PeopleTmRemain  int
//		PeopleAtmRemain int
//	}
//	TABLE_EVENT struct {
//		EventName     string
//		PeopleTdNum   int
//		PeopleTdList  []string
//		PeopleTmNum   int
//		PeopleTmList  []string
//		PeopleAtmNum  int
//		PeopleAtmList []string
//	}
//)

func MemberEvent2TableData(memberevent []MEMBEREVENT, serverConf config.Config, log *zap.Logger) TABLE_DATA {
	type eventWithTime struct {
		event []TABLE_EVENT
		times []int
	}
	var Dashboard TABLE_DATA
	var Event []TABLE_EVENT
	now := time.Now()
	eventConf := serverConf.EVENTS
	log.Debug("----------------MemberEvent2TableData-------------")
	log.Debug("Config Event Name List{ ")
	for i := 0; i < len(eventConf); i++ {
		Event = append(Event, TABLE_EVENT{
			EventName:     eventConf[i].NAME,
			PeopleTdNum:   0,
			PeopleTdList:  []string{},
			PeopleTmNum:   0,
			PeopleTmList:  []string{},
			PeopleAtmNum:  0,
			PeopleAtmList: []string{},
			times:         eventConf[i].TIME,
		})
		log.Debug("\t" + eventConf[i].NAME + ", ")
	}
	log.Debug("}\n")
	for k := 0; k < len(Event); k++ {
		for i := 0; i < len(memberevent); i++ {
			name := memberevent[i].Name
			eventList := memberevent[i].Event[0]
			brk := false
			for j := 0; j < len(eventList) && (!brk); j++ {
				log.Debug("Compare " + name + " today event name { " + eventList[j].Name + ", " + Event[k].EventName + "}")
				if eventList[j].Name == Event[k].EventName {
					if len(Event[k].times) == 0 {
						log.Debug("Today " + eventList[j].Name + " Add People " + name)
						Event[k].PeopleTdNum++
						Event[k].PeopleTdList = append(Event[k].PeopleTdList, name)
						brk = true
					}
					for l := 0; l < len(Event[k].times); l++ {
						log.Debug("Compare " + name + " today event time { " + strconv.Itoa(Event[k].times[l]) + ", " + strconv.Itoa(now.Hour()) + "}")
						if now.Hour() == Event[k].times[l] {
							log.Debug("Today " + eventList[j].Name + " Add People " + name)
							Event[k].PeopleTdNum++
							Event[k].PeopleTdList = append(Event[k].PeopleTdList, name)
							brk = true
						}
					}
				}
			}
		}
		for i := 0; i < len(memberevent); i++ {
			name := memberevent[i].Name
			eventList := memberevent[i].Event[1]
			for j := 0; j < len(eventList); j++ {
				log.Debug("Compare " + name + " tomorrow event name { " + eventList[j].Name + ", " + Event[k].EventName + "}")
				if eventList[j].Name == Event[k].EventName {
					log.Debug("Tomorrow " + eventList[j].Name + " Add People " + name)
					Event[k].PeopleTmNum++
					Event[k].PeopleTmList = append(Event[k].PeopleTmList, name)
					break
				}
			}
		}
		for i := 0; i < len(memberevent); i++ {
			name := memberevent[i].Name
			eventList := memberevent[i].Event[2]
			for j := 0; j < len(eventList); j++ {
				log.Debug("Compare " + name + " after tomorrow event name { " + eventList[j].Name + ", " + Event[k].EventName + "}")
				if eventList[j].Name == Event[k].EventName {
					log.Debug("After Tomorrow " + eventList[j].Name + " Add People " + name)
					Event[k].PeopleAtmNum++
					Event[k].PeopleAtmList = append(Event[k].PeopleAtmList, name)
					break
				}
			}
		}
	}
	Dashboard.RowEvents = Event
	Dashboard.PeopleTdTotal = len(serverConf.MEMBERS)
	Dashboard.PeopleTmTotal = Dashboard.PeopleTdTotal
	Dashboard.PeopleAtmTotal = Dashboard.PeopleTdTotal
	for i := 0; i < len(Event); i++ {
		Dashboard.PeopleTdRemain = Event[i].PeopleTdNum + Dashboard.PeopleTdRemain
		Dashboard.PeopleTmRemain = Event[i].PeopleTmNum + Dashboard.PeopleTmRemain
		Dashboard.PeopleAtmRemain = Event[i].PeopleAtmNum + Dashboard.PeopleAtmRemain
	}
	Dashboard.PeopleTdRemain = Dashboard.PeopleTdTotal - Dashboard.PeopleTdRemain
	Dashboard.PeopleTmRemain = Dashboard.PeopleTmTotal - Dashboard.PeopleTmRemain
	Dashboard.PeopleAtmRemain = Dashboard.PeopleAtmTotal - Dashboard.PeopleAtmRemain

	log.Debug("{People Today Remain: " + strconv.Itoa(Dashboard.PeopleTdRemain) + ", People Today Total  " + strconv.Itoa(Dashboard.PeopleTdTotal) + "}")
	log.Debug("{People Tomorrow Remain: " + strconv.Itoa(Dashboard.PeopleTmRemain) + ", People Tomorrow Total  " + strconv.Itoa(Dashboard.PeopleTmTotal) + "}")
	log.Debug("{People After Tomorrow Remain: " + strconv.Itoa(Dashboard.PeopleAtmRemain) + ", People After Tomorrow Total  " + strconv.Itoa(Dashboard.PeopleAtmTotal) + "}")
	log.Debug("Remove Empty Event {")
	for i := 0; i < len(Dashboard.RowEvents); i++ {
		if Dashboard.RowEvents[i].PeopleTdNum == 0 && Dashboard.RowEvents[i].PeopleTmNum == 0 && Dashboard.RowEvents[i].PeopleAtmNum == 0 {
			log.Debug("EventName: " + Dashboard.RowEvents[i].EventName + ",")
			copy(Dashboard.RowEvents[i:], Dashboard.RowEvents[i+1:])
			Dashboard.RowEvents[len(Dashboard.RowEvents)-1] = TABLE_EVENT{}
			Dashboard.RowEvents = Dashboard.RowEvents[:len(Dashboard.RowEvents)-1]
			i--
		}
	}
	log.Debug("}")
	return Dashboard

}

func GetMemberEventList(calId string, log *zap.Logger) [][]DataStruct.EVENT {
	now := time.Now()
	var eventList [][]DataStruct.EVENT
	eventList = append(eventList, CALAPI.GetSpecifyDay(calId, now, log))
	eventList = append(eventList, CALAPI.GetSpecifyDay(calId, now.AddDate(0, 0, 1), log))
	eventList = append(eventList, CALAPI.GetSpecifyDay(calId, now.AddDate(0, 0, 2), log))
	return eventList
}

func PrintableEvent(eventList [][]DataStruct.EVENT, serverConf config.Config) [][]DataStruct.EVENT {
	now := time.Now()
	filterEvent := serverConf.EVENTS
	var retEventList [][]DataStruct.EVENT
	for i := 0; i < len(eventList); i++ {
		var retEventArray []DataStruct.EVENT = nil
		for j := 0; j < len(eventList[i]); j++ {
			for k := 0; k < len(filterEvent); k++ {
				if eventList[i][j].Name == filterEvent[k].NAME {
					if now.Sub(eventList[i][j].StartTime) >= 0 && now.Sub(eventList[i][j].StartTime) <= 0 {
						if len(filterEvent[k].TIME) == 0 {
							retEventArray = append(retEventArray, eventList[i][j])
						}
						for l := 0; l < len(filterEvent[k].TIME); l++ {
							if filterEvent[k].TIME[l] == now.Hour() {
								retEventArray = append(retEventArray, eventList[i][j])
							}
						}
					} else if now.Sub(eventList[i][j].StartTime) < 0 {
						retEventArray = append(retEventArray, eventList[i][j])
					}
				}
			}
		}
		retEventList = append(retEventList, retEventArray)
	}
	return retEventList
}

func listElementRemove(list []DataStruct.EVENT, index int) []DataStruct.EVENT {
	list[index] = list[len(list)-1]
	list[len(list)-1] = DataStruct.EVENT{}
	list = list[:len(list)-1]
	return list
}

////TODO:rewrite for collcate all member obj
//func ThreeDaysDateToTableData(s THREE_DATYS_DATA) TABLE_DATA {
//	var table_data TABLE_DATA
//	table_data.PeopleTdTotal = s.today.all_people
//	table_data.PeopleTdRemain = s.today.remain_people
//	table_data.PeopleAtmTotal = s.day_after_tomorrow.all_people
//	table_data.PeopleAtmRemain = s.day_after_tomorrow.remain_people
//	table_data.PeopleTmTotal = s.tomorrow.all_people
//	table_data.PeopleTmRemain = s.tomorrow.remain_people
//	for i := 0; i < len(s.today.events); i++ {
//		var row_event TABLE_EVENT
//		row_event.EventName = s.today.events[i].event_name
//		row_event.PeopleTdNum = s.today.events[i].people_num
//		row_event.PeopleTdList = s.today.events[i].people_list
//		for j := 0; j < len(s.tomorrow.events); j++ {
//			if s.tomorrow.events[j].event_name == row_event.EventName {
//				row_event.PeopleTmList = s.tomorrow.events[j].people_list
//				row_event.PeopleTmNum = s.tomorrow.events[j].people_num
//				s.tomorrow.events = append(s.tomorrow.events[:j], s.tomorrow.events[j+1:]...)
//			}
//		}
//		for k := 0; k < len(s.day_after_tomorrow.events); k++ {
//			if s.day_after_tomorrow.events[k].event_name == row_event.EventName {
//				row_event.PeopleAtmList = s.day_after_tomorrow.events[k].people_list
//				row_event.PeopleAtmNum = s.day_after_tomorrow.events[k].people_num
//				s.day_after_tomorrow.events = append(s.day_after_tomorrow.events[:k], s.day_after_tomorrow.events[k+1:]...)
//			}
//		}
//		table_data.RowEvents = append(table_data.RowEvents, row_event)
//	}
//	for i := 0; i < len(s.tomorrow.events); i++ {
//		var row_event TABLE_EVENT
//		row_event.EventName = s.tomorrow.events[i].event_name
//		row_event.PeopleTmNum = s.tomorrow.events[i].people_num
//		row_event.PeopleTmList = s.tomorrow.events[i].people_list
//		for k := 0; k < len(s.day_after_tomorrow.events); k++ {
//			if s.day_after_tomorrow.events[k].event_name == row_event.EventName {
//				row_event.PeopleAtmList = s.day_after_tomorrow.events[k].people_list
//				row_event.PeopleAtmNum = s.day_after_tomorrow.events[k].people_num
//				s.day_after_tomorrow.events = append(s.day_after_tomorrow.events[:k], s.day_after_tomorrow.events[k+1:]...)
//			}
//		}
//		table_data.RowEvents = append(table_data.RowEvents, row_event)
//	}
//	for i := 0; i < len(s.day_after_tomorrow.events); i++ {
//		var row_event TABLE_EVENT
//		row_event.EventName = s.day_after_tomorrow.events[i].event_name
//		row_event.PeopleAtmNum = s.day_after_tomorrow.events[i].people_num
//		row_event.PeopleAtmList = s.day_after_tomorrow.events[i].people_list
//		table_data.RowEvents = append(table_data.RowEvents, row_event)
//	}
//	return table_data
//}
//
////TODO:
//// rewrite in member obj
//// unit test
//func renew(date time.Time) DAY_DATA {
//	var event1 EVENT_DATA
//	event1.event_name = "任管"
//	event1.people_num = 1
//	event1.people_list = append(event1.people_list, "董悅言")
//	var event2 EVENT_DATA
//	event2.event_name = "休假"
//	event2.people_num = 2
//	event2.people_list = append(event2.people_list, "董悅言", "陳居億")
//	var event3 EVENT_DATA
//	event3.event_name = "補休"
//	event3.people_num = 3
//	event3.people_list = append(event3.people_list, "董悅言", "陳居億", "陳齊修")
//
//	var day DAY_DATA
//	if date.Day() == 6 {
//		day.all_people = 8
//		day.events = append(day.events, event1)
//		day.remain_people = day.all_people - day.events[0].people_num
//	} else if date.Day() == 7 {
//		day.all_people = 8
//		day.events = append(day.events, event1)
//		day.events = append(day.events, event2)
//		day.remain_people = day.all_people - day.events[0].people_num - day.events[1].people_num
//	} else if date.Day() == 8 {
//		day.all_people = 8
//		day.events = append(day.events, event1)
//		day.events = append(day.events, event2)
//		day.events = append(day.events, event3)
//		day.remain_people = day.all_people - day.events[0].people_num - day.events[1].people_num - day.events[2].people_num
//	}
//	return day
//}
