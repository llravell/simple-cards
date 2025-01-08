package dto

type CreateOrUpdateModuleRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type QuizletImportRequest struct {
	ModuleName      string `json:"module_name"       validate:"required,max=100"`
	QuizletModuleID string `json:"quizlet_module_id" validate:"required"`
}
