package instance

import "encoding/json"

type response struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}

func JsonResponse(status int, message any) []byte {
	con, err := json.Marshal(response{
		Status: status,
		Message: message,
	})
	if err != nil {
		return []byte("can't serve requests at the moment")
	}
	return con
}
