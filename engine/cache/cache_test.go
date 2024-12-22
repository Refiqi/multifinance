package cache

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDoubleBuffer(t *testing.T) {
	Convey("double buffer", t, func() {
		doubleBufferCacheLru := NewDoubleBufferLru(DoubleBufferLruConfig{
			CacheSize:        3,
			CacheExpiryMSec:  10000,
			CacheRefreshMSec: 3000,
		}).(*DoubleBuffer)

		//Convey("fill cache", func() {
		value, err := doubleBufferCacheLru.Get("satu", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "satu", 1)
			return 1, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 1)

		value, err = doubleBufferCacheLru.Get("dua", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "dua", 2)
			return 2, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 2)

		value, err = doubleBufferCacheLru.Get("tiga", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "tiga", 3)
			return 3, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 3)
		//})

		//Convey("get from cache", func() {
		value, err = doubleBufferCacheLru.Get("satu", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "satu", 4)
			return 4, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 1)

		value, err = doubleBufferCacheLru.Get("dua", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "dua", 5)
			return 5, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 2)

		value, err = doubleBufferCacheLru.Get("tiga", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "tiga", 6)
			return 6, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 3)
		//})

		//Convey("refresh expired", func() {
		time.Sleep(4 * time.Second)

		value, err = doubleBufferCacheLru.Get("satu", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "satu", 7)
			return 7, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 1)

		value, err = doubleBufferCacheLru.Get("dua", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "dua", 8)
			return 8, nil
		})

		So(err, ShouldBeNil)
		So(value, ShouldEqual, 2)

		value, err = doubleBufferCacheLru.Get("tiga", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "tiga", 9)
			return 9, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 3)
		//})

		// refreshed cache
		time.Sleep(1 * time.Second)

		value, err = doubleBufferCacheLru.Get("satu", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "satu", 10)
			return 10, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 7)

		value, err = doubleBufferCacheLru.Get("dua", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "dua", 11)
			return 11, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 8)

		value, err = doubleBufferCacheLru.Get("tiga", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "tiga", 12)
			return 12, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 9)

		// expired cache
		time.Sleep(11 * time.Second)

		value, err = doubleBufferCacheLru.Get("satu", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "satu", nil)
			return nil, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldBeNil)

		value, err = doubleBufferCacheLru.Get("dua", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "dua", 13)
			return 13, nil
		})
		So(err, ShouldBeNil)
		So(value, ShouldEqual, 13)

		value, err = doubleBufferCacheLru.Get("tiga", func() (data interface{}, err error) {
			fmt.Printf("cache refresh %v, %v\n", "tiga", "err")
			return nil, fmt.Errorf("error")
		})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "error")
		So(value, ShouldEqual, nil)
		//})
	})
}