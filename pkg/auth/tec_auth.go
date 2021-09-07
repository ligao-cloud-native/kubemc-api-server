package auth

import "fmt"

const TecProviderName = "tec"

type TecProvider struct {
	BaseProvider
	authAPIUrl string
}

func NewTecProvider(apiAddr string) *TecProvider {
	if apiAddr == "" {
		return nil
	}

	return &TecProvider{authAPIUrl: apiAddr}
}

func (p *TecProvider) Name() string {
	return TecProviderName
}

func (*TecProvider) Auth(user, pwd string) (valid, privileged bool, err error) {
	//TODO: call auth system
	return false, false, fmt.Errorf("error")
}
