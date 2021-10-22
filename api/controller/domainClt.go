package controller

import (
	"encoding/json"
	"fmt"
	"interview1710/api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// return all ctgs
func GetAllDomain(w http.ResponseWriter, r *http.Request) {

	domains, err := models.AllDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(domains)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func DeleteOneDomain(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	err := models.DeleteDomain(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s\n", "Deleted domain.id "+id)
}

func CreateDomain(w http.ResponseWriter, r *http.Request) {

	domains, err := models.CreateDomain(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(domains)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func UpdateDomain(w http.ResponseWriter, r *http.Request) {

	domains, err := models.UpdateDomain(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(domains)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func DomainBasedOnCatId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	limitStr := vars["limit"]
	limitNum, _ := strconv.Atoi(limitStr)

	offsetStr := vars["offset"]
	offsetNum, _ := strconv.Atoi(offsetStr)

	domains, err := models.DomainCategory(id, limitNum, offsetNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(domains)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func SearchDomain(w http.ResponseWriter, r *http.Request) {

	result, err := models.SearchFullText(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}
