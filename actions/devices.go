package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// CreateDevice handles the request to post a new device
func CreateDevice(c buffalo.Context) error {
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

	if err := tx.Find(&models.User{}, device.UserID); err != nil && device.UserID.Valid == true {

		return c.Error(http.StatusNotFound, err) // should I return an err or set the id to null
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

	return c.Redirect(http.StatusSeeOther, "/devices/%v/details", device.ID)

}

// NewDevice handles the request to render the create device view
func NewDevice(c buffalo.Context) error {

	c.Set("device", &models.Device{})
	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	users := &[]models.User{}

	if err := tx.All(users); err != nil {
		return err
	}

	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("/devices/new.plush.html"))
}

// ShowDevice handles the request to get one device
func ShowDevice(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	device := &models.Device{}

	if err := tx.Eager().Find(device, c.Param("device_id")); err != nil {

		return c.Error(http.StatusNotFound, err)
	}

	c.Set("device", device)

	return c.Render(http.StatusOK, r.HTML("/devices/show.plush.html"))

}

// UpdateDevice hanlde the rquest to update a device
func UpdateDevice(c buffalo.Context) error {
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

	if err := tx.Find(&models.User{}, device.UserID); err != nil {
		return c.Error(http.StatusNotFound, err)
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

		return c.Render(http.StatusUnprocessableEntity, r.HTML("/devices/edit.plush.html"))
	}

	c.Flash().Add("success", "Device updated succesfully")

	return c.Redirect(http.StatusSeeOther, "/devices/%v/details", device.ID)

}

// ListDevice handles request to list devices
func ListDevice(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	devices := &[]models.Device{}

	field := c.Param("search")

	cost, err := strconv.ParseInt(field, 10, 64)

	if err != nil {
		cost = 0
	}

	if field != "" {
		query := tx.Where("model = ? OR manufacturer = ? OR cost = ?", field, field, cost)
		if err := query.All(devices); err != nil {
			return err
		}
	} else {
		if err := tx.All(devices); err != nil {
			return err
		}
	}

	c.Set("devices", devices)

	return c.Render(http.StatusOK, r.HTML("/devices/list.plush.html"))

}

// DeleteDevice handles request to delete devices
func DeleteDevice(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	device := &models.Device{}

	if err := tx.Find(device, c.Param("device_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(device); err != nil {
		return err
	}
	c.Flash().Add("success", "Device deleted succesfully")

	return c.Redirect(http.StatusSeeOther, "/devices")
}

// EditDevice handle the request to render device to be edited
func EditDevice(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("No transaction found")
	}

	device := &models.Device{}

	if err := tx.Eager().Find(device, c.Param("device_id")); err != nil {

		return c.Error(http.StatusNotFound, err)
	}

	users := &[]models.User{}

	if err := tx.All(users); err != nil {
		return err
	}

	c.Set("users", *users)
	c.Set("device", device)
	c.Set("osOptions", models.OS{"Android", "iOS", "Windows"})

	return c.Render(http.StatusOK, r.HTML("/devices/edit.plush.html"))

}
