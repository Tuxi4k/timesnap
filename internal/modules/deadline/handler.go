package deadline

import (
	"errors"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type request struct {
	Title    *string    `json:"title"`
	Status   *Status    `json:"status"`
	Priority *Priority  `json:"priority"`
	DueDate  *time.Time `json:"due_date"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Post("/", h.Create)
	r.Patch("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

// @Summary 		Get all deadlines
// @Description		Returns list of all deadlines
// @Tags			deadlines
// @Produce			json
// @Success			200	{array}		Deadline
// @Failure			500	{object}	map[string]string
// @Router			/deadlines		[get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	tasks, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(tasks)
}

// @Summary			Get deadline by id
// @Description		Returns deadline by your id
// @Tags			deadlines
// @Produce			json
// @Param			id	path	integer	true "Deadline id"
// @Success			200	{object}	Deadline
// @Failure			400	{object}	map[string]string
// @Failure			404	{object}	map[string]string
// @Failure			500	{object}	map[string]string
// @Router			/deadlines/{id}	[get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	task, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(task)
}

// @Summary			Create deadline
// @Tags			deadlines
// @Accept			json
// @Produce			json
// @Success			201	{object}	map[string]string
// @Failure			400	{object}	map[string]string
// @Failure			422	{object}	map[string]string
// @Failure			500	{object}	map[string]string
// @Router			/deadlines		[post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req request

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	input := Input{
		Title:    req.Title,
		Status:   req.Status,
		Priority: req.Priority,
		DueDate:  req.DueDate,
	}

	task, err := h.service.Create(input)
	if err != nil {
		errs, ok := err.(validation.Errors)
		if ok {
			return c.Status(422).JSON(errs)
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(task)
}

// @Summary			Updates deadline
// @Description		Partial updates deadline by your id
// @Tags			deadlines
// @Accept			json
// @Produce			json
// @Param			id	path	integer	true "Deadline id"
// @Param			deadline	body	request	false "Updatable fields"
// @Success			200	{object}	Deadline
// @Failure			400	{object}	map[string]string
// @Failure			422	{object}	map[string]string
// @Failure			404	{object}	map[string]string
// @Failure			500	{object}	map[string]string
// @Router			/deadlines/{id}	[patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req request

	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	input := Input{
		Title:    req.Title,
		Status:   req.Status,
		Priority: req.Priority,
		DueDate:  req.DueDate,
	}

	task, err := h.service.Update(uint(id), input)
	if err != nil {
		errs, ok := err.(validation.Errors)
		if ok {
			return c.Status(422).JSON(fiber.Map{"error": errs})
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(task)
}

// @Summary			Delete deadline
// @Description		Delete deadline by your id
// @Tags			deadlines
// @Accept			json
// @Produce			json
// @Param			id	path	integer	true "Deadline id"
// @Success			204	"No content"
// @Failure			400	{object}	map[string]string
// @Failure			404	{object}	map[string]string
// @Failure			500	{object}	map[string]string
// @Router			/deadlines/{id}	[delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
