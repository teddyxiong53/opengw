package setting

import (
	"fmt"
	"testing"
)

func TestNTPAddHost(t *testing.T){
	//跳过本函数
	//t.SkipNow()

	NTPAddHost("ntp1.aliyun.com")
	NTPAddHost("ntp2.aliyun.com")

	NTPRemoveHost("ntp1.aliyun.com")

	hostArray := NTPGetHost()
	fmt.Printf("NTPAddr %+v\n",hostArray)

	NTPGetTime()
}


