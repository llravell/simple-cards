package dto

type AuthRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
