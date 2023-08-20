package dtos

type PunqContext struct {
	Id            string `json:"id" validate:"required"`
	ContextBase64 string `json:"contextBase64" validate:"required"`
}
