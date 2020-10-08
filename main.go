package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Parameter struct {
	Panjang int `json:"panjang"`
	Lebar   int `json:"lebar"`
	Tinggi  int `json:"tinggi"`
}

type Hasil struct {
	Bangun string `json:"bangun"`
	Volume int    `json:"volume"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hitung-volume", Volume)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func (p *Parameter) RumusVolume() int {
	return p.Panjang * p.Lebar * p.Tinggi
}

func (p *Parameter) cekbangun() string {
	if p.Panjang == p.Lebar && p.Lebar == p.Tinggi {
		return "kubus"
	} else {
		return "balok"
	}
}

func Volume(w http.ResponseWriter, r *http.Request) {

	var hasilHitung []Hasil
	var parameter []Parameter
	if r.Method != "POST" {
		WarpAPIError(w, r, "cannot read body", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		WarpAPIError(w, r, "cannot read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &parameter)
	if err != nil {
		WarpAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range parameter {
		hasilHitung = append(hasilHitung, Hasil{
			Bangun: v.cekbangun(),
			Volume: v.RumusVolume(),
		})
	}

	WrapAPIData(w, r, hasilHitung, http.StatusOK, "success")

}

func WarpAPIError(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":           code,
		"error_type":     http.StatusText(code),
		"errror_details": message,
	})
	if err == nil {
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't warp API error : $p", err))
	}
}

func WarpAPISuccess(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't warp API error : $p", err))
	}
}

func WrapAPIData(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
		"data":   data,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't warp API error : $p", err))
	}
}
