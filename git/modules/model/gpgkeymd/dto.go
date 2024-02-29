package gpgkeymd

type InsertGpgKeyReqDTO struct {
	Name      string
	Account   string
	PubKeyId  string
	Content   string
	EmailList []string
	Expiry    int64
}

type SimpleGpgKeyDTO struct {
	Id        int64    `json:"id"`
	Name      string   `json:"name"`
	PubKeyId  string   `json:"pubKeyId"`
	Expiry    int64    `json:"expiry"`
	EmailList []string `json:"emailList"`
}
