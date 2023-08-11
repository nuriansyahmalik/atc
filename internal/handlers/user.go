package handlers

import (
	"encoding/json"
	"github.com/evermos/boilerplate-go/internal/domain/user"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"net/http"
)

type UserHandler struct {
	UserService    user.UserService
	AuthMiddleware *middleware.Authentication
}

func ProvideUserHandler(userService user.UserService, authMiddleware *middleware.Authentication) UserHandler {
	return UserHandler{UserService: userService, AuthMiddleware: authMiddleware}
}

func (h *UserHandler) Router(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Post("/login", h.Login)
	})
}

// CreateUser create a new user
// @Summary Create a new user
// @Description this endpoint create a new user
// @Tags user/user
// @Security JWTAuthentication
// @Param user body user.UserRequestFormat true "The User to be created."
// @Produce json
// @Success 201 {object} response.Base{data=user.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users/ [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat user.UserRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.NewV4()
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	user, err := h.UserService.Create(requestFormat, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, user)
}

// Login logins a new user
// @Summary Login a new user
// @Description this endpoint create a new user
// @Tags user/user
// @Security JWTAuthentication
// @Param user body user.LoginRequestFormat true "The user to be login."
// @Produce json
// @Success 200 {object} response.Base{data=user.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat user.LoginRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, failure.BadRequest(err))
		return
	}

	foo, err := h.UserService.Login(requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, foo)
}
