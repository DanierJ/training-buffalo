package actions

import (
	"fmt"
	"net/http"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// ListUser default implementation.
func ListUser(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	users := &[]models.User{}

	if err := tx.All(users); err != nil {
		return err
	}

	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("/users/list.plush.html"))

}

// CreateUser default implementation.
func CreateUser(c buffalo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	verrs, err := tx.ValidateAndCreate(user)

	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("/users/new.plush.html"))
	}

	c.Flash().Add("success", "User added successfully")
	return c.Redirect(http.StatusSeeOther, "/users/%v/details", user.ID)

}

// UpdateUser default implementation.
func UpdateUser(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(user); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(user)

	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("/users/edit.plush.html"))
	}

	c.Set("user", user)

	c.Flash().Add("success", "User updated succesfully")

	return c.Redirect(http.StatusSeeOther, "/users/%v/details", user.ID)

}

// DeleteUser default implementation.
func DeleteUser(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return err
	}

	c.Flash().Add("success", "User deleted succesfully")

	return c.Redirect(http.StatusSeeOther, "/users")
}

// ShowUser default implementation.
func ShowUser(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return fmt.Errorf("no transaction found")
	}

	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("user", user)

	return c.Render(http.StatusOK, r.HTML("/users/show.plush.html"))
}
