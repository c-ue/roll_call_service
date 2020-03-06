package member

import (
	"time"
)

type (
	EVENT struct {
		name      string
		startTime time.Time
		endTime   time.Time
	}
	TOKEN struct {
		url string
	}
	MEMBER struct {
		name   string
		token  TOKEN
		events []EVENT
	}
)

func NewMember(name string, token TOKEN) *MEMBER {
	p := new(MEMBER)
	p.name = name
	p.token = token
	return p
}

func (p *MEMBER) Name() string {
	return p.name
}

func (p *MEMBER) Token() TOKEN {
	return p.token
}

//TODO:get spec time event
func (p *MEMBER) GetEvent(t time.Time) string {

}

func (p *MEMBER) ChangeName(name string) {
	p.name = name
}

func (p *MEMBER) ChangeToken(t TOKEN) {
	p.token = t
}

//TODO:update this week event
// "c00ao9qlffo3okr0s33ima5kas@group.calendar.google.com"
func (p *MEMBER) updateNext7DayEvent() {
}
