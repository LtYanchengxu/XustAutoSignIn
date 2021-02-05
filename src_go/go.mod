module XustAutoSignIn

go 1.13

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/image => github.com/golang/image v0.0.0-20201208152932-35266b937fa6
	golang.org/x/net => github.com/golang/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sync => github.com/golang/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys => github.com/golang/sys v0.0.0-20210113181707-4bcb84eeeb78
	golang.org/x/tools => github.com/golang/tools v0.0.0-20210114065538-d78b04bdf963
)

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/robfig/cron/v3 v3.0.1
	go.uber.org/zap v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.0.3
	gorm.io/driver/postgres v1.0.6
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.11
)
