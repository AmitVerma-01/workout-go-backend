package api

import (
	"encoding/json"
	"errors"
	"go_beginner/internals/store"
	"go_beginner/utils"
	"log"
	"net/http"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error:: Reading user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := uh.userStore.GetUserByID(userId)
	if err != nil {
		uh.logger.Printf("Error:: Getting user by ID: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		uh.logger.Printf("Error:: User not found for ID: %d", userId)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"user": user,
	})
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error:: Decoding create user request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}
	if err := uh.validateRegisterRequest(&req); err != nil {
		uh.logger.Printf("Error:: Validating create user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": err.Error(),
		})
		return
	}

	user := &store.User{
		Name:  req.Name,
		Email: req.Email,
		Bio:   req.Bio,
	}
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("Error:: Setting user password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to set user password",
		})
		return
	}
	createdUser, err := uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("Error:: Creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to create user",
		})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"user": createdUser,
	})
}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error:: Reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid user ID",
		})
		return
	}
	type UpdateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Bio   string `json:"bio"`
	}
	var req UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error:: Decoding update user request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid request body",
		})
		return
	}
	user := &store.User{
		ID:    userId,
		Name:  req.Name,
		Email: req.Email,
		Bio:   req.Bio,
	}
	err = uh.userStore.UpdateUser(userId, user)
	if err != nil {
		uh.logger.Printf("Error:: Updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to update user",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"user": user,
	})
}

func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error:: Reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Invalid user ID",
		})
		return
	}
	err = uh.userStore.DeleteUser(userId)
	if err != nil {
		uh.logger.Printf("Error:: Deleting user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to delete user",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (uh *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.userStore.GetUsers()
	if err != nil {
		uh.logger.Printf("Error:: Getting users: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to get users",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"users": users,
	})
}

func (uh *UserHandler) validateRegisterRequest(req *createUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !utils.MatchRegex(emailRegex, req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	if req.Bio == "" {
		return errors.New("bio is required")
	}
	return nil
}

func (uh *UserHandler) HandleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		uh.logger.Println("Error:: Email query parameter is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "Email query parameter is required",
		})
		return
	}
	user, err := uh.userStore.GetUserByEmail(email)
	if err != nil {
		uh.logger.Printf("Error:: Getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "Failed to get user by email " + err.Error(),
		})
		return
	}
	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{
			"error": "User not found",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"user": user,
	})
}
