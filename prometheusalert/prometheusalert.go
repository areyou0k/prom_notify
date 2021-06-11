package prometheusalert

import (
	"fmt"
	"strconv"
	"time"
)

type Labels struct {
	Alertname        string `json:"alertname,omitempty"`
	Instance         string `json:"instance,omitempty"`
	Severity         string `json:"severity,omitempty"`
	Env              string `json:"env,omitempty"`
	Job              string `json:"job,omitempty"`
	Region           string `json:"region,omitempty"`
	Service          string `json:"service,omitempty"`
	Application      string `json:"application,omitempty"`
	Status           string `json:"status,omitempty"`
	Method           string `json:"method,omitempty"`
	Uri              string `json:"uri,omitempty"`
	ExportedInstance string `json:"exported_instance,omitempty"`
}

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type Alerts struct {
	Status string
	Labels Labels `json:"labels"`
	// Labels       map[string]string `json:"labels"`
	Annotations  Annotations `json:"annotations"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	GeneratorUrl string      `json:"generatorURL"`
}

type Prometheus struct {
	Status      string
	Alerts      []Alerts
	Externalurl string `json:"externalURL"`
}

func HandlingTime(timeStr string) string {
	cstZone := time.FixedZone("GMT", 8*3600) // 东八
	ts, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	sinceTime := time.Since(ts).Minutes()
	ts = ts.In(cstZone)
	timeInfo := ""
	if sinceTime > float64(10) {
		timeInfo = fmt.Sprintf("%s %smin ago", ts.Format("01-02 15:04"), strconv.Itoa(int(sinceTime)))
	} else {
		timeInfo = ts.Format("01-02 15:04")
	}
	return timeInfo
}
