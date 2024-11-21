package kindconfig

import (
	"github.com/NorskHelsenett/ror/pkg/config/configconsts"
	"github.com/NorskHelsenett/ror/pkg/config/rorversion"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var (
	VaultSecret *vault.Secret
	Role        string
	Version     string = "1.1.0"
	Commit      string = "FFFFF"
)

func init() {

	viper.SetDefault(configconsts.VAULT_URL, "http://localhost:8200")
	viper.SetDefault("KUBECTL_BASE_URL", "https://127.0.0.1")
	viper.GetString(configconsts.ROLE)
	viper.SetDefault(configconsts.VERSION, Version)
	viper.SetDefault(configconsts.COMMIT, Commit)
	viper.AutomaticEnv()
}

func Load() {
	environment := viper.GetString(configconsts.ENVIRONMENT)
	rlog.Info("loaded environment", rlog.String("Environment", environment))

	_ = viper.WriteConfig()
}

func GetRorVersion() rorversion.RorVersion {
	return rorversion.NewRorVersion(viper.GetString(configconsts.VERSION), viper.GetString(configconsts.COMMIT))
}
