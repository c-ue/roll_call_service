package calendar_api

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"roll_call_service/server/member/DataStruct"
	"time"
)

func GetSpecifyDay(calendarId string, date time.Time, log *zap.Logger) []DataStruct.EVENT {
	//TODO: pop it out to member
	log.Debug("---------------- GOOGLE CALENDAR api -------------")
	log.Debug("Get Calendar Id: " + calendarId)
	loc, _ := time.LoadLocation("Asia/Taipei")
	b, err := ioutil.ReadFile("config/calendar_account_credentials.json")
	if err != nil {
		log.Fatal("Unable to read client secret file")

	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatal("Unable to parse client secret file to config")
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatal("Unable to retrieve Calendar client")
	}
	date = date.In(loc)
	specDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	events, err := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(specDate.Format(time.RFC3339)).TimeMax(specDate.AddDate(0, 0, 1).
		Format(time.RFC3339)).OrderBy("startTime").Do()
	if err != nil {
		log.Fatal("Unable to retrieve the calendarId's events: " + err.Error())
		log.Fatal("CalendarId: " + calendarId)
	}
	var event_list []DataStruct.EVENT
	for _, item := range events.Items {
		var startT time.Time
		startT, err := formatEventDateTimeToGolangTimeTime(item.Start.Date, item.Start.DateTime, loc)
		if err != nil {
			log.Fatal("StartDateTime in formatEventDateTimeToGolangTimeTime error.[" + item.Start.Date + "][" + item.Start.DateTime + "]")
		}
		var endT time.Time
		endT, err = formatEventDateTimeToGolangTimeTime(item.End.Date, item.End.DateTime, loc)
		if err != nil {
			log.Fatal("EndDateTime in formatEventDateTimeToGolangTimeTime error.[" + item.Start.Date + "][" + item.Start.DateTime + "]")
		}
		event_list = append(event_list, DataStruct.EVENT{
			Name:      item.Summary,
			StartTime: startT,
			EndTime:   endT,
		})
		log.Debug("{Event Name: \"" + item.Summary + "\", Start Time: " + item.Start.Date + ", End Time:" + item.Start.DateTime + "}")
	}
	return event_list
}

func formatEventDateTimeToGolangTimeTime(date string, datetime string, loc *time.Location) (time.Time, error) {
	var ret time.Time
	var err error
	if date != "" {
		ret, err = time.ParseInLocation("2006-01-02", date, loc)
		if err != nil {
			log.Fatal("Date ParseInLocation Error. [" + date + "]\n")
		}
	} else if datetime != "" {
		ret, err = time.ParseInLocation(time.RFC3339, datetime, loc)
		if err != nil {
			log.Fatal("Date ParseInLocation Error. [" + datetime + "]\n")
		}
	}
	return ret, err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
