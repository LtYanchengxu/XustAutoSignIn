package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"XustAutoSignIn/config"
	"XustAutoSignIn/utils"
	"XustAutoSignIn/variable"
)

var runServer sync.Once

const (
	afternoonSignInCronString = "0 30 11 * * ?"
	nightSignInCronString     = "0 30 17 * * ?"
)

func Run() {
	runServer.Do(func() {
		cfg := config.GetConfig()
		// init server
		binding.Validator = new(ginValidator) // update validator
		gin.SetMode(gin.ReleaseMode)
		variable.Engine = gin.Default()
		// middleware
		// api and cron
		baseApiRouter()
		if cfg.XustAutoSignIn {
			xustAutoSignInApi()
			variable.Cron.AddFunc(afternoonSignInCronString, utils.XustSignInAfternoon)
			variable.Cron.AddFunc(nightSignInCronString, utils.XustSignInNight)
		}
		// run
		variable.Cron.Start()
		variable.Engine.Run(":" + strconv.Itoa(cfg.GetServerPort()))
	})
}
