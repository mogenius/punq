package dtos

type PunqUser struct {
	Id          string `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
}
