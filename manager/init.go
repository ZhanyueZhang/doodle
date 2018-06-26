package manager

import (
	"github.com/dearcode/crab/http/client"
	"github.com/dearcode/crab/http/server"
	"github.com/dearcode/crab/log"
	"github.com/dearcode/crab/orm"
	"github.com/juju/errors"

	"github.com/dearcode/doodle/manager/config"
	"github.com/dearcode/doodle/util/rbac"
)

var (
	rbacClient *rbac.Client
	mdb        *orm.DB
	httpClient *client.HTTPClient
)

// ServerInit 初始化HTTP接口.
func ServerInit() error {
	if err := config.Load(); err != nil {
		return err
	}
	mdb = orm.NewDB(config.Manager.DB.IP, config.Manager.DB.Port, config.Manager.DB.Name, config.Manager.DB.User, config.Manager.DB.Passwd, config.Manager.DB.Charset, 10)

	rbacClient = rbac.New(config.Manager.RBAC.Host, config.Manager.RBAC.Token)

	httpClient = client.New().SetLogger(log.GetLogger())

	if err := server.RegisterPath(&domain{}, "/domain"); err != nil {
		return errors.Trace(err)
	}

	if err := server.RegisterPath(&account{}, "/account"); err != nil {
		return errors.Trace(err)
	}

	if err := server.RegisterPrefix(&debug{}, "/debug/pprof/"); err != nil {
		return errors.Trace(err)
	}

	server.RegisterPrefix(&static{}, "/static/")
	server.RegisterPrefix(&static{}, "/")

	server.RegisterPath(&resource{}, "/resource/")
	server.RegisterPath(&resourceInfo{}, "/resource/info")
	server.RegisterPath(&resourceRole{}, "/resource/role/")

	server.RegisterPath(&cluster{}, "/cluster/")
	server.RegisterPath(&clusterInfo{}, "/cluster/info/")
	server.RegisterPath(&clusterNode{}, "/cluster/node/")

	server.RegisterPath(&role{}, "/role/")
	server.RegisterPath(&roleUser{}, "/role/user/")
	server.RegisterPath(&roleInfo{}, "/role/info/")
	server.RegisterPath(&userRole{}, "/user/role/")

	server.RegisterPath(&serviceInfo{}, "/service/info/")
	server.RegisterPath(&service{}, "/service/")

	server.RegisterPath(&nodes{}, "/nodes/")

	server.RegisterPath(&interfaceAction{}, "/interface/")
	server.RegisterPath(&interfaceRegister{}, "/interface/register/")
	server.RegisterPath(&interfaceRun{}, "/interface/run")
	server.RegisterPath(&interfaceInfo{}, "/interface/info")
	server.RegisterPath(&interfaceDeploy{}, "/interface/deploy")

	server.RegisterPath(&variableInfo{}, "/variable/infos")
	server.RegisterPath(&variable{}, "/variable/")

	server.RegisterPath(&appInfo{}, "/application/info")
	server.RegisterPath(&appInfos{}, "/application/infos")
	server.RegisterPath(&app{}, "/application/")

	server.RegisterPath(&relation{}, "/relation/")

	server.RegisterPath(&docs{}, "/docs/")

	server.RegisterPath(&statsSumAction{}, "/stats/sum/")
	server.RegisterPath(&statsTopApplication{}, "/stats/top/app/")
	server.RegisterPath(&statsTopInterface{}, "/stats/top/iface/")
	server.RegisterPath(&statsErrors{}, "/stats/error/")

	return nil
}
