package store

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type setRequest struct {
	Value string `json:"value"`
}

// Set  godoc
//
//	@Tags			Store
//	@Summary		Set
//	@Description	set value for key
//	@Param			key		path	string	true	"storage key"
//	@Param 			request 		body 	setRequest true "new value"
//	@Success		200
//	@Failure		400,500
//	@Router			/store/{key}	[put]
func (i *Implementation) Set(w http.ResponseWriter, r *http.Request) {
	key, err := getKeyFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req setRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		i.log.Error(fmt.Sprintf("set: decode request: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = i.store.Set(key, req.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
