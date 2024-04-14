package store

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type joinRequest struct {
	NodeID string `json:"node_id"`
	Addr   string `json:"addr"`
}

// Join godoc
//
//	@Tags			Store
//	@Summary		Join
//	@Description	join distributed storage cluster
//	@Param			key		path	string	true	"storage key"
//	@Param 			request 		body 	joinRequest true "join storage params"
//	@Success		200
//	@Failure		400,500
//	@Router			/store	[post]
func (i *Implementation) Join(w http.ResponseWriter, r *http.Request) {
	var req joinRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		i.log.Error(fmt.Sprintf("join: decode request: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = i.store.Join(req.NodeID, req.Addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
