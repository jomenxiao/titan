package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	pdAddr    string
	token     string
	namespace string
)

func main() {
	flag.StringVar(&pdAddr, "pd-addr", "", "pd address, example: tikv://10.18.47.111:2379,10.18.47.124:2379,10.18.47.113:2379")
	flag.StringVar(&token, "token", "", "client token")
	flag.StringVar(&namespace, "namespace", "", "biz name")
	flag.Parse()
	if token != "" && namespace == "" {
		namespace = strings.Split(token, "-")[0]
	}
	if namespace == "" {
		fmt.Printf("token is invalid or namespace is null\n")
		return
	}
	pdAddrList := strings.Split(pdAddr, ",")
	if len(pdAddrList) < 1 {
		fmt.Printf("pd address error")
		return
	}
	regionsInfoUrl := fmt.Sprintf("http://%s%s", pdAddrList[len(pdAddrList)-1], regionListApi)
	data, err := httpGet(regionsInfoUrl)
	if err != nil {
		fmt.Printf("http pd address %s error %v", regionsInfoUrl, err)
		return
	}
	size, regionIdList, err := getSize(namespace, data.Regions)
	if err != nil {
		fmt.Printf("get region size error %v", err)
	}

	fmt.Printf("pd : [ %s ] \n", pdAddr)
	fmt.Printf("namespace : %s \n", namespace)
	fmt.Printf("size : %d MB\n", size)
	fmt.Printf("region id list :\n %v \n", regionIdList)
}
