//Whois will concact the RIPE.net and get from its APIs the ISP of the ip address passed as argumets
package main

import (
	"fmt"
	"os"
	"sync"
	"whois/returnisp"
)

//Functuion runner permits to get the rusults concorrently
func runner(ipaddr string, wg *sync.WaitGroup) {
	//recuper il ISP dell' ip
	isp, country := returnisp.ReturnISPandStore(ipaddr)
	//scrive a video l'ip e il sup ISP
	fmt.Println(ipaddr, isp, country)
	//decrementa di una unità il contatore del waitgroup wg
	wg.Done()
}

func main() {
	ip := os.Args[1:]
	//fmt.Println(ip)

	//Creates a wait group to manage Goroutines
	var wg sync.WaitGroup

	//For each ip passed as argument a Goroutin of func runner is created
	for _, ipaddr := range ip {
		//Adds one counter to the waitgroup for each goroutine created
		wg.Add(1)
		//Fa partire una goroutine a cui è collegato il pointer alla variabile wg per modificarne il contenuto
		go runner(ipaddr, &wg)
	}
	//waits for all Goroutines created to finish before exiting
	wg.Wait()
	os.Exit(0)
}
