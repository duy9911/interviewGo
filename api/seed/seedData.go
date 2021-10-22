package seed

import (
	"encoding/json"
	"fmt"
	"interview1710/api/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Custom struct to collect domain
type DataDomain struct {
	Datas Info `json:"data,omitempty"`
}

type Info struct {
	Data []Domain `json:"data,omitempty"`
}
type Domain struct {
	CategoryName string `json:"category1,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

// Custom struct to collect domain keyword
type DataKey struct {
	Data DomainInfor `json:"data,omitempty"`
}

type DomainInfor struct {
	DomainTag Tag `json:"data,omitempty"`
}
type Tag struct {
	Keywords []string `json:"tags,omitempty"`
}

// Custom categories struct
type Categories struct {
	Categories []models.Category `json:"data"`
}

func SeedData() {
	wg := sync.WaitGroup{}
	data := make(chan DataDomain, 82)
	wg.Add(2)
	scrapSiteCategories()
	go scrapSiteInfo(&wg, data)
	go scrapKeyWord(&wg, data)
	wg.Wait()
	log.Println("Successfully write to DB")

}

func scrapSiteCategories() {
	urlCategories := "https://localapi.trazk.com/webdata/websiteapicat.php?task=getCategories"
	body := httpGetData(urlCategories)

	categoriesJS := Categories{}
	jsonErr := json.Unmarshal(body, &categoriesJS)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for _, c := range categoriesJS.Categories {
		if _, err := models.CreateCategoryBasedName(c.Name); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Seed All Category Name > Done")
}

// Scrap keyword based on domain recieve and add to db
func scrapKeyWord(wg *sync.WaitGroup, data chan DataDomain) {
	defer wg.Done()
	for {
		v, ok := <-data
		if !ok {
			return
		}
		for _, d := range v.Datas.Data {
			domain := d.Domain
			fmt.Println(domain)
			urlKeyAnalysis := fmt.Sprintf("https://localapi.trazk.com/webdata/v3.1.php?task=getHeader&domain=" + domain + "&userToken=bUp3eG9lN3JsTjB0eDIrclRBeHl1NmszWlNhaGZ0WTZUQVpqZTM0SklIRT06OvMSvNfabZvlIf2yIVsJirg")
			body := httpGetData(urlKeyAnalysis)

			keywordJS := DataKey{}
			jsonErr := json.Unmarshal(body, &keywordJS)
			if jsonErr != nil {
				fmt.Println("Eror while unmarshal body keyword")
				continue
			}
			fmt.Println(keywordJS.Data.DomainTag.Keywords)
			fmt.Println(">>>>>>>>>>>", len(data))

			category, err := models.OneCategoryName(d.CategoryName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tags := keywordJS.Data.DomainTag.Keywords
			tagsJson, _ := json.Marshal(tags)

			site := models.SiteInfo{
				Domain:     d.Domain,
				CategoryID: int(category.ID),
				Tags:       tagsJson,
			}
			models.CreateDomainBasedCatName(site)
		}
	}
}

func scrapSiteInfo(wg *sync.WaitGroup, data chan DataDomain) {
	defer wg.Done()
	defer close(data)
	startNum := 0
	for {
		url := "https://localapi.trazk.com/webdata/v3.php?task=getTopWebsiteInVietnam&userToken=ZGdZVktsdE91by9qOUtndjc4MjYwTHdQeXllT3NKTS9ZUHVzdThJYTNWST06OhMNb7G48NOo6noCn1JFw0I&&draw=3&columns%5B0%5D%5Bdata%5D=function&columns%5B0%5D%5Bname%5D=&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=false&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=function&columns%5B1%5D%5Bname%5D=&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=false&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=countryTraffic&columns%5B2%5D%5Bname%5D=&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=false&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=uniqueUser&columns%5B3%5D%5Bname%5D=&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=false&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=time_on_site&columns%5B4%5D%5Bname%5D=&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=false&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=pages&columns%5B5%5D%5Bname%5D=&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=false&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B6%5D%5Bdata%5D=DesktopVsMobile&columns%5B6%5D%5Bname%5D=&columns%5B6%5D%5Bsearchable%5D=true&columns%5B6%5D%5Borderable%5D=false&columns%5B6%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B6%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B7%5D%5Bdata%5D=function&columns%5B7%5D%5Bname%5D=&columns%5B7%5D%5Bsearchable%5D=true&columns%5B7%5D%5Borderable%5D=false&columns%5B7%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B7%5D%5Bsearch%5D%5Bregex%5D=false&start=" + strconv.Itoa(startNum) + "&length=50&search%5Bvalue%5D=&search%5Bregex%5D=false&_=1634468481758"

		body := httpGetData(url)
		domainJS := DataDomain{}
		jsonErr := json.Unmarshal(body, &domainJS)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		if domainJS.Datas.Data == nil {
			fmt.Println("Fetch All Domain > Done")
			return
		}
		fmt.Println(domainJS)
		startNum += 50
		data <- domainJS
	}
}

func httpGetData(url string) []byte {

	expireTime := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "spacecount-tutorial")
	res, getErr := expireTime.Do(req)

	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}
