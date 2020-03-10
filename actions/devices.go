package actions

import (
	"fmt"
	"net/http"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/x/responder"
)

// Create handles the request to post a new device
func Create(c buffalo.Context) error {
	device := &models.Device{}

	// Binding the values
	if err := c.Bind(device); err != nil {
		return err
	}

	// Getting db connection from context
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validating send data
	verrs, err := tx.ValidateAndCreate(device)

	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Setting errors
			c.Set("errors", verrs)

			c.Set("device", device)
			c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/devices/new.plush.html"))

		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Flash().Add("success", "Device created succesfully")

		c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

		return c.Redirect(http.StatusSeeOther, "/devices/new")

	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(device))
	}).Respond(c)

}

// New handles the request to render the create device view
func New(c buffalo.Context) error {

	c.Set("device", &models.Device{})
	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	return c.Render(http.StatusOK, r.HTML("/devices/new.plush.html"))
}
