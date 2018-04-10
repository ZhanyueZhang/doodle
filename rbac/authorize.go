package rbac

import (
	"net/http"

	"github.com/dearcode/crab/http/server"
	"github.com/dearcode/crab/log"
	"github.com/dearcode/crab/orm"
)

type authorize struct {
	Name     string
	Password string
	Salt     string
	Callback string
	Error    string
}

func (a authorize) GET(w http.ResponseWriter, r *http.Request) {
	if err := server.ParseURLVars(r, &a); err != nil {
		server.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	execute(w, a)
}

func (a authorize) POST(w http.ResponseWriter, r *http.Request) {
	if err := server.ParseFormVars(r, &a); err != nil {
		a.Error = err.Error()
		execute(w, a)
		return
	}

	db, err := mdb.GetConnection()
	if err != nil {
		a.Error = err.Error()
		execute(w, a)
		return
	}

	var ac account

	if err = orm.NewStmt(db, "account").
		Where("name='%v' and md5(concat(password, '%s')) = '%s'", a.Name, a.Salt, a.Password).
		Query(&ac); err != nil {
		a.Error = err.Error()
		log.Errorf("query account error:%v", err)
		execute(w, a)
		return
	}

	token := ac.token()

	w.Header().Add("X-Authorize-Token", token)

	log.Debugf("user:%v, token:%v", a.Name, token)

	w.Header().Add("Location", a.Callback)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
