package entity

type Module struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	UserUUID string `json:"user_uuid"`
}

type ModuleWithCards struct {
	Module
	Cards []*Card `json:"cards"`
}
