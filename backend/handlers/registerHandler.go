package handlers

import (
	"beehive_api/handlers"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"strings"
	"github.com/jackc/pgx/v5"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
}