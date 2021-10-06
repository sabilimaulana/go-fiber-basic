package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Todo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: 1, Name: "Walk the dog", Completed: false},
	{ID: 2, Name: "Walk the cat", Completed: false},
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	SetupApiV1(app)

	log.Fatal(app.Listen(":8080"))
}

func SetupApiV1(app *fiber.App) {
	v1 := app.Group("/v1")
	SetupTodosRoutes(v1)
}

func SetupTodosRoutes(grp fiber.Router) {
	todosRoutes := grp.Group("/todos")

	todosRoutes.Get("/", GetTodos)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

func GetTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}

func CreateTodo(c *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse json body",
		})
	}

	todo := Todo{
		ID:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	return c.Status(fiber.StatusOK).JSON(todo)
}

func GetTodo(c *fiber.Ctx) error {
	paramsID := c.Params("id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Id must be number",
		})
	}

	for _, todo := range todos {
		if id == todo.ID {
			return c.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func DeleteTodo(c *fiber.Ctx) error {
	paramsID := c.Params("id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Id must be number",
		})
	}

	for i, todo := range todos {
		if id == todo.ID {
			todos = append(todos[0:i], todos[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func UpdateTodo(c *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse json body",
		})
	}

	paramsID := c.Params("id")
	id, err := strconv.Atoi(paramsID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Id must be number",
		})
	}

	var todo Todo

	for _, t := range todos {
		if id == t.ID {
			todo = t
		}
	}

	if todo.ID == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	for i, t := range todos {
		if id == t.ID {
			todos[i] = todo
		}
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}
