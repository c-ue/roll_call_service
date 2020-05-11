package main

import (
	"fmt"
	"go.uber.org/zap"
	"roll_call_service/server/logger"
	"roll_call_service/server/member/DataStruct"
	CalApi "roll_call_service/server/member/calendar_api"
	"time"
)

func main() {
	var log *zap.Logger = logger.Console()
	var t []DataStruct.EVENT
	t = CalApi.GetSpecifyDay("v8tor6mf5g97dil7e57bsjksfg@group.calendar.google.com", time.Now(), log)
	for i := 0; i < len(t); i++ {
		fmt.Print(t[i].Name, "=", t[i].StartTime.Format(time.RFC3339), "=", t[i].EndTime.Format(time.RFC3339), "\n")
	}
	fmt.Printf("%v", t)
}
