package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/models"
)

func (c *Routes) crudTokensListHandler(w http.ResponseWriter, r *http.Request) {
	var tokens []models.RegisterToken

	if err := c.db.Model(&models.RegisterToken{}).Find(&tokens).Error; err != nil {
		c.log.WithError(err).Error("unable to get tokens")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Routes) crudTokensAddHandler(w http.ResponseWriter, r *http.Request) {
	token := models.RegisterToken{}
	if err := c.db.Create(&token).Error; err != nil {
		c.log.WithError(err).Error("unable to get tokens")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		c.log.WithError(err).Error("unable to write builds to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
