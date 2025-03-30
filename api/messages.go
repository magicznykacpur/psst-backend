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
	ChatID     string `json:"chat_id"`
	Body       string `json:"body"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
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

	senderID, err := uuid.Parse(msgRequest.SenderId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "sender id malformed")
		return
	}

	receiverID, err := uuid.Parse(msgRequest.ReceiverId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "receiver id malformed")
		return
	}

	message, err := cfg.DB.CreateMessage(r.Context(),
		database.CreateMessageParams{
			ChatID:     chat.ID,
			Body:       msgRequest.Body,
			SenderID:   senderID,
			ReceiverID: receiverID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't send message to chat")
		return
	}

	log.Printf("message %s sent to chat %s...", message.ID.String(), chat.ID.String())

	w.WriteHeader(http.StatusCreated)
}

type messageResponse struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Body       string    `json:"body"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
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

	if chat.SenderID != userId && chat.ReceiverID != userId {
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
				ID:         msg.ID,
				CreatedAt:  msg.CreatedAt,
				Body:       msg.Body,
				SenderID:   msg.SenderID,
				ReceiverID: msg.ReceiverID,
			},
		)
	}

	log.Printf("retrieved all messages from chat %s...", chat.ID.String())

	respondWithJSON(w, http.StatusOK, msgsResponse)
}

type deleteMessageRequest struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
}

func (cfg *ApiConfig) HandlerDeleteMessageFromChat(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.Header.Get("User-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user id malformed")
		return
	}

	defer r.Body.Close()

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't read request bytes")
		return
	}

	var deleteMsgReq deleteMessageRequest
	err = json.Unmarshal(reqBytes, &deleteMsgReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal request body")
		return
	}

	if deleteMsgReq.MessageID == (uuid.UUID{}) || deleteMsgReq.ChatID == (uuid.UUID{}) {
		respondWithError(w, http.StatusBadRequest, "request must contain message id and chat id")
		return
	}

	msg, err := cfg.DB.GetMessageWhereChatAndUser(r.Context(),
		database.GetMessageWhereChatAndUserParams{
			SenderID: userId,
			ChatID:   deleteMsgReq.ChatID,
		},
	)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusBadRequest, "message does not exist, or does not belong to user")
		return
	}
	if err != nil && !strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusBadRequest, "couldn't delete message")
		return
	}

	err = cfg.DB.DeleteMessage(r.Context(),
		database.DeleteMessageParams{
			ID:     msg.ID,
			ChatID: msg.ChatID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete message")
		return
	}

	log.Printf("deleted message %s from chat %s", deleteMsgReq.MessageID.String(), deleteMsgReq.ChatID.String())

	w.WriteHeader(http.StatusOK)
}
