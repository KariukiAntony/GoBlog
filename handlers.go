package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateNoteSchema
	// decode the request body into a payload schema
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	// validate the decoded struct
	errors := ValidateStruct(&payload)
	if errors != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errors)
	}
	now := time.Now()
	newNote := Note{
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Published: payload.Published,
		CreatedAt: now,
		UpdatedAt: now}
	fmt.Println(newNote)

	// create a new note
	result := DB.Create(&newNote)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "failed",
				"message": "Title already exists, please use another title",
			})
			return
		}

		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	// send a success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "created",
		"data": map[string]interface{}{
			"note": newNote,
		},
	})
}

func FindNoteById(w http.ResponseWriter, r *http.Request) {
	noteId := r.PathValue("noteId")
	var note Note
	result := DB.First(&note, "id=?", noteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{
				"status":  "failed",
				"message": fmt.Sprintf("Note with id %v was not found", noteId),
			}
			json.NewEncoder(w).Encode(response)
		}

		w.WriteHeader(http.StatusBadGateway)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":  "failed",
			"message": result.Error.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	// note was found
	response := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"note": note,
		},
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
