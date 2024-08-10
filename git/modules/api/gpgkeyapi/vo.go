package gpgkeyapi

type CreateGpgKeyReqVO struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type GpgKeyVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	KeyId   string `json:"keyId"`
	Expired string `json:"expired"`
	Created string `json:"created"`
	Email   string `json:"email"`
	SubKeys string `json:"subKeys"`
}
