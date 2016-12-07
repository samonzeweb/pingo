package pingo

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResolveIPv4(t *testing.T) {
	Convey("Calling resolve with", t, func() {

		Convey("a valid IPv4 address,", func() {
			target := "127.0.0.1"
			address, err := resolve(target, IP)
			So(err, ShouldBeNil)
			So(address.String(), ShouldEqual, target)
		})

		Convey("a valid hostname (for IPv4)", func() {
			target := "localhost"
			address, err := resolve(target, IPv4)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, "127.0.0.1")
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

func TestResolveIPv6(t *testing.T) {
	Convey("Calling resolve with", t, func() {

		if os.Getenv("IPV6_UNAVAILABLE") != "" {
			t.Skip("IPv6 is not available for the tests.")
		}

		Convey("a valid IPv6 address,", func() {
			target := "::1"
			address, err := resolve(target, IP)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, target)
		})

		Convey("a valid hostname (for IPv6)", func() {
			target := "ip6-localhost"
			address, err := resolve(target, IPv6)
			So(err, ShouldBeNil)
			So(address, ShouldNotBeNil)
			So(address.String(), ShouldEqual, "::1")
		})
	})
}

func TestSimplePingIPv4(t *testing.T) {
	Convey("Calling Ping with", t, func() {

		Convey("127.0.0.1 as target", func() {
			t, err := SimplePing("127.0.0.1", IPv4, time.Second)
			So(err, ShouldBeNil)
			So(t, ShouldBeLessThanOrEqualTo, time.Second)
		})

		Convey("8.8.8.8 as target", func() {
			t, err := SimplePing("8.8.8.8", IPv4, time.Second)
			So(err, ShouldBeNil)
			So(t, ShouldBeLessThanOrEqualTo, time.Second)
			So(t, ShouldBeGreaterThan, 0)
		})

		Convey("hostname and any IP as target", func() {
			t, err := SimplePing("www.debian.org", IP, time.Second)
			So(err, ShouldBeNil)
			So(t, ShouldBeLessThanOrEqualTo, time.Second)
			So(t, ShouldBeGreaterThan, 0)
		})

		Convey("hostname and IPv4 as target", func() {
			t, err := SimplePing("www.debian.org", IPv4, time.Second)
			So(err, ShouldBeNil)
			So(t, ShouldBeLessThanOrEqualTo, time.Second)
			So(t, ShouldBeGreaterThan, 0)
		})

	})
}

func TestSimplePingIPv6(t *testing.T) {
	Convey("Calling Ping with", t, func() {

		if os.Getenv("IPV6_UNAVAILABLE") != "" {
			t.Skip("IPv6 is not available for the tests.")
		}

		Convey("::1 as target", func() {
			if os.Getenv("TRAVIS_CI") != "" {
				t.Log("Disabled as Travis CI does not support ::1")
				t, err := SimplePing("::1", IPv6, time.Second)
				So(err, ShouldBeNil)
				So(t, ShouldBeLessThanOrEqualTo, time.Second)
			}
		})

		Convey("hostname and IPv6 as target", func() {
			t, err := SimplePing("www.debian.org", IPv6, time.Second)
			So(err, ShouldBeNil)
			So(t, ShouldBeLessThanOrEqualTo, time.Second)
			So(t, ShouldBeGreaterThan, 0)
		})
	})
}
