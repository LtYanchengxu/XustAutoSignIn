package flags

import (
	"flag"
	"sync"
)

type Flags struct {
	Help    bool
	Version bool
	Test    bool
	Config  string
}

var (
	flags     *Flags
	flagsOnce sync.Once
)

func init() {
	flags = new(Flags)
	flag.BoolVar(&flags.Help, "h", false, "help")
	flag.BoolVar(&flags.Version, "V", false, "version")
	flag.BoolVar(&flags.Test, "t", false, "test a config file")
	flag.StringVar(&flags.Config, "config", "", "config file")
}

func GetFlags() *Flags {
	flagsOnce.Do(func() {
		flag.Parse()
		if flags.Help == true {
			printUsage()
		}
		if flags.Version == true {
			printVersion()
		}
	})
	return flags
}

func (this *Flags) GetConfigFile() string {
	return this.Config
}

func (this *Flags) CheckTestFlag() {
	if this.Test == true {
		printTest()
	}
}
