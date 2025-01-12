package dto

type CreateCardRequest struct {
	Term    string `json:"term"    validate:"required"`
	Meaning string `json:"meaning" validate:"required"`
}

type UpdateCardRequest struct {
	Term    string `json:"term"    validate:"required_without=Meaning"`
	Meaning string `json:"meaning" validate:"required_without=Term"`
}
