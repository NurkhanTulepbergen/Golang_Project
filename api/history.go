package api

import (
	"Golang_Project/pkg/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Handler to add history for a user
func (api *API) AddHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var history model.History
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = api.HistoryModel.AddHistory(history.UserID, history.UserName, history.OrdersList)
	if err != nil {
		http.Error(w, "Failed to add history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler to get history for a user
func (api *API) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := &model.HistoryFilter{}
	filter.UserID, _ = strconv.Atoi(r.URL.Query().Get("user_id"))
	filter.SortBy = r.URL.Query().Get("sort_by")
	filter.Order = r.URL.Query().Get("order")
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))

	// Get history from the model
	historyList, err := api.HistoryModel.GetHistory(filter)
	if err != nil {
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}

	// Encode the response as JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historyList)
}

// Handler to delete history for a user
func (api *API) DeleteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = api.HistoryModel.DeleteHistory(userID)
	if err != nil {
		http.Error(w, "Failed to delete history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handler to update history for a user
func (api *API) UpdateHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var history model.History
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = api.HistoryModel.UpdateHistory(history.UserID, history.UserName, history.OrdersList)
	if err != nil {
		http.Error(w, "Failed to update history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
