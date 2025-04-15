package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MexcCredentialHandler struct {
	DB *gorm.DB
}

func NewMexcCredentialHandler(db *gorm.DB) *MexcCredentialHandler {
	return &MexcCredentialHandler{DB: db}
}

// POST /api/mexc-credentials
func (h *MexcCredentialHandler) AddCredential(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	var req struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
		Label     string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	encKey, err := crypto.Encrypt(req.ApiKey)
	if err != nil {
		http.Error(w, "Encryption error", http.StatusInternalServerError)
		return
	}
	encSecret, err := crypto.Encrypt(req.ApiSecret)
	if err != nil {
		http.Error(w, "Encryption error", http.StatusInternalServerError)
		return
	}
	cred := &entity.MexcApiCredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		ApiKey:    encKey,
		ApiSecret: encSecret,
		Label:     req.Label,
	}
	if err := h.DB.Create(cred).Error; err != nil {
		http.Error(w, "Failed to save credential", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GET /api/mexc-credentials
func (h *MexcCredentialHandler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	var creds []entity.MexcApiCredential
	if err := h.DB.Where("user_id = ?", userID).Find(&creds).Error; err != nil {
		http.Error(w, "Failed to fetch credentials", http.StatusInternalServerError)
		return
	}
	// Do not return secrets
	resp := make([]map[string]interface{}, 0, len(creds))
	for _, c := range creds {
		resp = append(resp, map[string]interface{}{
			"id":         c.ID,
			"label":      c.Label,
			"created_at": c.CreatedAt,
			"updated_at": c.UpdatedAt,
		})
	}
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/mexc-credentials/{id}
func (h *MexcCredentialHandler) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.MexcApiCredential{}).Error; err != nil {
		http.Error(w, "Failed to delete credential", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Helper: get user ID from context (replace with your actual auth logic)
func getUserIDFromContext(ctx context.Context) string {
	// Example: return ctx.Value("user_id").(string)
	return "user_id_placeholder"
}
