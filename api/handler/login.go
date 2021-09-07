package handler

import (
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/auth"
	"k8s.io/klog/v2"
	mathrand "math/rand"
	"os"
	"sync"
	"time"
)

var (
	letters = []rune("abcdefgjijklmnopqrstuvwxyz0123456789")
	lettersLen = len(letters)
	rand = mathrand.New(mathrand.NewSource(time.Now().UTC().UnixNano()))
	mutex sync.Mutex

	provider auth.Provider
)

type AccessToken struct {
	Token string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	TokenType string `json:"token_type"`
}


// init auth provider
func init() {
	authSystem := os.Getenv("CONTROL_PLANE_USERMNG_API")
	if authSystem != "" {
		provider = auth.NewTecProvider(authSystem)
	}
}

func AuthAccess(user, pwd string) bool {
	// use admin user login
	if user == "admin" && pwd == os.Getenv("CONTROL_ADMIN_PWD") {
		klog.Infof("login user %s", user)
		return true
	}

	// not admin user, to auth user by auth provider
	if provider == nil {
		klog.Error("no auth provider")
		return false
	}
	klog.Infof("use auth provider %s", provider.Name())
	valid, _, err := provider.Auth(user, pwd)
	if err != nil {
		klog.Error(err)
	}

	return valid

}

func GetAccessToken(tokenLen int, tokenTTL int) AccessToken {
	b := make([]rune, tokenLen)
	for i := range b {
		b[i] = letters[intn(lettersLen)]
	}

	return AccessToken{
		Token: string(b),
		ExpiresIn: tokenTTL,
		TokenType: "bearer",
	}

}


func intn(max int) int{
	mutex.Lock()
	defer mutex.Unlock()
	return rand.Intn(max)
}