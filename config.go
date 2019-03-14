package oauth2

type Config struct {
	TokenType          string
	RedirectUriSep     string
	RedirectAllowEmpty bool
	AllowGrantType     map[GrantType]GrantTypeConfig
	ManagerConfig      ManagerConfig
}

type ManagerConfig struct {
	TokenGcInterval  int64
	TokenStoreName   string
	TokenStoreConfig string
	GeneratorName    string
}

type GrantTypeConfig struct {
	CodeExpire         int64
	AccessTokenExpire  int64
	RefreshTokenExpire int64
	IsGenerateRefresh  bool
	IsResetRefreshTime bool
}

var DefaultOauth2Config = Config{
	TokenType:          "Bearer",
	RedirectUriSep:     "|",
	RedirectAllowEmpty: false,
	AllowGrantType: map[GrantType]GrantTypeConfig{
		AuthorizationCode: {
			AccessTokenExpire:  12 * 3600,
			RefreshTokenExpire: 72 * 3600,
			CodeExpire:         300,
			IsGenerateRefresh:  true,
			IsResetRefreshTime: false,
		},
		Implicit: {
			AccessTokenExpire:  6 * 3600,
			IsGenerateRefresh:  false,
			IsResetRefreshTime: false,
		},
		PasswordCredentials: {
			AccessTokenExpire:  6 * 3600,
			IsGenerateRefresh:  false,
			IsResetRefreshTime: false,
		},
		ClientCredentials: {
			AccessTokenExpire:  6 * 3600,
			IsGenerateRefresh:  false,
			IsResetRefreshTime: false,
		},
		RefreshToken: {
			AccessTokenExpire:  12 * 3600,
			RefreshTokenExpire: 72 * 3600,
			IsGenerateRefresh:  true,
			IsResetRefreshTime: false,
		},
	},
	ManagerConfig: ManagerConfig{
		TokenGcInterval:  7 * 24 * 3600,
		TokenStoreName:   "mem",
		GeneratorName:    "default",
		TokenStoreConfig: "",
	},
}

func NewDefaultConfig() *Config {
	return &DefaultOauth2Config
}
