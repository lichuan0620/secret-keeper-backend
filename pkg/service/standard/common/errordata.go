package common

import (
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type paramElement struct {
	// TestMessageFmtArg is test message format arg. The example value is "test_action"
	TestMessageFmtArg template.HTML
	// TestMessageFmtArg is test message format arg. The example value is "action"
	MessageParams string
}

// errorItem defines error template data struct
type errorItem struct {
	Code     string
	HTTPCode int32
	Message  string
	Comment  string

	// LCCode lowercase the first words of Code
	LCCode string
	// FmtMessage replace {{.*}} with %s of Message
	FmtMessage string
	// TestMessageArgsJoin is full args string used as value of test.go template. The example value is `"test_action", "test_version"`
	TestMessageArgsJoin template.HTML
	// TestMessage is used as value of test.go template
	TestMessage string
	// MessageFmtJoin is full params string for fmt used as value of test.go template. The example value is `test_action, test_version, test_action`
	MessageFmtJoin string
	// MessageParamsJoin is full params string used as value of test.go template. The example value is `test_action, test_version`
	MessageParamsJoin string
	// ParamElements are param element which is uniq for element
	ParamElements []paramElement
}

func removeDuplicateElement(src []string) []string {
	result := make([]string, 0, len(src))
	temp := map[string]struct{}{}
	for _, item := range src {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (e *errorItem) genExtra() {
	if e.Code != "" {
		e.LCCode = strings.ToLower(e.Code[:1]) + e.Code[1:]
	}
	if e.Message != "" {
		exp := regexp.MustCompile("{{([_a-zA-Z][_a-zA-Z0-9]*)}}")
		e.FmtMessage = exp.ReplaceAllString(e.Message, "%s")
		submatch := exp.FindAllStringSubmatch(e.Message, -1)
		testArgsString := make([]string, 0)
		messageFmtParams := make([]string, 0)
		testArgsInterface := make([]interface{}, 0)
		argsNum := strings.Count(e.FmtMessage, "%s")
		if argsNum != len(submatch) {
			panic("args num should be equal to sub match len")
		}
		for i := 1; i <= argsNum; i++ {
			testArgsString = append(testArgsString, fmt.Sprintf(`"test_%s"`, submatch[i-1][1]))
			messageFmtParams = append(messageFmtParams, submatch[i-1][1])
			testArgsInterface = append(testArgsInterface, fmt.Sprintf("test_%s", submatch[i-1][1]))
		}
		e.TestMessage = fmt.Sprintf(e.FmtMessage, testArgsInterface...)
		e.TestMessageArgsJoin = template.HTML(strings.Join(removeDuplicateElement(testArgsString), ", "))
		e.MessageFmtJoin = strings.Join(messageFmtParams, ", ")
		e.MessageParamsJoin = strings.Join(removeDuplicateElement(messageFmtParams), ", ")
		e.ParamElements = func() []paramElement {
			rdTestMessageFmtArg := removeDuplicateElement(testArgsString)
			rdMessageParams := removeDuplicateElement(messageFmtParams)

			ret := make([]paramElement, 0)
			for i := range rdTestMessageFmtArg {
				ret = append(ret, paramElement{
					TestMessageFmtArg: template.HTML(rdTestMessageFmtArg[i]),
					MessageParams:     rdMessageParams[i],
				})
			}

			return ret
		}()
	}
}

// ErrorData defile error items type which implements DataUnmarshaler
type ErrorData []errorItem

// UnmarshalData unmarshal data to error items
func (d *ErrorData) UnmarshalData(_ string, data []byte) error {
	*d = make(ErrorData, 0)
	if err := json.Unmarshal(data, d); err != nil {
		return errors.WithMessage(err, "json unmarshal")
	}
	for i := range *d {
		(*d)[i].genExtra()
	}
	return nil
}
