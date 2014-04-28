package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rsf"
)

func main() {
	ak := flag.String("ak", "", "access key")
	sk := flag.String("sk", "", "secret key")
	bucket := flag.String("bucket", "", "bucket")
	out := flag.String("o", "", "output file")
	flag.Parse()
	if ak == nil || sk == nil || bucket == nil || out == nil {
		log.Fatalln("invalid args")
		return
	}
	conf.ACCESS_KEY = *ak
	conf.SECRET_KEY = *sk
	list := []rsf.ListItem{}
	client := rsf.New(nil)
	marker := ""
	var err error
	for {
		var ret []rsf.ListItem
		ret, marker, err = client.ListPrefix(nil, *bucket, "", marker, 1000)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error occured!", err)
				break
			}
		}
		list = append(list, ret...)
	}
	data, err := json.MarshalIndent(list, "", "")
	if err != nil {
		log.Fatalln("list data error", err)
		return
	}
	err = ioutil.WriteFile(*out, data, 0)
	if err != nil {
		log.Fatalln("write file failed")
	}
}
