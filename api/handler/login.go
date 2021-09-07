package handler

import (
	"github.com/ligao-cloud-native/kubemc-api-server/pkg/auth"
	"k8s.io/klog/v2"
	"os"
)

var provider auth.Provider

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
