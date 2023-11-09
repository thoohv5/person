package template

import _ "embed"

//go:embed request.tpl
var request string

func GetRequest() string {
	return request
}

//go:embed response.tpl
var response string

func GetResponse() string {
	return response
}

//go:embed service.tpl
var service string

func GetService() string {
	return service
}

//go:embed controller.tpl
var controller string

func GetController() string {
	return controller
}

//go:embed condition.tpl
var condition string

func GetCondition() string {
	return condition
}

//go:embed dao.tpl
var dao string

func GetDao() string {
	return dao
}

//go:embed model.tpl
var model string

func GetModel() string {
	return model
}
