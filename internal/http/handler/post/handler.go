package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/middleware"
	"github.com/rxmy43/support-platform/internal/http/response"
	"github.com/rxmy43/support-platform/internal/modules/post"
)

type PostHandler struct {
	postService *post.PostService
}

func NewPostHandler(postService *post.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid request format", apperror.CodeInvalidRequestFormFormat))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.ToJSON(w, r, apperror.BadRequest("failed to read file", apperror.CodeFileNotFound))
		return
	}
	defer file.Close()

	creatorID, _ := strconv.ParseUint(r.FormValue("creator_id"), 10, 64)

	req := post.PostCreateRequest{
		CreatorID: uint(creatorID),
		Text:      r.FormValue("text"),
		File:      file,
		Header:    header,
	}

	if err := h.postService.Create(r.Context(), req); err != nil {
		response.ToJSON(w, r, err)
		return
	}

	response.ToJSON(w, r, "Post has been created!")
}

func (h *PostHandler) GenerateCaption(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tone string `json:"tone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ToJSON(w, r, apperror.BadRequest("invalid request json format", apperror.CodeInvalidRequestJSONFormat))
		return
	}
	defer r.Body.Close()

	caption, err := h.postService.GenerateCaption(r.Context(), req.Tone)
	if err != nil {
		response.ToJSON(w, r, err)
		return
	}

	resp := map[string]string{"caption": caption}
	response.ToJSON(w, r, resp)
}

func (h *PostHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	var cursor *uint

	if cursorStr != "" {
		parsed, err := strconv.ParseUint(cursorStr, 10, 64)
		if err != nil {
			response.ToJSON(w, r, apperror.BadRequest("invalid cursor", apperror.CodeUnknown))
			return
		}

		temp := uint(parsed)
		cursor = &temp
	}

	var userID *uint
	if userRole := middleware.GetUserRole(r.Context()); userRole != "" && userRole == "creator" {
		userID = middleware.GetUserID(r.Context())
	}

	posts, nextCursor, appErr := h.postService.FindAll(r.Context(), cursor, userID)
	if appErr != nil {
		response.ToJSON(w, r, appErr)
	}

	data := make([]any, len(posts))
	for i, p := range posts {
		data[i] = p
	}

	resp := response.SuccessPaginateResponse{
		Status:     response.StatusSuccess,
		Data:       data,
		NextCursor: nextCursor,
	}

	response.ToJSON(w, r, resp)
}
