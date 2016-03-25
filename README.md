# pingo

pingo a a Go package to do ICMP pings (send ICMP echo and wait for a reply). You can do it with IPv4 and IPv6.

To use pingo you need root privileges.

# Install

```
go get bitbucket.org/samonzeweb/pingo
```

# Example

```
package main

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/samonzeweb/pingo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give a hostname as argument")
		os.Exit(1)
	}

	t, err := pingo.SimplePing(os.Args[1], pingo.IP, time.Second)
	if err != nil {
		if err == pingo.ErrTimeOut {
			fmt.Printf("Time out : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Error : %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Got a response from %s, in %d ms\n", os.Args[1], t/time.Millisecond)
}
```

# Development

The pingo package depends on external packages ( [golang.org/x/net...](https://godoc.org/golang.org/x/net) ) you have to install yourself if you cloned the git repository or if you need to run tests :

```
go get
go get github.com/smartystreets/goconvey
```

# Licence

Released under the MIT License, see `LICENSE.txt` for more informations.
