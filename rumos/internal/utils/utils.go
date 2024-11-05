package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Node struct {
	AirportID uint
	Price     uint
	Path      []string // Usado para demonstrar o caminho de IDs de voo
}

// PriorityQueue implementa heap.Interface e possui uma slice de *Node.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Nós com menor preço têm prioridade
	return pq[i].Price < pq[j].Price
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	// Adiciona um novo elemento à fila
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	// Remove e retorna o elemento de menor prioridade
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Filter is a generic function that filters a slice of items based on a predicate function.
// It returns a new slice containing only the elements for which the predicate function returns true.
//
// Parameters:
//   - items: The slice of items to be filtered.
//   - pred: A function that takes an item of type T and returns a boolean value.
//
// Return:
//   - A new slice containing only the elements for which the predicate function returns true.
func Filter[T any](items []T, pred func(T) bool) []T {
	var res []T
	for _, v := range items {
		if pred(v) {
			res = append(res, v)
		}
	}
	return res
}

// Find is a generic function that finds the first item in a slice that satisfies a predicate function.
// It returns a pointer to the found item, or nil if no item satisfies the predicate.
//
// Parameters:
//   - items: The slice of items to be searched.
//   - pred: A function that takes an item of type T and returns a boolean value.
//
// Return:
//   - A pointer to the first item in the slice that satisfies the predicate, or nil if no item satisfies the predicate.
func Find[T any](items []T, pred func(T) bool) *T {
	for _, v := range items {
		if pred(v) {
			return &v
		}
	}
	return nil
}

func PrintMap[V any](m map[string]V) string {
	str := "\n{"
	for k, v := range m {
		str += fmt.Sprintf("\n%v: %v", k, v)
	}
	return str + "\n}"
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
