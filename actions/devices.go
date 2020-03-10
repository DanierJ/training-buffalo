package actions

import (
	"fmt"
	"net/http"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
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

	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	if verrs.HasAny() {
		// Setting errors
		c.Set("errors", verrs)

		c.Set("device", device)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("/devices/new.plush.html"))
	}

	c.Flash().Add("success", "Device created succesfully")

	return c.Redirect(http.StatusSeeOther, "/devices/new")

}

// New handles the request to render the create device view
func New(c buffalo.Context) error {

	c.Set("device", &models.Device{})
	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	return c.Render(http.StatusOK, r.HTML("/devices/new.plush.html"))
}
