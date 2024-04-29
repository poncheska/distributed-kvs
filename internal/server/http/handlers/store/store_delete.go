package store

import (
	"errors"
	"net/http"

	storeerrors "distributed-kvs/internal/store/errors"
)

// Delete  godoc
//
//	@Tags			Store
//	@Summary		Delete
//	@Description	delete value by key
//	@Param			key		path	string	true	"storage key"
//	@Success		204
//	@Failure		400,500
//	@Router			/store/{key}	[delete]
func (i *Implementation) Delete(w http.ResponseWriter, r *http.Request) {
	key, err := getKeyFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = i.store.Delete(key)
	if err != nil {
		if errors.Is(err, storeerrors.ErrActionUnavailable) {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
