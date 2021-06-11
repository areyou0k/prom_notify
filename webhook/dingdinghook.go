package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"promnotify/prometheusalert"

	"github.com/buger/jsonparser"

	"github.com/blinkbean/dingtalk"
)

type DingDing struct {
	*NotifyBase
}

func NewDingDing(key string, token, member []string) *DingDing {
	return &DingDing{
		&NotifyBase{
			key:    key,
			token:  token,
			member: member,
		},
	}
}

func (d *DingDing) SendMsg() {
	msg := d.msgFormat()
	cli := dingtalk.InitDingTalk(d.token, d.key)
	title := fmt.Sprintf("%s [%s] [%s]",
		d.Alter.Labels.Alertname, d.Alter.Labels.Severity, d.Alter.Status)
	err := cli.SendMarkDownMessageBySlice(title, msg, dingtalk.WithAtMobiles(d.member))
	if err != nil {
		log.Println("dingtalk sendMsg err: ", err)
	}
	// cli.SendMarkDownMessage(title, msg, cli.WithAtMobiles([]string{testPhone}))
}

func (d *DingDing) msgFormat() []string {
	var notifyTitle string
	if d.Alter.Status == "firing" {
		notifyTitle = fmt.Sprintf("### <font color=#FF6100 size=3>%s [%s] [%s]</font>",
			d.Alter.Labels.Alertname, d.Alter.Labels.Severity, d.Alter.Status)
	} else {
		notifyTitle = fmt.Sprintf("### <font color=#00C957 size=3>%s [%s] [%s]</font>",
			d.Alter.Labels.Alertname, d.Alter.Labels.Severity, d.Alter.Status)
	}

	labels, _ := json.Marshal(d.Alter.Labels)
	var notifyLabels string
	jsonparser.ObjectEach(labels, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if string(key) != "alertname" && string(key) != "severity" {
			notifyLabels += fmt.Sprintf("> **%s**: %s \n\n", string(key), string(value))
		}

		return nil
	})
	timeStr := ""
	if d.Alter.Status == "firing" {
		timeAlarmStart := prometheusalert.HandlingTime(d.Alter.StartsAt)
		timeStr = "> **时间**: " + timeAlarmStart + "\n"
	} else {

		timeAlarmStart := prometheusalert.HandlingTime(d.Alter.StartsAt)
		timeAlarmEnd := prometheusalert.HandlingTime(d.Alter.EndsAt)
		timeStr = "> **开始**: " + timeAlarmStart + "\n\n" + "> **结束**: " + timeAlarmEnd + "\n"
	}
	notifyDescription := fmt.Sprintf("> **Description** : |\n\n %s \n", d.Alter.Annotations.Description)
	notifySummary := fmt.Sprintf("> **Summary** : |\n\n %s \n", d.Alter.Annotations.Summary)
	notifyLink := fmt.Sprintf("- [Prometheus](%s)", d.Alter.GeneratorUrl)
	msg := []string{
		notifyTitle,
		notifyLabels,
		timeStr,
		notifySummary,
		notifyDescription,
		notifyLink,
	}
	return msg
}
