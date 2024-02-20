package ap

type PublicKey struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}
