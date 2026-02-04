package api

import (
	"encoding/json"
	redisclient "internals/redis"
	"log"
	"net/http"
	"time"
)

type SetRequest struct {
	Key   string `json:"key"`
	Field string `json:"field"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"` // seconds
}

// POST /set
func SetHandler(w http.ResponseWriter, r *http.Request) {
	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Key == "" || req.Field == "" {
		http.Error(w, "key and field are required", http.StatusBadRequest)
		return
	}

	ttl := time.Duration(req.TTL) * time.Second
	if err := redisclient.SetKey(req.Key, req.Field, req.Value, ttl); err != nil {
		http.Error(w, "failed to set key", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"message": "hash field set successfully",
		"key":     req.Key,
		"field":   req.Field,
		"ttl":     req.TTL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


// GET /get?key=abc
func GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	val, err := redisclient.GetKey(key)
	ttl,err:=redisclient.GetTTL(key)

	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	j,err:=json.Marshal(val)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Println("TTL remaining:", ttl)
}
