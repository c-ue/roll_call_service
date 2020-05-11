package member

//
//import (
//	"roll_call_service/server/member/DataStruct"
//	"roll_call_service/server/member/calendar_api"
//	"time"
//)
//
////TODO:ReDesign member obj
//type MEMBER struct {
//	calendarId string
//	events []DataStruct.EVENT
//	focusEvents []DataStruct.EVENT
//	printableEvent []string
//}
//
//func (m *MEMBER)SetCalendarId(id string){
//	m.calendarId = id
//}
//
//func (m *MEMBER)GetCalendarId() string{
//	return m.calendarId
//}
//
//func (m *MEMBER)UpdataEvents(dates []time.Time){
//	m.events = nil;
//	for i:=0;i<len(dates);i++ {
//		m.events = append(m.events, calendar_api.GetSpecifyDay(m.calendarId, dates[i])...);
//	}
//}
//
//func (m *MEMBER)GetEvents()[]DataStruct.EVENT{
//	return m.events
//}
//
//func (m *MEMBER)SetFocusEvents(events []DataStruct.EVENT){
//	m.focusEvents = events
//}
//
//func (m *MEMBER)GetFocusEvents()[]DataStruct.EVENT{
//	return m.focusEvents
//}
//
//func (m *MEMBER)GetPrintableEvent()[]string{
//	m.printableEvent = nil;
//	for i:=0;i<len(m.events);i++{
//		for j:=0;j<len(m.focusEvents);j++{
//			if m.events[i].Name == m.focusEvents[j].Name{
//				if
//			}
//		}
//	}
//	return
//}
//
//func (m *MEMBER)NewMember()*MEMBER{
//	m = new(MEMBER)
//	return m
//}
//
//
//
////input: GetSpecDay event
////input: configure event set
////proc: event filter
////proc: day event to time event(with configure setting)
////output: GetSpecTime event
////TODO:Dsign group obj member
////input: multi member
////output: group SpecTime event
//
//////TODO:update this week event
////// "c00ao9qlffo3okr0s33ima5kas@group.calendar.google.com"
