package common

// ErrorBase defines error base struct
type ErrorBase struct {
	HTTPCode int32
	Code     string
	Message  string
	Data     map[string]string

	DataPreset map[string]string
}

// GetHTTPCode returns http code of the error
func (e *ErrorBase) GetHTTPCode() int32 {
	return e.HTTPCode
}

// GetCode returns code of the error
func (e *ErrorBase) GetCode() string {
	return e.Code
}

// GetMessage returns message of the error
func (e *ErrorBase) GetMessage() string {
	return e.Message
}

// GetData returns data map of the error
func (e *ErrorBase) GetData() map[string]string {
	return func(mObj ...map[string]string) map[string]string {
		var ret map[string]string
		for _, m := range mObj {
			for k, v := range m {
				if ret == nil {
					ret = make(map[string]string)
				}
				ret[k] = v
			}
		}
		return ret
	}(e.DataPreset, e.Data)
}
