package network

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/pkg/errors"
)

func TestParseEndpoint(t *testing.T) {
	type TestCase struct {
		Raw           string
		DefaultScheme string
		Parsed        string
	}
	testGroups := []*struct {
		Name  string
		Error bool
		Cases []*TestCase
	}{
		{
			Name: "empty-scheme",
			Cases: []*TestCase{
				{Raw: "localhost", Parsed: "http://localhost"},
				{Raw: "localhost:80", Parsed: "http://localhost:80"},
				{Raw: "localhost:30000", Parsed: "http://localhost:30000"},
				{Raw: "foo.default.svc.cluster.local", Parsed: "http://foo.default.svc.cluster.local"},
				{Raw: "127.0.0.1", Parsed: "http://127.0.0.1"},
				{Raw: "127.0.0.1:80", Parsed: "http://127.0.0.1:80"},
				{Raw: "127.0.0.1:30000", Parsed: "http://127.0.0.1:30000"},
			},
		},
		{
			Name: "default-scheme",
			Cases: []*TestCase{
				{Raw: "localhost", DefaultScheme: "http", Parsed: "http://localhost"},
				{Raw: "https://localhost", DefaultScheme: "http", Parsed: "https://localhost"},
				{Raw: "localhost:80", DefaultScheme: "https", Parsed: "https://localhost:80"},
				{Raw: "http://localhost:80", DefaultScheme: "https", Parsed: "http://localhost:80"},
				{Raw: "foo.default.svc.cluster.local", DefaultScheme: "postgres", Parsed: "postgres://foo.default.svc.cluster.local"},
				{Raw: "http://foo.default.svc.cluster.local", DefaultScheme: "postgres", Parsed: "http://foo.default.svc.cluster.local"},
				{Raw: "127.0.0.1", DefaultScheme: "http", Parsed: "http://127.0.0.1"},
				{Raw: "https://127.0.0.1", DefaultScheme: "http", Parsed: "https://127.0.0.1"},
				{Raw: "127.0.0.1:80", DefaultScheme: "https", Parsed: "https://127.0.0.1:80"},
				{Raw: "http://127.0.0.1:80", DefaultScheme: "https", Parsed: "http://127.0.0.1:80"},
				{Raw: "127.0.0.1:30000", DefaultScheme: "postgres", Parsed: "postgres://127.0.0.1:30000"},
				{Raw: "https://127.0.0.1:30000", DefaultScheme: "postgres", Parsed: "https://127.0.0.1:30000"},
			},
		},
		{
			Name: "non-empty-scheme",
			Cases: []*TestCase{
				{Raw: "hTtP://localhost", Parsed: "http://localhost"},
				{Raw: "hTtP://localhost:80", Parsed: "http://localhost:80"},
				{Raw: "hTtP://localhost:30000", Parsed: "http://localhost:30000"},
				{Raw: "hTtP://foo.default.svc.cluster.local", Parsed: "http://foo.default.svc.cluster.local"},
				{Raw: "hTtP://127.0.0.1", Parsed: "http://127.0.0.1"},
				{Raw: "hTtP://127.0.0.1:80", Parsed: "http://127.0.0.1:80"},
				{Raw: "hTtP://127.0.0.1:30000", Parsed: "http://127.0.0.1:30000"},
				{Raw: "https://localhost", Parsed: "https://localhost"},
				{Raw: "https://localhost:80", Parsed: "https://localhost:80"},
				{Raw: "https://localhost:30000", Parsed: "https://localhost:30000"},
				{Raw: "https://foo.default.svc.cluster.local", Parsed: "https://foo.default.svc.cluster.local"},
				{Raw: "https://127.0.0.1", Parsed: "https://127.0.0.1"},
				{Raw: "https://127.0.0.1:80", Parsed: "https://127.0.0.1:80"},
				{Raw: "https://127.0.0.1:30000", Parsed: "https://127.0.0.1:30000"},
				{Raw: "postgres://postgres:Pwd123456@localhost/raven", Parsed: "postgres://postgres:Pwd123456@localhost/raven"},
				{Raw: "postgres://postgres:Pwd123456@localhost:80/raven", Parsed: "postgres://postgres:Pwd123456@localhost:80/raven"},
				{Raw: "postgres://postgres:Pwd123456@localhost:30000/raven", Parsed: "postgres://postgres:Pwd123456@localhost:30000/raven"},
				{Raw: "postgres://postgres:Pwd123456@foo.default.svc.cluster.local/raven", Parsed: "postgres://postgres:Pwd123456@foo.default.svc.cluster.local/raven"},
				{Raw: "postgres://postgres:Pwd123456@127.0.0.1/raven", Parsed: "postgres://postgres:Pwd123456@127.0.0.1/raven"},
				{Raw: "postgres://postgres:Pwd123456@127.0.0.1:80/raven", Parsed: "postgres://postgres:Pwd123456@127.0.0.1:80/raven"},
				{Raw: "postgres://postgres:Pwd123456@127.0.0.1:30000/raven", Parsed: "postgres://postgres:Pwd123456@127.0.0.1:30000/raven"},
			},
		},
		{
			Name: "escape",
			Cases: []*TestCase{
				{Raw: "测试.bytedance.com", Parsed: "http://%E6%B5%8B%E8%AF%95.bytedance.com"},
				{Raw: "bytedance.com/测试", Parsed: "http://bytedance.com/%E6%B5%8B%E8%AF%95"},
			},
		},
		{
			Name:  "invalid-address",
			Error: true,
			Cases: []*TestCase{
				{Raw: ""},
				{Raw: "http:localhost"},
				{Raw: "http:/localhost"},
				{Raw: "http:///localhost"},
				{Raw: "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"},
				{Raw: "<tag>"},
			},
		},
	}
	for _, tg := range testGroups {
		t.Run(tg.Name, func(tt *testing.T) {
			handleResult := func(parsed *url.URL, expectation string, err error) error {
				if err != nil {
					return errors.Wrap(err, "unexpected parse error")
				}
				if parsedString := parsed.String(); parsedString != expectation {
					return errors.Errorf("got: %s; want: %s", parsedString, expectation)
				}
				return nil
			}
			if tg.Error {
				handleResult = func(parsed *url.URL, _ string, err error) error {
					if err == nil {
						return errors.Errorf("parsing \"%s\" succeeded while expecting error", parsed)
					}
					return nil
				}
			}
			for _, tc := range tg.Cases {
				parsed, err := ParseEndpoint(tc.Raw, tc.DefaultScheme)
				if err = handleResult(parsed, tc.Parsed, err); err != nil {
					tt.Error(err)
				}
			}
		})
	}
}

func TestParseBasicAuth(t *testing.T) {
	const passwordFile = "raven-unit-test-password.tmp"
	const passwordFileBody = "Pwd123456"
	const password = "UnitTestPw123"
	const username = "admin"
	if err := ioutil.WriteFile(passwordFile, []byte(passwordFileBody), 0660); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(passwordFile)
	}()
	tcs := []struct {
		Username       string
		Password       string
		PasswordFile   string
		ExpectPassword string
		ExpectError    bool
	}{
		{},
		{
			Username:       username,
			Password:       password,
			ExpectPassword: password,
		},
		{
			Username:       username,
			PasswordFile:   passwordFile,
			ExpectPassword: passwordFileBody,
		},
		{
			Username:     username,
			Password:     password,
			PasswordFile: passwordFile,
			ExpectError:  true,
		},
		{
			PasswordFile: passwordFile,
			ExpectError:  true,
		},
		{
			Password:    password,
			ExpectError: true,
		},
		{
			Username:     username,
			PasswordFile: "non-exist-file.tmp",
			ExpectError:  true,
		},
	}
	for i := range tcs {
		tc := &tcs[i]
		uname, passwd, err := ParseBasicAuth(tc.Username, tc.Password, tc.PasswordFile)
		if err != nil {
			if !tc.ExpectError {
				t.Errorf("tc %d, unexpected error %v", i, err)
			}
		} else if tc.ExpectError {
			t.Errorf("tc %d, expecting error but function call succeeded", i)
		} else {
			if uname != tc.Username {
				t.Errorf("tc %d, username missmatch, expecting %s, got %s", i, tc.Username, uname)
			}
			if passwd != tc.ExpectPassword {
				t.Errorf("tc %d, password missmatch, expecting %s, got %s", i, tc.Password, passwd)
			}
		}
	}
}
