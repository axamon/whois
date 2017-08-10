//Package returnisp restituisce l'Internet Service Provider dell'ip passato utilizzando le API del RIPE
package returnisp

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

var m sync.RWMutex
var gobfile = "ipisp.db"

//Verficifa se esiste il file di appoggio ipisp.db e se non esiste lo crea, inoltre instanzia la mappa listaipdafile
func init() {
	listaipdafile := make(map[string]string)
	if _, err := os.Stat(gobfile); os.IsNotExist(err) {
		file, err := os.Create(gobfile)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		// ipinit := "10.0.0.1"
		// ispinit := "private"
		// listaipdafile[ipinit] = ispinit
		//fmt.Println(listaipdafile[ipinit])
		//avvia il salvataggio della mappa sul file
		savemapgob(listaipdafile, &m)
	}
}

type mufile struct {
	gobfile string
	mu      sync.RWMutex
}

func savemapgob(data map[string]string, m *sync.RWMutex) {
	m.Lock()
	dataFile, err := os.Create(gobfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(data)
	m.Unlock()
	dataFile.Close()
}

func readmapfromgob(gobfile string) (res map[string]string) {
	var data map[string]string
	dataFile, err := os.Open(gobfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataFile.Close()
	return data
}

//ReturnISP restituisce l'Internet Servie Provider dell'ip passato utilizzando le API del RIPE
func ReturnISP(ip string) (isp, country string) {

	type ripe3 struct {
		Service struct {
			Name string
			Test []struct {
				Value string
			}
		}
		Parameters struct {
			Inverse struct{}
			Type    struct{}
			Flags   struct{}
			Queries struct {
				Query []struct {
					Value string
				}
			}
			Sources struct{}
		}
		Objects struct {
			Object []struct {
				Type string
				Link struct {
					Type string
					Href string
				}
				Source struct {
					ID string
				}
				Primary struct {
					Attribute []struct {
						Name  string
						Value string
					}
				}
				Attributes struct {
					Attribute []struct {
						Name  string
						Value string
					}
				}
			}
		}
	}

	resp, err := http.Get("http://rest.db.ripe.net/search.json?query-string=" + ip)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		risposta, err2 := ioutil.ReadAll(resp.Body)

		if err2 != nil {
			panic(err2)
		}

		var data ripe3
		err := json.Unmarshal([]byte(risposta), &data)
		if err != nil {
			panic(err)
		}

		for _, u := range data.Objects.Object {
			for _, l := range u.Attributes.Attribute {
				if l.Name == "netname" {
					//fmt.Println(l.Value)
					isp = l.Value
				}
				if l.Name == "country" {
					//fmt.Println(l.Value)
					country = l.Value
				}
			}
		}

	}
	return isp, country
}

//ReturnISPandStore recupera il ISP dell'ip e lo salva nel file di appoggio gobfile
func ReturnISPandStore(ip string) (isp, country string) {
	listaipdafile := make(map[string]string)

	listaipdafile = readmapfromgob(gobfile)

	//inizio := time.Now()
	//fmt.Print(ip)

	if ip == "author" {
		fmt.Println("Author: Alberto Bregliano, all rights reserved")
		os.Exit(0)
	}

	//verifiche che sia un IPv4 corretto
	testInput := net.ParseIP(ip)
	if testInput.To4() == nil {
		fmt.Printf("%v non è un IPv4 valido\n", ip)
		os.Exit(1)
	}

	//cerca nella mappa listaipdafile se c'è un valore associato all'ip, se è vuoto non c'è
	if len(listaipdafile[ip]) != 0 {
		isp = listaipdafile[ip]
		//fmt.Println(listaipdafile[ip])
		//fmt.Println(time.Since(inizio), "da file")
		return isp, country
	}

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	isp, country = ReturnISP(ip)
	//fmt.Println(isp)
	listaipdafile[ip] = isp
	savemapgob(listaipdafile, &m)
	//fmt.Println(time.Since(inizio), "da Api RIPE")
	// 	wg.Done()
	// }()
	// wg.Wait()
	return isp, country
}
