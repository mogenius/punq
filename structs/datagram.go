package structs

import (
	"fmt"
	"time"

	"github.com/mogenius/punq/utils"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type Datagram struct {
	Id        string      `json:"id" validate:"required"`
	Pattern   string      `json:"pattern" validate:"required"`
	Payload   interface{} `json:"payload,omitempty"`
	Err       string      `json:"err,omitempty"`
	CreatedAt time.Time   `json:"-"`
}

func CreateDatagramRequest(request Datagram, data interface{}) Datagram {
	datagram := Datagram{
		Id:        request.Id,
		Pattern:   request.Pattern,
		Payload:   data,
		CreatedAt: request.CreatedAt,
	}
	return datagram
}

func CreateDatagramFrom(pattern string, data interface{}) Datagram {
	datagram := Datagram{
		Id:        uuid.New().String(),
		Pattern:   pattern,
		Payload:   data,
		CreatedAt: time.Now(),
	}
	return datagram
}

func CreateDatagram(pattern string) Datagram {
	datagram := Datagram{
		Id:        uuid.New().String(),
		Pattern:   pattern,
		CreatedAt: time.Now(),
	}
	return datagram
}

func CreateDatagramAck(pattern string, id string) Datagram {
	datagram := Datagram{
		Id:        id,
		Pattern:   pattern,
		CreatedAt: time.Now(),
	}
	return datagram
}

func CreateEmptyDatagram() Datagram {
	datagram := Datagram{
		Id:        uuid.New().String(),
		Pattern:   "",
		CreatedAt: time.Now(),
	}
	return datagram
}

func (d *Datagram) DisplayBeautiful() {
	IDCOLOR := color.New(color.FgWhite, color.BgBlue).SprintFunc()
	PATTERNCOLOR := color.New(color.FgBlack, color.BgYellow).SprintFunc()
	TIMECOLOR := color.New(color.FgWhite, color.BgRed).SprintFunc()
	PAYLOADCOLOR := color.New(color.FgBlack, color.BgHiGreen).SprintFunc()

	fmt.Printf("%s %s\n", IDCOLOR("ID:      "), d.Id)
	fmt.Printf("%s %s\n", PATTERNCOLOR("PATTERN: "), color.BlueString(d.Pattern))
	fmt.Printf("%s %s\n", TIMECOLOR("TIME:    "), time.Now().Format(time.RFC3339))
	fmt.Printf("%s %s\n", TIMECOLOR("Duration:"), DurationStrSince(d.CreatedAt))

	// f := colorjson.NewFormatter()
	// f.Indent = 2
	// s, _ := f.Marshal(d.Payload)
	// PrettyPrintString(d.Payload)

	fmt.Printf("%s %s\n\n", PAYLOADCOLOR("PAYLOAD: "), PrettyPrintString(d.Payload))
}

func (d *Datagram) DisplayReceiveSummary() {
	fmt.Println()
	fmt.Printf("%s%s%s (%s)\n", utils.FillWith("RECEIVED", 23, " "), utils.FillWith(d.Pattern, 60, " "), color.BlueString(d.Id), DurationStrSince(d.CreatedAt))
}

func (d *Datagram) DisplaySentSummary() {
	fmt.Printf("%s%s%s (%s)\n", utils.FillWith("SENT", 23, " "), utils.FillWith(d.Pattern, 60, " "), color.BlueString(d.Id), DurationStrSince(d.CreatedAt))
}

func (d *Datagram) DisplaySentSummaryEvent(kind string, reason string, msg string, count int32) {
	fmt.Printf("%s%s: %s/%s -> %s (Count: %d)\n", utils.FillWith("SENT", 23, " "), d.Pattern, kind, reason, msg, count)
}

func (d *Datagram) DisplayStreamSummary() {
	fmt.Printf("%s%s%s\n", utils.FillWith("STREAMING", 23, " "), utils.FillWith(d.Pattern, 60, " "), color.BlueString(d.Id))
}
