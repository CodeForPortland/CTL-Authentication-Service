package reporter

type Status struct {
	ReportType ReportType
	Message string
}
type Reporter struct {
	Updater chan Report
	Status Status
	Errors int
	Checks int
}

type ReportType int
type Payload map[string]interface{}

type Report interface {
	ReportType() ReportType
}

type Dict map[string]interface{}

const (
	INITIALIZING        ReportType = 1
	LISTENING           ReportType = 2
	EVENT               ReportType = 3
	CLOSING             ReportType = 4
)

var rtStrings = map[ReportType]string {
	INITIALIZING:           "INITIALIZING",
	LISTENING:              "LISTENING",
	EVENT:                  "EVENT",
	CLOSING:                "CLOSING",
}

func (rt ReportType) String() string {return rtStrings[rt]}

func NewReport(t ReportType) Report {
	switch t {
	case INITIALIZING:
		return &Initializing{}
	}
	return nil
}

func (r Reporter) MakeReport(rt ReportType, message string, payload interface{}) {
	rpt := NewReport(rt)
	r.Status = Status{rpt.ReportType(), message}
	r.Updater <- rpt
	if rt < 0 {
		r.Errors++
	}
	r.Checks++
}

type Initializing struct {
	ID      string
	Details Dict
}

func (rpt *Initializing) ReportType() ReportType {return INITIALIZING}

type Listening struct {
	ID      string
	Details Dict
}

func (rpt *Listening) ReportType() ReportType {return LISTENING}

type Closing struct {
	ID      string
	Details Dict
}

func (rpt *Closing) ReportType() ReportType {return CLOSING}

type Event struct {
	ID      string
	Details Dict
}

func (rpt *Event) ReportType() ReportType {return EVENT}

