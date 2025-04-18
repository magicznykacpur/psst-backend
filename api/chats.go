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

type createChatRequest struct {
	ReceiverId string `json:"receiver_id"`
}

const noRows = "sql: no rows in result set"

func (cfg *ApiConfig) HandlerCreateChat(w http.ResponseWriter, r *http.Request) {
	senderId, err := uuid.Parse(r.Header.Get("User-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user id malformed")
		return
	}

	defer r.Body.Close()

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var chatRequest createChatRequest
	err = json.Unmarshal(requestBytes, &chatRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if chatRequest.ReceiverId == "" {
		respondWithError(w, http.StatusBadRequest, "missing chat receiver's id")
		return
	}

	receiverId, err := uuid.Parse(chatRequest.ReceiverId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "receivers's id malformed")
		return
	}

	receiver, err := cfg.DB.GetUserByID(r.Context(), receiverId)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusBadRequest, "receiver does not exist")
		return
	}

	chat, err := cfg.DB.CreateChatWith(r.Context(),
		database.CreateChatWithParams{
			SenderID:   senderId,
			ReceiverID: receiver.ID,
		},
	)
	if err != nil && strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
		respondWithError(w, http.StatusBadRequest, "chat with this user already exists")
		return
	}
	if err != nil && !strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chat")
		return
	}

	for client, _ := range cfg.Hub.Clients {
		if client.UserID == senderId || client.UserID == receiverId {
			broadcastFor := struct {
				Clients []uuid.UUID
				Message []byte
			}{
				[]uuid.UUID{senderId, receiverId},
				[]byte("new chat created"),
			}
			bytes, err := json.Marshal(broadcastFor)
			if err != nil {
				log.Printf("cannot marshall broadcast for: %v", err)
			}

			cfg.Hub.BroadcastFor <- bytes
		}
	}

	log.Printf("chat %s created between %s and %s...", chat.ID, senderId, receiverId)

	w.WriteHeader(http.StatusCreated)
}

type chatResponse struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Sender    userResponse `json:"sender"`
	Receiver  userResponse `json:"receiver"`
}

func (cfg *ApiConfig) HandlerGetAllUsersChats(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.Header.Get("User-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user id malformed")
		return
	}

	chats, err := cfg.DB.GetChatsByUser(r.Context(), userId)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusNotFound, "no chats found")
		return
	}

	chatsResponse := []chatResponse{}
	for _, chat := range chats {
		chatsResponse = append(chatsResponse,
			chatResponse{
				ID:        chat.ID,
				CreatedAt: chat.CreatedAt,
				UpdatedAt: chat.UpdatedAt,
				Sender: userResponse{
					ID:       chat.SenderID.String(),
					Username: chat.SenderUsername,
				},
				Receiver: userResponse{
					ID:       chat.ReceiverID.String(),
					Username: chat.ReceiverUsername,
				},
			},
		)
	}

	log.Printf("retrieved all users chats for %s...", userId)

	respondWithJSON(w, http.StatusOK, chatsResponse)
}

func (cfg *ApiConfig) HandlerDeleteChat(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.Header.Get("User-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user id malformed")
		return
	}

	chatId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chat id malformed")
		return
	}

	chat, err := cfg.DB.GetChatByIdAndSender(r.Context(),
		database.GetChatByIdAndSenderParams{
			ID:       chatId,
			SenderID: userId,
		},
	)
	if err != nil && strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusForbidden, "chat doesn't exist or doesn't belong to user")
		return
	}
	if err != nil && !strings.Contains(err.Error(), noRows) {
		respondWithError(w, http.StatusForbidden, "couldn't find chat")
		return
	}

	err = cfg.DB.DeleteChat(r.Context(), database.DeleteChatParams{ID: chat.ID, SenderID: chat.SenderID})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chat")
		return
	}

	log.Printf("deleted chat %s...", chat.ID.String())

	w.WriteHeader(http.StatusOK)
}
