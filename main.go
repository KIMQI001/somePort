package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)
var addr = flag.String("addr", "localhost:8081", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func main() {
	//get the data
	FIL:=GetFILPrice()
	fmt.Printf("%s",FIL)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	//http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))

}


func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		if string(message) == "futures"{
			FILPrice:= GetFILPrice()
			err = c.WriteMessage(mt, []byte(FILPrice))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		//if string(message) == "time"{
		//	CuTime,_:= GetCurrentTime()
		//	err = c.WriteMessage(mt, []byte(CuTime))
		//	if err != nil {
		//		log.Println("write:", err)
		//		break
		//	}
		//}

		//if string(message) == "minute"{
		//	_,_,CuTime:= GetCurrentTime()
		//	err = c.WriteMessage(mt, []byte(CuTime[2]))
		//	if err != nil {
		//		log.Println("write:", err)
		//		break
		//	}
		//}

	}
}

//获取编码格式
func determinEncodeing(r io.Reader) encoding.Encoding {
	bytes,err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e,_,_ := charset.DetermineEncoding(bytes,"")
	return e
}
func GetFILPrice() string {
	// get the data
	resp,err := http.Get("https://mifengcha.com/coin/filecoin?browsermode=pc")
	if err!= nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		fmt.Println("Status is not OK",resp.StatusCode)
		return ""
	}

	e := determinEncodeing(resp.Body)
	//println(resp.Body)
	utf8Obj := transform.NewReader(resp.Body,e.NewDecoder())
	all,err := ioutil.ReadAll(utf8Obj)
	if err!=nil{
		panic(nil)
	}
	//println(string(all))
	match:= regexp.MustCompile(`<span data-v-53a9ede8>([\d])(.)([\d]{4})`)
	FILStr := match.FindAllString(string(all),-1)
	matchFIL:= regexp.MustCompile(`>..+[\d]`)

	FStr := matchFIL.FindAllString(FILStr[2],-1)
	matchF:= regexp.MustCompile(`[\d].....`)
	FIL := matchF.FindAllString(FStr[0],-1)
	return FIL[0]
	//return ""
}

//func GetCurrentTime() (string,string,[]string){
//
//	CurrentTime :=fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
//
//
//	TimeStampInt:=[]int{0,0,0}
//	monthStr:=time.Now().Month()
//	monthS:=fmt.Sprintf("%d",monthStr)
//	monthI,_:=strconv.Atoi(monthS)
//	var month string
//	if monthI< 10{
//		month=fmt.Sprintf("0%s",monthS)
//	}else {
//		month=monthS
//	}
//	TimeStampInt[0]=time.Now().Day()
//	TimeStampInt[1]=time.Now().Hour()
//	TimeStampInt[2]=time.Now().Minute()
//
//	TimeStamp:=[]string{"","",""}
//	for i,_ :=range TimeStampInt{
//		if TimeStampInt[i]<10{
//			TimeStamp[i]=fmt.Sprintf("0%s",strconv.Itoa(TimeStampInt[i]))
//		}else {
//			TimeStamp[i]=fmt.Sprintf("%s",strconv.Itoa(TimeStampInt[i]))
//		}
//	}
//
//	return CurrentTime,month,TimeStamp
//}