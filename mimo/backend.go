package mimo

import (
	"fmt"
	"log"
	"net/url"
)

type Backend struct {
	Url    string `json:url`
	Master bool   `json:"master"`
	Weight int
}

func (back *Backend) init() error {
	u, err := url.Parse(back.Url)
	eMsg := "backend url [" + back.Url + "] "
	if err != nil {
		log.Println(eMsg+" parse failed:", err)
		return err
	}
	if u.Scheme != "http" {
		log.Println(eMsg + " ,schema is not http")
		return fmt.Errorf("schema is not http")
	}
	if back.Master {
		back.Weight = 1
	}
	return nil
}

type Backends []*Backend

func (backs *Backends) init() error {
	masterTotal := 0
	for _, bak := range *backs {
		err := bak.init()
		if err != nil {
			continue
		}
		if bak.Master {
			masterTotal++
		}
	}

	if masterTotal == 0 {
		for _, bak := range *backs {
			bak.Weight = 1
		}
	}
	return nil
}

func (backs *Backends) GetMasterIndex() int {
	masterTotal := 0
	for _, bak := range *backs {
		if bak.Weight > 0 {
			masterTotal++
		}
	}
	if masterTotal < 1 {
		return 0
	}

	indexM := randR.Int() % masterTotal
	i := 0
	for n, bak := range *backs {
		if bak.Weight > 0 && i >= indexM {
			i++
			return n
		}
	}
	return 0

}