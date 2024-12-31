package noteservice

import (
	"net/http"

	"github.com/momoh-yusuf/note-app/utils"
)

func HandleNoteCreation(w http.ResponseWriter, r *http.Request) {

	utils.CustomResponseInJson(w, http.StatusCreated, "Hello", nil)
}
