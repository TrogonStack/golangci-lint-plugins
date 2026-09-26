package wrongfile

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {} // want "wrongfile declares Handler in other.go, not wrongfile.go"
