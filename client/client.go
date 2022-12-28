package main

/* TODO:
O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

*/

import "fmt"

func main() {
	fmt.Println("client on")
}
