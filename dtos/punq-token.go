package dtos

type PunqToken struct {
	Token string `json:"token" validate:"required"`
}

func CreateToken(token string) *PunqToken {
	return &PunqToken{
		Token: token,
	}
}
