package controller

import (
	"encoding/json"
	"fmt"
	"interview1710/api/cache"
	"interview1710/api/models"
	"net/http"

	"github.com/gorilla/mux"
)

// return all ctgs
func GetAllCategory(w http.ResponseWriter, r *http.Request) {

	categories, err := models.AllCategory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctgC := cache.Get(id)
	fmt.Println("done")

	if ctgC == nil {
		ctg, err := models.GetOne(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cache.Set(ctg.ID, ctg)

		uj, err := json.Marshal(ctg)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println(err)
		}
		fmt.Fprintf(w, "%s \n", uj)

	}

	uj, err := json.Marshal(&ctgC)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%s \n", uj)
}

func DeleteOneCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	err := models.DeleteCategory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", "Deleted Categoriy"+id)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {

	categories, err := models.CreateCategory(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	categories, err := models.UpdateCategory(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uj, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s\n", uj)
}
