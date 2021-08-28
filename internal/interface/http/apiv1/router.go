package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type Register interface {
	Register(r *gin.RouterGroup)
}

func register(r *gin.RouterGroup, register ...Register) {
	for _, v := range register {
		v.Register(r)
	}
}

// Swagger spec:
// @title       云猫api文档
// @version     1.0
// @BasePath    /api/v1

func NewRouter(r *gin.Engine, db *database.Database, kv kvdb.KvDb) error {

	v1 := r.Group("/api/v1")

	system := NewSystem(kv)

	register(v1, system)

	return nil
}
