package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rsf"
)

type Item struct {
	Name     string    `json:"name"`
	Etag     string    `json:"etag"`
	Size     int64     `json:"size"`
	Time     time.Time `json:"time"`
	MimeType string    `json:"mime"`
	User     string    `json:"user"`
}

func convert(list []rsf.ListItem) (ret []Item, size int64) {
	for _, v := range list {
		var it Item
		it.Name = v.Key
		it.Etag = v.Hash
		it.Size = v.Fsize
		it.Time = time.Unix(v.PutTime/1e7, 0)
		it.MimeType = v.MimeType
		it.User = v.EndUser
		ret = append(ret, it)
		size += v.Fsize
	}
	return
}

func main() {
	ak := flag.String("ak", "", "access key")
	sk := flag.String("sk", "", "secret key")
	bucket := flag.String("bucket", "", "bucket")
	out := flag.String("o", "", "output file")
	flag.Parse()
	if *ak == "" || *sk == "" || *bucket == "" || *out == "" {
		flag.PrintDefaults()
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
		if err != nil && err != io.EOF {
			log.Fatalln("error occured!", err)
			return
		}
		list = append(list, ret...)
		if err == io.EOF {
			break
		}
	}
	items, size := convert(list)
	log.Println("total size", size)
	data, err := json.MarshalIndent(items, "", "")
	if err != nil {
		log.Fatalln("list data error", err)
		return
	}
	err = ioutil.WriteFile(*out, data, os.ModePerm)
	if err != nil {
		log.Fatalln("write file failed")
	}
}
