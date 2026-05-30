package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

type RandomHandler struct{}

func NewRandomHandler(router *http.ServeMux) {
	random := &RandomHandler{}
	router.HandleFunc("/random", random.random)
}

func (h *RandomHandler) random (w http.ResponseWriter, r *http.Request){
	rand := strconv.Itoa(rand.Intn(6)+1)
	w.Write([]byte(rand))
}
