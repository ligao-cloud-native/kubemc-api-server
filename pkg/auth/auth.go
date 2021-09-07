package auth

import "fmt"

type Provider interface {
	Name() string
	Auth(user, pwd string) (valid, privileged bool, err error)
}

type BaseProvider struct{}

func (*BaseProvider) Name() string {
	return "template"
}

func (*BaseProvider) Auth(user, pwd string) (valid, privileged bool, err error) {
	return false, false, fmt.Errorf("error")
}
