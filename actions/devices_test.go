package actions

import (
	"fmt"
	"net/http"

	"github.com/DanierJ/div_manager/models"
)

func (as *ActionSuite) Test_Create_Device() {
	// setup a Device model
	d := &models.Device{Manufacturer: "T manufacturer", Make: "T make", Model: "T model", Storage: "T storage", Cost: 120, OS: "Android", ImageURL: "T image url", IsNew: true} // make a POST /widgets request
	res := as.HTML("/devices").Post(d)
	// assert that the response status code was 302 as.Equal(201, res.Code)
	as.Equal(http.StatusSeeOther, res.Code)
	// retrieve the first Widget from the database
	err := as.DB.First(d)
	as.NoError(err)
	as.NotZero(d.ID)
	// assert the Widget title was saved correctly
	as.Equal("T manufacturer", d.Manufacturer)
	// assert the redirect was sent to the place expected
	as.Equal(fmt.Sprintf("/devices/new"), res.Location())

}
