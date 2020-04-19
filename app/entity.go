package app

import (
	"bytes"
	"net/http"
	"toconline-casper-go/model"
)

// GetEntityList retives all entities from current user
func (a *App) GetEntityList() ([]*model.Entity, error) {

	httpHeader := http.Header{}
	httpHeader.Add("Content-Type", "application/json")
	httpHeader.Add("Accept", "application/json")

	entityList := a.Request.SetMessage(
		model.VerbGet,
		model.MessageTarget{
			Target:  model.TargetHTTP,
			URL:     "https://cdb-master.toconline.pt/cdb/entity/list",
			Headers: nil,
		},
		nil,
	)

	message, err := a.CasperRequest(entityList)

	if err != nil {
		return nil, err
	}

	return model.EntityListFromJSON(bytes.NewReader(message.Params)), nil
}

// EntitySwitch ...
func (a *App) EntitySwitch(entityID int) (*model.EntityFiscalYears, error) {
	entitySwitch := a.Request.SetMessage(
		model.VerbPut,
		model.MessageTarget{
			Target:  model.TargetHTTP,
			URL:     "https://cdb-master.toconline.pt/cdb/entity/switch",
			Headers: nil,
		},
		map[string]interface{}{
			"action":       "switch",
			"to_entity_id": entityID,
		},
	)

	message, err := a.CasperRequest(entitySwitch)

	if err != nil {
		return nil, err
	}

	return model.EntityFiscalYearsFromJSON(bytes.NewReader(message.Params)), nil
}

// GetEntityReport ...
func (a *App) GetEntityReport() {
	// entityReport := a.Request.SetMessage(
	// 	model.VerbPut,
	// 	map[string]interface{}{
	// 		"target":   "job-queue",
	// 		"tube":     "job-controller",
	// 		"ttr":      3000,
	// 		"validaty": 259200,
	// 	},
	// 	[]byte{`
	// 		"tube":"job-controller",
	// 		"disable_notification":true,
	// 		"destination_tube":"casper-print-queue-hd",
	// 		"notification_title":"Geração do Relatório - Balancetes",
	// 		"follow_up":interface{
	// 		   "tube":"job-follow-up"
	// 		},
	// 		"timeout":259200,
	// 		"validity":259200,
	// 		"ttr":3000,
	// 		"locale":"pt_PT",
	// 		"name":"balanceteperiodoacumulado_514244127_19_04_2020",
	// 		"public_link":{
	// 		   "path":"download"
	// 		},
	// 		"documents":[
	// 		   {
	// 			  "jrxml":"default/trial_balance_period.jrxml",
	// 			  "name":"balanceteperiodoacumulado_514244127_19_04_2020",
	// 			  "title":"balanceteperiodoacumulado_514244127_19_04_2020",
	// 			  "jsonapi":{
	// 				 "urn":"https://app1.toconline.pt/trial_balance_period/0?filter[tbs_type]=Balancete&filter[tbs_level]=Movimento&filter[period_name]=Janeiro%20(2020)&filter[balance_type]=Saldo%20das%20somas&filter[period_month]=1&filter[level]=9&filter[trial_balance_sh]=Per%C3%ADodo%2C%20Acumulado&filter[balance_type_id]=account_balance&filter[include_zeroes]=false&filter[only_opening_balance]=false&filter[include_opening_balance]=true&filter[include_closing]=false&filter[include_yearly_balance]=false&filter[show_short_descriptions]=false&filter[hide_third_party_accounts]=false"
	// 			  }
	// 		   }
	// 		],
	// 		"overridable_system_variables":{
	// 		   "CERTIFIED_SOFTWARE_NOTICE":"Emitido por TOConline - https://www.toconline.pt"
	// 		},
	// 		"x_brand":"toconline",
	// 		"x_product":"toconline",
	// 		"app_digest":"202004170031"
	// 	`},
	// )

	_, err := a.CasperRequestRaw([]byte(`2:PUT:{"target":"job-queue","tube":"job-controller","ttr":3000,"validity":259200}:{"tube":"job-controller","disable_notification":true,"destination_tube":"casper-print-queue-hd","conflicts_tubes_messages":[{"tube":"casper-print-queue-hd","title":"Geração de relatórios duplicada","sub_title":"Neste momento o software já se encontra a gerar um relatório.","message":"Pode acompanhar o progresso da geração desse mesmo relatório na área de notificações. Assim que este terminar poderá aceder ao mesmo numa aba nova e proceder à criação de um novo."}],"notification_title":"Geração do Relatório - Balancetes","follow_up":{"tube":"job-follow-up"},"timeout":259200,"validity":259200,"ttr":3000,"locale":"pt_PT","name":"balanceteperiodoacumulado_514244127_19_04_2020","public_link":{"path":"download"},"documents":[{"jrxml":"default/trial_balance_period.jrxml","name":"balanceteperiodoacumulado_514244127_19_04_2020","title":"balanceteperiodoacumulado_514244127_19_04_2020","jsonapi":{"urn":"https://app1.toconline.pt/trial_balance_period/0?filter[tbs_type]=Balancete&filter[tbs_level]=Movimento&filter[period_name]=Janeiro%20(2020)&filter[balance_type]=Saldo%20das%20somas&filter[period_month]=1&filter[level]=9&filter[trial_balance_sh]=Per%C3%ADodo%2C%20Acumulado&filter[balance_type_id]=account_balance&filter[include_zeroes]=false&filter[only_opening_balance]=false&filter[include_opening_balance]=true&filter[include_closing]=false&filter[include_yearly_balance]=false&filter[show_short_descriptions]=false&filter[hide_third_party_accounts]=false","prefix":null,"user_id":null,"entity_id":null,"entity_schema":null,"sharded_schema":null,"subentity_prefix":null,"subentity_schema":null}}],"overridable_system_variables":{"CERTIFIED_SOFTWARE_NOTICE":"Emitido por TOConline - https://www.toconline.pt"},"user_id":null,"entity_id":null,"entity_schema":null,"sharded_schema":null,"subentity_id":null,"subentity_schema":null,"subentity_prefix":null,"user_email":null,"role_mask":null,"module_mask":null,"x_brand":"toconline","x_product":"toconline","app_digest":"202004170031"}`))

	if err != nil {
		return
	}

	entityReportGenerate := a.Request.SetMessage(
		model.VerbSubscribe,
		map[string]interface{}{
			"target": "job",
			"tube":   "casper-print-queue-hd",
			"id":     36392470,
		},
		nil,
	)

	_, err = a.CasperRequest(entityReportGenerate)

	if err != nil {
		return
	}
}
