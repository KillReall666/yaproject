package hashmiddleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
)

func NewHashMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hashFunc := func(w http.ResponseWriter, r *http.Request) {
			if key != "" && strings.Contains(r.URL.Path, "/updates/") {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					log.Println("err when reading body", err)
				}

				hash := r.Header.Get("HashSHA256")

				if !verifySHA256Hash(body, hash, key) {
					http.Error(w, "несовпадение хэша", http.StatusBadRequest)
					log.Println("ошибка: хэш не совпадает")
					return
				}

				r.Body = io.NopCloser(bytes.NewReader(body))
				w.Header().Set("HashSHA256", hash)
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hashFunc)
	}
}

func verifySHA256Hash(data []byte, hash string, key string) bool {
	computedHash := computeSHA256Hash(data, key)
	return computedHash == hash
}

func computeSHA256Hash(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
