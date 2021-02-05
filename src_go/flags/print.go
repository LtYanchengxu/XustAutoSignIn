package flags

import (
	"flag"
	"fmt"
	"os"

	"XustAutoSignIn/variable"
)

func appNameString() string {
	return "XustAutoSignIn"
}

func versionString() string {
	return "0.0.1"
}

func copyrightString() string {
	return `BSD 3-Clause License
Copyright (c) 2020, Hang 'hang_c' Chen
          (c) 2021, XUST kcsoft`
}

func printUsage() {
	_, _ = fmt.Fprintf(os.Stdout, `%s %s

%s

Options:
`, appNameString(), versionString(), copyrightString())
	flag.PrintDefaults()
	os.Exit(variable.SUCCESS)
}

func printVersion() {
	_, _ = fmt.Fprintf(os.Stdout, `%s
Version: %s

%s
`, appNameString(), versionString(), copyrightString())
	os.Exit(variable.SUCCESS)
}

func printTest() {
	_, _ = fmt.Fprintf(os.Stdout, "Check Config File %s Successful\n", flags.Config)
	os.Exit(variable.SUCCESS)
}
