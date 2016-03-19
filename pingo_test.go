package pingo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResolve(t *testing.T) {
	Convey("Call resolve with", t, func() {

		Convey("a valid IPv4 address,", func() {
			target := "127.0.0.1"
			address, err := resolve(target, IP)
			So(err, ShouldBeNil)
			So(address.String(), ShouldEqual, target)
		})

		Convey("a valid IPv6 address,", func() {
			target := "::1"
			address, err := resolve(target, IP)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, target)
		})

		Convey("a valid hostname (for IPv4)", func() {
			target := "localhost"
			address, err := resolve(target, IPv4)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, "127.0.0.1")
		})

		Convey("a valid hostname (for IPv6)", func() {
			target := "localhost"
			address, err := resolve(target, IPv6)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, "::1")
		})

		Convey("a valid hostname (for any kind of IP)", func() {
			target := "localhost"
			address, err := resolve(target, IP)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
		})

		Convey("an invalid IP address", func() {
			target := "111.222.333.444"
			_, err := resolve(target, IP)
			So(err, ShouldNotBeNil)
		})
		Convey("an invalid hostname", func() {
			target := "dummy.host.nowhere"
			_, err := resolve(target, IP)
			So(err, ShouldNotBeNil)
		})
		Convey("an invalid IP version", func() {
			target := "localhost"
			_, err := resolve(target, "lol")
			So(err, ShouldNotBeNil)
		})
	})

}
