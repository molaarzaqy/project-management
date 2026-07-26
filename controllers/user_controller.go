package controllers

import (
	"math"
	"strconv"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user  := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "gagal parsing data", err.Error())
	}
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "register failed", err.Error())
	}
	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "register success", userResp)
}
// Login godoc
// @Summary User login
// @Description Masuk ke dalam sistem menggunakan email dan password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}
	user, err := c.service.Login(body.Email, body.Password)
	if err != nil{
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token ,_ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	
	return utils.Success(ctx, "Login Successfully", fiber.Map{
		"access_token": token,
		"refresh_token":refreshToken,
		"user": userResp,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "data not found", err.Error())
	}
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil {
		return utils.BadRequest(ctx, "internal server error", err.Error())
	}
	return utils.Success(ctx, "user found", userResp)
}

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page -1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")
	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "failed to retrieve data", err.Error())
	}
	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)
	meta := utils.PaginationMeta{
		Page: page,
		Limit: limit,
		Total: int(total),
		TotalPages: int(math.Ceil(float64(total)/float64(limit))),
		Filter: filter,
		Sort: sort,
	}
	if total == 0 {
		utils.NotFoundPagination(ctx, "data users not found", userResp, meta)
	}
	return utils.SuccessPagination(ctx, "data users found", userResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "invalid ID format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to parse data", err.Error())
	}
	user.PublicID = publicID
	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "failed to update user", err.Error())
	}
	userUpdated, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to get data", err.Error())
	}
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "something went wrong", err.Error())
	}
	return utils.Success(ctx, "success to update data", userResp)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ :=strconv.Atoi(ctx.Params("id"))
	if err := c.service.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Failed to Delete data", err.Error())
	}
	return utils.Success(ctx, "User deleted successfully", id)
}