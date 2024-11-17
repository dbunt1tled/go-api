package token

type Tokens struct {
	ID           int64  `json:"id,omitempty" jsonapi:"primary,credential"`
	AccessToken  string `json:"accessToken" jsonapi:"attr,accessToken"`
	RefreshToken string `json:"refreshToken" jsonapi:"attr,refreshToken"`
}
