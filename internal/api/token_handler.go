package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DavidGudovic/api_exercise/internal/store"
	"github.com/DavidGudovic/api_exercise/internal/tokens"
	"github.com/DavidGudovic/api_exercise/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)

	if err != nil || user == nil {
		h.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to verify password"})
		return
	}

	if !passwordsDoMatch {
		_ = utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, tokens.ScopeAuth, 24*time.Hour)

	if err != nil {
		h.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJson(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create token"})
		return
	}

	_ = utils.WriteJson(w, http.StatusCreated, utils.Envelope{"token": token.Plaintext})
}
