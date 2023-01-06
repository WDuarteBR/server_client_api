package main

/* TODO:
O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Cambio struct {
	Bid string
}

func main() {

	cambio, err := GetFromServer()
	if err != nil {
		panic(err)
	}

	err = SaveToFile(cambio)
	if err != nil {
		panic(err)
	}

	fmt.Println("Arquivo cotacao.txt criado com sucesso!")

}

func GetFromServer() (*Cambio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cam Cambio
	err = json.Unmarshal(respJson, &cam)
	if err != nil {
		return nil, err
	}

	return &cam, nil
}

func SaveToFile(c *Cambio) error {

	arq, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	defer arq.Close()

	_, err = arq.Write([]byte("Dólar: " + c.Bid))
	if err != nil {
		return err
	}

	return nil
}
