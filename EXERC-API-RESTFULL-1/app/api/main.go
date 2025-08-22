package main

import (
	"fmt"      //nativo
	"net/http" //nativo
	"strings"  //nativo
)

func getOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Este endpoint aceita somente GET")
}

func postHandler(w http.ResponseWriter, response *http.Request) {
	fmt.Fprintln(w, "POST autorizado! Você tem permissão.")
}

// ///////// gin-gonose alternativa
func withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token não informado", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "meuToken123" {
			http.Error(w, "Token inválido", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func main() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")

	})

	_, erros := fmt.Println("ola")
	fmt.Print(erros)
	http.HandleFunc("/get-only", getOnlyHandler)

	//middleware
	http.HandleFunc("/post", withAuth(postHandler))

	fmt.Println("Servidor Rodando na porta 7878")
	http.ListenAndServe(":7878", nil)
}

/////number, varchar, number
//500 usuarios simultaneo
//50 iniciais , incremento de 10, 200cap max simultaneas
//1000 dados do banco
