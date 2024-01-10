package handler

import (
	"fmt"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type OmiseHandler struct {
	omiseService ports.OmiseService
	userService  ports.UserService
}

func NewOmiseHandler(omiseService ports.OmiseService, userService ports.UserService) *OmiseHandler {
	return &OmiseHandler{
		omiseService: omiseService,
		userService:  userService,
	}
}

type Token struct {
	Name            string     `json:"name"`
	Number          string     `json:"number"`
	ExpirationMonth time.Month `json:"expiration_month"`
	ExpirationYear  int        `json:"expiration_year"`
}

type ChargeCreditCardInput struct {
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
}

func (h *OmiseHandler) ChargeCreditCard(c *fiber.Ctx) error {
	chargeInput := &ChargeCreditCardInput{}
	if err := c.BodyParser(chargeInput); err != nil {
		fmt.Println(err)
	}
	amount := chargeInput.Amount
	token := chargeInput.Token
	userID := c.Params("userID")

	charge, e := h.omiseService.ChargeCreditCard(amount, token, userID)
	if e != nil {
		return c.Status(500).SendString(e.Error())
	}
	fmt.Printf("[%s]charge: %s  amount: %s %d\n", charge.Status, charge.ID, charge.Currency, charge.Amount)
	if charge.Status == "failed" {
		fmt.Println(charge.FailureCode, charge.FailureMessage)
		return c.Status(500).JSON(fiber.Map{"failure_code": charge.FailureCode, "message": charge.FailureMessage})
	}
	return c.JSON(fiber.Map{"Charge ID": charge.ID, "Amount": charge.Amount, "Status": charge.Status, "Charge": charge})
}

type ChargeBankingInput struct {
	Source string `json:"source"`
	Amount int64  `json:"amount"`
}

func (h *OmiseHandler) ChargeBanking(c *fiber.Ctx) error {
	chargeInput := &ChargeBankingInput{}
	if err := c.BodyParser(chargeInput); err != nil {
		fmt.Println(err)
	}
	amount := chargeInput.Amount
	source := chargeInput.Source
	userID := c.Params("userID")

	charge, err := h.omiseService.ChargeBanking(amount, source, userID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fmt.Printf("[%s]charge: %s  amount: %s %d\n", charge.Status, charge.ID, charge.Currency, charge.Amount)
	if charge.Status == "failed" {
		fmt.Println(charge.FailureCode, charge.FailureMessage)
		return c.Status(500).JSON(fiber.Map{"failure_code": charge.FailureCode, "message": charge.FailureMessage})
	}
	return c.JSON(fiber.Map{"Charge ID": charge.ID, "Amount": charge.Amount, "Currency": charge.Currency, "AuthorizeURI": charge.AuthorizeURI, "Status": charge.Status, "Charge": charge})
}

func (h *OmiseHandler) CreateToken(c *fiber.Ctx) error {
	token := &Token{}
	if err := c.BodyParser(token); err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	result, err := h.omiseService.CreateToken(token.Name, token.Number, token.ExpirationMonth, token.ExpirationYear)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(result)
	return c.JSON(result)
}

func (h *OmiseHandler) ListCustomers(c *fiber.Ctx) error {
	result, err := h.omiseService.ListCustomers()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(result)
	return c.JSON(result)
}

type AttachCardRequest struct {
	CustomerID string `json:"customer_id"`
	Card       string `json:"card"`
}

func (h *OmiseHandler) AttchCardToCustomer(c *fiber.Ctx) error {
	req := &AttachCardRequest{}
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
	}

	customer, err := h.omiseService.AttachCardToCustomer(req.CustomerID, req.Card)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(customer)
	return c.JSON(customer)
}

func (h *OmiseHandler) RetrieveCharge(c *fiber.Ctx) error {
	chargeID := c.Params("charge_id")

	charge, err := h.omiseService.RetrieveCharge(chargeID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}

	fmt.Println(charge)
	return c.JSON(charge)
}

func (h *OmiseHandler) HandleWebhook(c *fiber.Ctx) error {

	return c.SendStatus(int(c.Body()[0]))
}

func (h *OmiseHandler) GetCharges(c *fiber.Ctx) error {
	result, err := h.omiseService.GetCharges()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}

func (h *OmiseHandler) GetTransaction(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")
	result, err := h.omiseService.GetTransaction(transactionID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}

func (h *OmiseHandler) GetCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customerToken")
	result, err := h.omiseService.GetCustomer(customerID)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}
