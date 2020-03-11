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

	return c.Redirect(http.StatusCreated, "/devices/new")

}

// New handles the request to render the create device view
func New(c buffalo.Context) error {

	c.Set("device", &models.Device{})
	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	return c.Render(http.StatusOK, r.HTML("/devices/new.plush.html"))
}

// Show handles the request to get one device
func Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	device := &models.Device{}

	if err := tx.Find(device, c.Param("device_id")); err != nil {

		return c.Error(http.StatusNotFound, err)
	}

	c.Set("device", device)

	return c.Render(http.StatusOK, r.HTML("/devices/show.plush.html"))

}

// Update hanlde the rquest to update a device
func Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	device := &models.Device{}

	if err := tx.Find(device, c.Param("device_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Todo to the html form elements
	if err := c.Bind(device); err != nil {
		return err
	}

	// Validating send data
	verrs, err := tx.ValidateAndUpdate(device)

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

	c.Flash().Add("success", "Device updated succesfully")

	return c.Redirect(http.StatusOK, "/devices/%v/details", device.ID)

}

// List handles request to list devices
func List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	devices := &[]models.Device{}

	if err := tx.All(devices); err != nil {
		return err
	}

	c.Set("devices", devices)

	return c.Render(http.StatusOK, r.HTML("/devices/list.plush.html"))

}
