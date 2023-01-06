package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type RespCotacao struct {
	Bid float64 `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", handlerCotacao)
	http.ListenAndServe(":8081", nil)
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

	db, err := sql.Open("sqlite3", "./db/cambio.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}
	defer db.Close()

	err = salvar(db, cot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	cotResp, err := getLastSave(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err)
		return
	}

	cotJson, err := json.Marshal(cotResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(cotJson)
}

func salvar(db *sql.DB, cot *Cotacao) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	stmt, err := db.Prepare("insert into cotacao(name, bid, create_date) values ($,$,$)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	strBid := &cot.Usdbrl.Bid
	bid, err := strconv.ParseFloat(*strBid, 8)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, &cot.Usdbrl.Code, &cot.Usdbrl.Name, bid, &cot.Usdbrl)
	if err != nil {
		return err
	}

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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	// log.Println(resp.Body)

	resJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cot Cotacao
	err = json.Unmarshal(resJson, &cot)
	if err != nil {
		return nil, err
	}

	return &cot, nil

}

func getLastSave(db *sql.DB) (*RespCotacao, error) {

	rows, err := db.Query("select bid from cotacao where id = (select max(id) from cotacao")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var respCot RespCotacao

	for rows.Next() {

		err = rows.Scan(&respCot.Bid)
		if err != nil {
			return nil, err
		}
	}

	return &respCot, nil

}
