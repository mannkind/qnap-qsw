package qnap

import "encoding/json"

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
	Mode     POEMode `json:"Mode"`
	Priority string  `json:"Priority"`
}

type POEMode int

var POEModes = struct {
	Unknown     POEMode
	Disabled    POEMode
	POE         POEMode
	POEPlus     POEMode
	POEPlusPlus POEMode
}{
	Unknown:     0,
	Disabled:    1,
	POE:         2,
	POEPlus:     3,
	POEPlusPlus: 4,
}

func (t POEMode) String() string {
	return [...]string{"unknown", "disable", "poeDot3af", "poePlusDot3at", "poePlusDot3bt"}[t]
}

func (t POEMode) FromString(kind string) POEMode {
	return map[string]POEMode{
		"unknown":  POEModes.Unknown,
		"disabled": POEModes.Disabled,
		"poe":      POEModes.POE,
		"poe+":     POEModes.POEPlus,
		"poe++":    POEModes.POEPlusPlus,
	}[kind]
}

func (t POEMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *POEMode) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = t.FromString(s)
	return nil
}
