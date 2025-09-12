package http

import (
	"strconv"

	"dunhayat-api/internal/products"
	"dunhayat-api/internal/products/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	listProductsUseCase usecase.ListProductsUseCase
	getProductUseCase   usecase.GetProductUseCase
}

func NewProductHandler(
	listProductsUseCase usecase.ListProductsUseCase,
	getProductUseCase usecase.GetProductUseCase,
) *ProductHandler {
	return &ProductHandler{
		listProductsUseCase: listProductsUseCase,
		getProductUseCase:   getProductUseCase,
	}
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	categoryStr := c.Query("category")
	var category *products.Category

	if categoryStr != "" {
		if cat, err := strconv.Atoi(categoryStr); err == nil {
			catEnum := products.Category(cat)
			category = &catEnum
		}
	}

	products, err := h.listProductsUseCase.Execute(c.Context(), category)
	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := fiber.Map{
		"data":  products,
		"count": len(products),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	product, err := h.getProductUseCase.Execute(c.Context(), productID)
	if err != nil {
		switch err.Error() {
		case "product not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		default:
			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(fiber.Map{
				"error": "Failed to get product",
			})
		}
	}

	response := fiber.Map{
		"data": product,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
