package core

import "fmt"

// KeeperHeaderName is the name of the custom Keeper header
const KeeperHeaderName = "X-KEEPER"

// KeeperVersion is the version of the build
var KeeperVersion = "undefined"

// KeeperHeaderValue is the value of the custom Keeper header
var KeeperHeaderValue = fmt.Sprintf("Version %s", KeeperVersion)

// KeeperUserAgent is the value of the user agent header sent to the backends
var KeeperUserAgent = fmt.Sprintf("Keeper Version %s", KeeperVersion)
