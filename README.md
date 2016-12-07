# pingo

pingo a a Go package to do ICMP pings (send ICMP echo and wait for a reply). You can do it with IPv4 and IPv6.

**WARNING : to use pingo you need root privileges.**

# Install

```
go get github.com/samonzeweb/pingo
```

# Example

```
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/samonzeweb/pingo"
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

The pingo package depends on [GoConvey](https://github.com/smartystreets/goconvey). Simply install all the necessary things with :

```
go get -t
```

The tests are naives. Testing some cases need disabling network card, changing routes, ... It's really hard to automate but easy to do manually if needed.

# Licence

Released under the MIT License, see `LICENSE.txt` for more informations.
