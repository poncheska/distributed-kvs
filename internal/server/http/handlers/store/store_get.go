package store

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type getResponse struct {
	Value string `json:"value"`
}

// Get  godoc
//
//	@Tags			Store
//	@Summary		Get
//	@Description	get value by key
//	@Produce		json
//	@Param			key		path	string	true	"storage key"
//	@Success		200			{object}	getResponse
//	@Failure		400,500
//	@Router			/store/{key}	[get]
func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
	key, err := getKeyFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := i.store.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(getResponse{Value: value})
	if err != nil {
		i.log.Error(fmt.Sprintf("get: encode response: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
