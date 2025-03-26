package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/magicznykacpur/psst-backend/internal/database"
)

type createMessageRequest struct {
	ChatID string `json:"chat_id"`
	Body   string `json:"body"`
}

func (cfg *ApiConfig) HandlerCreateMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var msgRequest createMessageRequest
	err = json.Unmarshal(requestBytes, &msgRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if msgRequest.ChatID == "" || msgRequest.Body == "" {
		respondWithError(w, http.StatusBadRequest, "message must contain chat id and body")
		return
	}

	chatID, err := uuid.Parse(msgRequest.ChatID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chat id malformed")
		return
	}

	chat, err := cfg.DB.GetChatById(r.Context(), chatID)
	if err != nil && strings.Contains(err.Error(), "sql: no rows in result set") {
		respondWithError(w, http.StatusBadRequest, "chat doesn't exist")
		return
	}

	message, err := cfg.DB.CreateMessage(r.Context(),
		database.CreateMessageParams{
			ChatID: chat.ID,
			Body:   msgRequest.Body,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt send message to chat")
		return
	}

	log.Printf("message %s sent to chat %s...", message.ID.String(), chat.ID.String())

	w.WriteHeader(http.StatusCreated)
}

type messageResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

func (cfg *ApiConfig) HandlerGetAllMessagesFromChat(w http.ResponseWriter, r *http.Request) {
	chatId, err := uuid.Parse(r.PathValue("chat_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chat id malformed")
		return
	}

	userId, err := uuid.Parse(r.Header.Get("User-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user id malformed")
		return
	}

	chat, err := cfg.DB.GetChatById(r.Context(), chatId)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusNotFound, "chat not found")
		return
	}

	if chat.SenderID != userId {
		respondWithError(w, http.StatusForbidden, "chat doesn't belong to user")
		return
	}

	messages, err := cfg.DB.GetMessagesByChatId(r.Context(), chatId)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusNotFound, "no messages found for this chat")
		return
	}

	msgsResponse := []messageResponse{}
	for _, msg := range messages {
		msgsResponse = append(msgsResponse,
			messageResponse{
				ID:        msg.ID,
				CreatedAt: msg.CreatedAt,
				Body:      msg.Body,
			},
		)
	}

	log.Printf("retrieved all messages from chat %s...", chat.ID.String())

	respondWithJSON(w, http.StatusOK, msgsResponse)
}
