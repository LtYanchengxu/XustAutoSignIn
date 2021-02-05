package variable

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Validate *validator.Validate
	DB       *gorm.DB
	Cache    *cache.Cache
	Logger   *zap.SugaredLogger
	Cron     *cron.Cron
	Engine   *gin.Engine
)

func init() {
	// init validator
	Validate = validator.New()
	// init Cache
	Cache = cache.New(5*time.Minute, 10*time.Minute)
	// init logger
	Logger = NewConsoleLogger()
	// init cron
	Cron = cron.New(cron.WithSeconds())
}
