package elasticDB

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"interview1710/api/models"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

type SearchConfigure struct {
	Text         string            `json:"text_search,omitempty"`
	Limit        int               `json:"limit,omitempty"`
	Omit         int               `json:"omit,omitempty"`
	TotalMatched int               `json:"totalMatch,omitempty"`
	DomainMatchs []models.SiteInfo `json:"domain_matched,omitempty"`
}

func NewElasticSearch() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(os.Getenv("ES_URL")),
		elastic.SetHealthcheckInterval(5*time.Second), // quit trying after 5 seconds
	)
	if err != nil {
		// (Bad Request): Failed to parse content to map if mapping bad
		fmt.Println("elastic.NewClient() ERROR: %v", err)
		log.Fatalf("quiting connection..")
	} else {
		// Print client information
		fmt.Println("client:", client)
		fmt.Println("client TYPE:", reflect.TypeOf(client), "\n")
	}
	return client
}

//fisrt index
func AddToEs() {
	ctx := context.Background()
	clientEs := NewElasticSearch()
	domains, errGetdomain := models.AllDomain()

	if errGetdomain != nil {
		// Handle error
		panic(errGetdomain)
	}
	for _, domain := range domains {

		esRespond, err := clientEs.Index().
			Index("domain").
			Type("keywords").
			Id(strconv.Itoa(int(domain.ID))).
			BodyJson(domain).
			Do(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Doaman %s Indexed with  %s to index %s \n", domain.Domain, esRespond.Id, esRespond.Index)
	}
	fmt.Println("Done...!")
}

//index when have a new domain
func AddOne(site models.SiteInfo) {
	ctx := context.Background()
	clientEs := NewElasticSearch()
	esRespond, err := clientEs.Index().
		Index("domain").
		Type("keywords").
		Id(strconv.Itoa(int(site.ID))).
		BodyJson(site).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Doaman %s Indexed with  %s to index %s \n", site.Domain, esRespond.Id, esRespond.Index)
	fmt.Println("Done...!")
}

func UpdateField(site models.SiteInfo) {
	ctx := context.Background()
	clientEs := NewElasticSearch()
	data := map[string]interface{}{
		"updated_at":  site.UpdatedAt,
		"domain":      site.Domain,
		"tags":        site.Tags,
		"category_id": site.CategoryID,
	}
	update, err := clientEs.
		Update().
		Index("domain").
		Type("keywords").
		Id("4").Doc(data).
		Do(ctx)

	if err != nil {
		errC := errors.New("update index " + fmt.Sprint(err))
		panic(errC)
	}
	fmt.Printf("New version of domain %q is now %d\n", update.Id, update.Version)
}

func SearchFullText(r *http.Request) (SearchConfigure, error) {
	ctx := context.Background()
	clientEs := NewElasticSearch()
	searchInfo, err := validateTextSearch(r)
	if err != nil {
		return searchInfo, err
	}
	tem2 := elastic.NewMatchQuery("tags", searchInfo.Text)
	searchResult, err := clientEs.Search().
		Index("domain"). // search in index "doamin"
		Query(tem2).     // specify the query
		From(searchInfo.Omit).Size(searchInfo.Limit).
		Do(ctx) // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Query with term {%v with  %v} took %d milliseconds\n", "tags", searchInfo.Text, searchResult.TookInMillis)
	var (
		ttyp models.SiteInfo
		site []models.SiteInfo
	)
	// detached each item in search result then print for each
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(models.SiteInfo); ok {
			fmt.Printf("Filed by %v: %s\n", t.Domain, t.Tags)
			site = append(site, t)
		}
	}
	// custom respond search
	results := SearchConfigure{
		Text:         searchInfo.Text,
		Limit:        searchInfo.Limit,
		Omit:         searchInfo.Omit,
		TotalMatched: searchInfo.TotalMatched,
		DomainMatchs: site,
	}

	// totalhist is another convenience function that works even when sths go wrong.
	fmt.Printf("Found a total of %d field\n", searchResult.TotalHits())
	return results, nil

}

func validateTextSearch(r *http.Request) (SearchConfigure, error) {
	searchCd := SearchConfigure{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&searchCd); err != nil {
		return searchCd, err
	}
	omit := &searchCd.Omit

	switch {
	case searchCd.Limit == 0:
		return searchCd, errors.New("limit is empty, check it")
	case omit == nil:
		return searchCd, errors.New("omit is empty, check it")
	case searchCd.Text == "":
		return searchCd, errors.New("text is empty, check it")
	}
	return searchCd, nil
}
