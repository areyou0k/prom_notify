package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"promnotify/prometheusalert"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/buger/jsonparser"
)

type AliSms struct {
	n             *NotifyBase
	SignName      string
	TemplateCode  string
	TemplateParam string
}

func NewAliSms(key, signName, templateCode, templateParam string,
	token, member []string) *AliSms {
	return &AliSms{
		SignName:      signName,
		TemplateParam: templateParam,
		TemplateCode:  templateCode,
		n: &NotifyBase{
			key:    key,
			token:  token,
			member: member,
		},
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (a *AliSms) CreateClient() (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(a.n.key),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(a.n.token[0]),
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func (a *AliSms) SendMsg() {
	client, _err := a.CreateClient()
	if _err != nil {
		log.Println(_err)
	}
	msg := map[string]string{
		"code": a.msgFormat(),
	}
	mjson, _ := json.Marshal(msg)
	member := strings.Join(a.n.member, ",")
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(member),
		SignName:      tea.String(a.SignName),
		TemplateCode:  tea.String(a.TemplateCode),
		TemplateParam: tea.String(string(mjson)),
	}

	ret, err := client.SendSms(sendSmsRequest)
	if *(ret.Body.Code) != "OK" || err != nil {
		log.Println("发送失败，添加到队列重试,RET: ", ret, "\nERR: ", err)
		if a.n.Alter.Labels.Severity == "critical" && a.n.Retry < 3 {
			// if a.n.Alter.Labels["severity"] == "critical" && a.n.Retry < 3 {
			Critical <- a
			a.n.Retry += 1
		} else if a.n.Retry < 3 {
			OtherLeve <- a
		}
	}
}

func (a *AliSms) msgFormat() string {
	// notifyTitle := fmt.Sprintf("%s [%s] [%s]\n", a.n.Alter.Labels["alertname"], a.n.Alter.Labels["severity"], a.n.Alter.Status)
	notifyTitle := fmt.Sprintf("\n%s [%s] [%s]\n", a.n.Alter.Labels.Alertname, a.n.Alter.Labels.Severity, a.n.Alter.Status)
	labels, _ := json.Marshal(a.n.Alter.Labels)
	var notifyLabels string
	jsonparser.ObjectEach(labels, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		if string(key) != "alertname" && string(key) != "severity" {
			notifyLabels += fmt.Sprintf("%s: %s \n", string(key), string(value))
		}

		return nil
	})

	notifySummary := fmt.Sprintf("Summary:  |\n    %s\n", a.n.Alter.Annotations.Summary)
	notifyDescription := fmt.Sprintf("Description:  |\n    %s", a.n.Alter.Annotations.Description)

	timeStr := ""
	if a.n.Alter.Status == "firing" {
		timeAlarmStart := prometheusalert.HandlingTime(a.n.Alter.StartsAt)
		timeStr = "时间: " + timeAlarmStart + "\n"
	} else {
		timeAlarmStart := prometheusalert.HandlingTime(a.n.Alter.StartsAt)
		timeAlarmEnd := prometheusalert.HandlingTime(a.n.Alter.EndsAt)
		timeStr = "开始: " + timeAlarmStart + "\n" + "结束: " + timeAlarmEnd + "\n"
	}

	msg := []string{
		notifyTitle,
		notifyLabels,
		timeStr,
		notifySummary,
		notifyDescription,
	}

	msgStr := strings.Join(msg, "\n")
	return msgStr
}
