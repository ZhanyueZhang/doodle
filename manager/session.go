package manager

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dearcode/crab/cache"
	"github.com/dearcode/crab/http/client"
	"github.com/dearcode/crab/log"
	"github.com/juju/errors"

	"github.com/dearcode/doodle/manager/config"
)

const (
	//session会话超时30分钟
	sessionTimeout = 1800
)

type sessionCache struct {
	cache *cache.Cache
}

func newSession() *sessionCache {
	return &sessionCache{cache: cache.NewCache(sessionTimeout)}
}

//getToken 读取用户cookie中token.
func (s *sessionCache) getToken(r *http.Request) (string, error) {
	c, err := r.Cookie(config.Manager.SSO.Key)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

//verify 调用sso接口验证token返回用户信息.
func (s *sessionCache) verify(r *http.Request, token string) (*userinfo, error) {
	resp := struct {
		Status  int
		Message string
		Data    userinfo
	}{}
	url := fmt.Sprintf("%s?token=%s", config.Manager.SSO.VerifyURL, token)
	buf, err := client.New().Get(url, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if err = json.Unmarshal(buf, &resp); err != nil {
		return nil, errors.Trace(err)
	}

	if resp.Status != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp.Data, nil
}

func (u userinfo) String() string {
	return u.Email
}

//loadInfo 加载资源与角色信息.
func (u *userinfo) loadInfo() error {
	res, err := userdb.loadResource(u.Email)
	if err != nil {
		return errors.Trace(err)
	}
	u.setResource(res)

	roles, err := userdb.loadRoles(u.Email)
	if err != nil {
		return errors.Trace(err)
	}

	u.setRoles(roles)

	u.IsAdmin = userdb.isAdmin(u.Email)

	return nil
}

func (s *sessionCache) User(r *http.Request) (*userinfo, error) {
	token, err := s.getToken(r)
	if err != nil {
		return nil, errors.Trace(err)
	}

	log.Debugf("token:%v", token)
	if val := s.cache.Get(token); val != nil {
		user := val.(*userinfo)
		log.Debugf("cache userinfo:%v", user)
		return user, user.loadInfo()
	}

	user, err := s.verify(r, token)
	if err != nil {
		return nil, errors.Trace(err)
	}

	log.Debugf("userinfo:%+v", user)

	if err = user.loadInfo(); err != nil {
		return nil, errors.Trace(err)
	}

	id, err := userdb.loadUserID(user.Email)
	if err != nil {
		return nil, errors.Trace(err)
	}
	user.UserID = id

	s.cache.Add(token, user)

	return user, nil
}
