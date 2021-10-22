package controller

import (
	"encoding/json"
	"fmt"
	"interview1710/api/cache"
	"interview1710/api/models"
	"net/http"
	"strconv"

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
	ctgC := cache.Get("category" + id)
	if ctgC == nil {
		ctg, err := models.GetOneCategory(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		cache.Set("category"+strconv.Itoa(int(ctg.ID)), ctg)
		uj, err := json.Marshal(ctg)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println(err)
		}
		fmt.Println("DB ")
		fmt.Fprintf(w, "%s \n", uj)
		return
	}
	fmt.Println("REDIS")
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
	//delete cache
	cache.Delete("category" + id)
	fmt.Fprintf(w, "%s\n", "Deleted Categoriy"+id)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {

	categories, err := models.CreateCategory(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// set new ctg in redisDB
	cache.Set("category"+strconv.Itoa(int(categories.ID)), categories)
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
