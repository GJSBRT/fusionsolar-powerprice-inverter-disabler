package fusionsolar

type PublicKeyResponse struct {
	Version       string `json:"version"`
	PubKey        string `json:"pubKey"`
	Timestamp     int64  `json:"timeStamp"`
	EnableEncrypt bool   `json:"enableEncrypt"`
}

type ValidateUserRequest struct {
	OrganizationName string `json:"organizationName"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

type ValidateUserResponse struct {
	ErrorCode           string      `json:"errorCode"`
	ErrorMsg            interface{} `json:"errorMsg"`
	RedirectURL         interface{} `json:"redirectURL"`
	RespMultiRegionName []string    `json:"respMultiRegionName"`
	ServiceURL          interface{} `json:"serviceUrl"`
	VerifyCodeCreate    bool        `json:"verifyCodeCreate"`
	TwoFactorStatus     interface{} `json:"twoFactorStatus"`
}

type KeepAliveResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Payload string      `json:"payload"`
}
