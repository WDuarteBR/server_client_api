package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

/* TODO:
 O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.

Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.

*/

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handlerCotacao)
	http.ListenAndServe(":8080", nil)
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {

	route := r.URL.Path
	if route != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cot, err := GetCotacao()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = salvar(cot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cotJson, err := json.Marshal(cot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(cotJson)

}

func salvar(cot *Cotacao) error {
	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	stmt, err := db.Prepare("insert into cota")

	return nil
}

func GetCotacao() (*Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	resp, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	var cot Cotacao
	err = json.Unmarshal(resp, &cot)
	if err != nil {
		return nil, err
	}

	return &cot, nil

}
