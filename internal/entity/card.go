package entity

type Card struct {
	UUID       string `json:"uuid"`
	Term       string `json:"term"`
	Meaning    string `json:"meaning"`
	ModuleUUID string `json:"module_uuid"`
}
