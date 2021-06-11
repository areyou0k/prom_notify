package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"promnotify/prometheusalert"
	yamlconfig "promnotify/yamlconfig"
)

type NotifyBase struct {
	Retry  int
	key    string
	token  []string
	member []string
	Alter  *prometheusalert.Alerts
}

var (
	Critical  = make(chan WebhookInterf, 1000)
	OtherLeve = make(chan WebhookInterf, 1000)
)

func Hierarchy(p *prometheusalert.Prometheus, notifyType, notifyAisle, notifyGroup string) {

	labels, _ := json.Marshal(p)
	log.Println(string(labels))
	switch notifyType {
	case "dingtalk":
		for _, alter := range p.Alerts {
			makeDingdingNotice(alter, notifyType, notifyAisle, notifyGroup)
		}
	case "telephone":
		for _, alter := range p.Alerts {
			makeTelephoneNotice(alter, notifyType, notifyAisle, notifyGroup)
		}
	default: //默认dingtalk
		for _, alter := range p.Alerts {
			makeDingdingNotice(alter, "default", "default", "default")
		}
	}
}

func WebhookWorker() {
	for {
		select {
		// case <-stopCh:
		// 	return
		case critical := <-Critical:
			critical.SendMsg()
		case otherLeve := <-OtherLeve:
		priority:
			for {
				select {
				case critical := <-Critical:
					critical.SendMsg()
				default:
					break priority
				}
			}
			otherLeve.SendMsg()
		}
	}
}

func makeDingdingNotice(alter prometheusalert.Alerts, notifyType, notifyAisle, notifyGroup string) {
	// debug
	ding := NewDingDing(
		// "dingtalk.config.argocd.key"
		yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.key", notifyType, notifyAisle)),
		//"dingtalk.config.argocd.token"
		yamlconfig.ConfigYml.GetStringSlice(fmt.Sprintf("%s.config.%s.token", notifyType, notifyAisle)),
		// "dingtalk.group.test"
		yamlconfig.ConfigYml.GetStringSlice(fmt.Sprintf("%s.group.%s", notifyType, notifyGroup)),
	)
	currentAlter := alter
	ding.Alter = &currentAlter
	if alter.Labels.Severity == "critical" {
		Critical <- ding
	} else {
		OtherLeve <- ding
	}
}

func makeTelephoneNotice(alter prometheusalert.Alerts, notifyType, notifyAisle, notifyGroup string) {
	aliSms := NewAliSms(
		// "telephone.config.ali_sms.access_key"
		yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.access_key", notifyType, notifyAisle)),
		//"telephone.config.ali_sms.sign_Name"
		yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.sign_Name", notifyType, notifyAisle)),
		yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.template_code", notifyType, notifyAisle)),
		yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.Template_param", notifyType, notifyAisle)),
		[]string{yamlconfig.ConfigYml.GetString(fmt.Sprintf("%s.config.%s.access_key_secret", notifyType, notifyAisle))},
		yamlconfig.ConfigYml.GetStringSlice(fmt.Sprintf("%s.group.%s", notifyType, notifyGroup)),
	)
	currentAlter := alter
	aliSms.n.Alter = &currentAlter
	if alter.Labels.Severity == "critical" {
		Critical <- aliSms
	} else {
		OtherLeve <- aliSms
	}
}
