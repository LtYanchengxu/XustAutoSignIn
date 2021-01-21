package router

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"XustAutoSignIn/config"
	"XustAutoSignIn/utils"
	"XustAutoSignIn/variable"
)

var runServer sync.Once

func Run() {
	runServer.Do(func() {
		cfg := config.GetConfig()
		// cron start
		variable.Cron.AddFunc("0 45 11 * * ?", utils.XustSignInAfternoon)
		variable.Cron.AddFunc("0 45 17 * * ?", utils.XustSignInNight)
		variable.Cron.Start()
		// init server
		gin.SetMode(gin.ReleaseMode)
		variable.Engine = gin.Default()
		// middleware
		// api
		initApiRouter()
		// run
		variable.Engine.Run(":" + strconv.Itoa(cfg.GetServerPort()))
	})
}
