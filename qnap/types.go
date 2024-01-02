package qnap

type loginResponse struct {
	ErrorCode    int                 `json:"error_code"`
	ErrorMessage string              `json:"error_message"`
	Result       loginResultResponse `json:"result"`
}

type loginResultResponse struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

type interfacesResponse struct {
	ErrorCode    int                        `json:"error_code"`
	ErrorMessage string                     `json:"error_message"`
	Result       []interfacesResultResponse `json:"result"`
}

type interfacesResultResponse struct {
	Key string                  `json:"key"`
	Val interfacesValueResponse `json:"val"`
}

type interfacesValueResponse struct {
	Mode     string `json:"Mode"`
	Priority string `json:"Priority"`
}
