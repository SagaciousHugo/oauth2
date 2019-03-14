package manager

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/sagacioushugo/oauth2"
	"github.com/sagacioushugo/oauth2/generator"
	"github.com/sagacioushugo/oauth2/store"
	"time"
)

func init() {
	Register("mem", &store.MemTokenStore{})
	Register("default", &generator.Default{})
}

var tokenStores = make(map[string]store.TokenStore)

var generators = make(map[string]generator.Generator)

func Register(name string, target interface{}) {
	if target == nil {
		panic("oauth2 manager: Register target is nil")
	}
	switch t := target.(type) {
	case store.TokenStore:
		if _, dup := tokenStores[name]; dup {
			panic("oauth2 manager: Register called twice for tokenStore " + name)
		}
		newTarget, _ := target.(store.TokenStore)
		tokenStores[name] = newTarget
	case generator.Generator:
		if _, dup := generators[name]; dup {
			panic("oauth2 manager: Register called twice for generator " + name)
		}
		newTarget, _ := target.(generator.Generator)
		generators[name] = newTarget
	default:
		panic(fmt.Sprintf("oauth2 manager: Unsupported register type %v", t))
	}
}

type Manager struct {
	config     *oauth2.ManagerConfig
	tokenStore store.TokenStore
	generator  generator.Generator
}

func NewManager(managerConfig *oauth2.ManagerConfig) *Manager {
	var ok bool
	var tokenStore store.TokenStore
	var generator generator.Generator

	tokenStore, ok = tokenStores[managerConfig.TokenStoreName]
	if !ok {
		panic(fmt.Errorf("oauth2 manager: unknow tokenStore %q", managerConfig.TokenStoreName))
	}
	if err := tokenStore.Init(managerConfig.TokenStoreConfig); err != nil {
		panic(err)
	}

	generator, ok = generators[managerConfig.GeneratorName]
	if !ok {
		panic(fmt.Errorf("oauth2 manager: unknow generator %q", managerConfig.GeneratorName))
	}

	return &Manager{
		config:     managerConfig,
		tokenStore: tokenStore,
		generator:  generator,
	}
}

func (manager *Manager) TokenGC() {
	manager.tokenStore.GC(manager.config.TokenGcInterval)
	time.AfterFunc(time.Duration(manager.config.TokenGcInterval)*time.Second, func() { manager.TokenGC() })
}

//tokenStore wrap
func (manager *Manager) TokenNew(ctx *context.Context) store.Token {
	return manager.tokenStore.NewToken(ctx)
}

func (manager *Manager) TokenCreate(token store.Token) error {
	return manager.tokenStore.Create(token)
}

func (manager *Manager) TokenGetByAccess(access string) (store.Token, error) {
	return manager.tokenStore.GetByAccess(access)
}

func (manager *Manager) TokenGetByRefresh(refresh string) (store.Token, error) {
	return manager.tokenStore.GetByRefresh(refresh)
}
func (manager *Manager) TokenGetByCode(code string) (store.Token, error) {
	return manager.tokenStore.GetByCode(code)
}

// generator wrap
func (manager *Manager) GenerateCode(token store.Token, req *oauth2.Request, ctx *context.Context) error {
	if code, err := manager.generator.Code(req, ctx); err != nil {
		return err
	} else {
		token.SetCode(code)
		return manager.tokenStore.Create(token)
	}
}

func (manager *Manager) GenerateToken(token store.Token, req *oauth2.Request, ctx *context.Context, isGenerateRefresh bool) error {
	if access, refresh, err := manager.generator.Token(req, ctx, isGenerateRefresh); err != nil {
		return err
	} else {
		if isGenerateRefresh {
			token.SetRefresh(refresh)
		}
		token.SetAccess(access)
		return manager.tokenStore.Create(token)
	}
}

func (manager *Manager) GenerateTokenAndDelToken(token store.Token, tokenToDel store.Token, req *oauth2.Request, ctx *context.Context, isGenerateRefresh bool) error {
	if access, refresh, err := manager.generator.Token(req, ctx, isGenerateRefresh); err != nil {
		return err
	} else {
		if isGenerateRefresh {
			token.SetRefresh(refresh)
		}
		token.SetAccess(access)
		return manager.tokenStore.CreateAndDel(token, tokenToDel)
	}
}
