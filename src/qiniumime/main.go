package main

import (
	"flag"
	"io"
	"log"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/rsf"
)

var ak *string
var sk *string

func changeMime(bucket, old, _new string, list []rsf.ListItem) {
	client := rs.New(nil)
	for _, v := range list {
		if v.MimeType == old {
			err := client.ChangeMime(nil, bucket, v.Key, _new)
			if err != nil {
				log.Println(v.Key, "modify mime fail", err)
			} else {
				log.Println(v.Key, "modified")
			}
		}
	}
}

func main() {
	ak = flag.String("ak", "", "access key")
	sk = flag.String("sk", "", "secret key")
	bucket := flag.String("bucket", "", "bucket")
	oldMime := flag.String("old", "", "old mime")
	newMime := flag.String("new", "", "new mime")
	flag.Parse()
	if *ak == "" || *sk == "" || *bucket == "" || *oldMime == "" || *newMime == "" {
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
	log.Println("total files", len(list))
	changeMime(*bucket, *oldMime, *newMime, list)
}
