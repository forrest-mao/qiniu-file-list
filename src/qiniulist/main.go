package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rsf"
)

func main() {
	ak := flag.String("ak", "", "access key")
	sk := flag.String("sk", "", "secret key")
	bucket := flag.String("bucket", "", "bucket")
	out := flag.String("o", "", "output file dir")
	flag.Parse()
	if ak == nil || sk == nil || bucket == nil || out == nil {
		log.Fatalln("invalid args")
		return
	}
	conf.ACCESS_KEY = *ak
	conf.SECRET_KEY = *sk

	client := rsf.New(nil)
	marker = ""
	err := os.MkdirAll(*out, 0)
	if err != nil {
		log.Fatalln("invalid out dir")
	}
	count := 0
	for {
		ret, marker, err := client.ListPrefix(nil, bucket, "", marker, 1000)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error occured!")
				return
			}
		}
		_path := path.Join(*out, "part"+strconv.Itoa(count))
		ioutil.WriteFile(_path, ret, 0)
	}
}
