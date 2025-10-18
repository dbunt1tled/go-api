package server

import (
	"github.com/dbunt1tled/go-api/internal/util/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelmResponse struct {
	ID int64  `json:"id,omitempty" jsonapi:"primary,helm"`
	DB string `json:"db" jsonapi:"attr,db"`
}

func Helm(c echo.Context) error {
	return helper.JSONAPIModel(c.Response(), &HelmResponse{ID: 1, DB: "test"}, http.StatusOK)
}

// TODO Implement check db connection
// func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
//
// 	w.Header().Set("Content-Type", "application/json")
//
// 	if err := json.NewEncoder(w).Encode(db.Stats()); err != nil {
//
// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
//
// 		return
//
// 	}
//
// }
