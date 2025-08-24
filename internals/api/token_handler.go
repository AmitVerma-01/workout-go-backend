package api

import (
	"encoding/json"
	"go_beginner/internals/store"
	"go_beginner/internals/tokens"
	"go_beginner/utils"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func  (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"Message": "Invalid request body",
		})
		return
	}

	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		h.logger.Printf("Error getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"Message": "Failed to get user",
		})
		return
	}
	println("Password",req.Password)
	passwordMatch , err := user.PasswordHash.Check(req.Password)
	if err != nil {
		h.logger.Printf("Error : password match, %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"Message": "Failed to check password",
		})
		return
	}
	if !passwordMatch {
		h.logger.Printf("Password mismatch for user %s", req.Email)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"Message": "Invalid email or password",
		})
		return
	}

	token, err := h.tokenStore.CrateNewToken(int(user.ID), 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("Error creating new token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"Message": "Failed to create token",
		})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"token": token,
	})
}

