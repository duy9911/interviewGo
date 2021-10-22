package routers

import (
	"interview1710/api/controller"
	"interview1710/api/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/home", http.StatusSeeOther)

}
func HandleRequests() {
	r := mux.NewRouter()
	r.Use(middleware.WithLogging)
	r.HandleFunc("/", index)
	//ctg
	r.HandleFunc("/api/categories", middleware.SetMiddlewareJSON(controller.GetAllCategory)).Methods("GET")
	r.HandleFunc("/api/categories/{id:[0-9]+}", middleware.SetMiddlewareJSON(controller.GetCategoryById)).Methods("GET")
	r.HandleFunc("/api/categories/{id:[0-9]+}", middleware.SetMiddlewareJSON(controller.DeleteOneCategory)).Methods("DELETE")
	r.HandleFunc("/api/categories", middleware.SetMiddlewareJSON(controller.CreateCategory)).Methods("POST")
	r.HandleFunc("/api/categories", middleware.SetMiddlewareJSON(controller.UpdateCategory)).Methods("PUT")

	//domain
	r.HandleFunc("/api/domains", middleware.SetMiddlewareJSON(controller.GetAllDomain)).Methods("GET")
	r.HandleFunc("/api/domains/{id:[0-9]+}", middleware.SetMiddlewareJSON(controller.DeleteOneDomain)).Methods("DELETE")
	r.HandleFunc("/api/domains", middleware.SetMiddlewareJSON(controller.CreateDomain)).Methods("POST")
	r.HandleFunc("/api/domains", middleware.SetMiddlewareJSON(controller.UpdateDomain)).Methods("PUT")
	//search
	r.HandleFunc("/api/domains/search", middleware.SetMiddlewareJSON(controller.SearchDomain)).Methods("POST")

	//show domain from one ctgID
	r.Path("/api/categories/{id:[0-9]+}/domain").
		Queries("limit", "{limit}").Queries("offset", "{offset}").
		HandlerFunc(middleware.SetMiddlewareJSON(controller.DomainBasedOnCatId)).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9999", r))
}
