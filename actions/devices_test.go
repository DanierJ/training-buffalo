package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_Create_Device_With_Valid_Values_And_Existing_User() {

	users := addUsers(as, 1)

	err := as.DB.Find(&users[0], users[0].ID)

	as.NoError(err)

	d := &models.Device{}

	form := url.Values{
		"Manufacturer": []string{"T manufacturer"},
		"Make":         []string{"T make"},
		"Model":        []string{"T model"},
		"Storage":      []string{"T storage"},
		"Cost":         []string{"120"},
		"OS":           []string{"Android"},
		"ImageURL":     []string{"T image url"},
		"IsNew":        []string{"true"},
		"UserID":       []string{users[0].ID.String()},
	}
	res := saveDevice(as, form, 0)

	// assert that the response status code was 302 as.Equal(201, res.Code)
	as.Equal(http.StatusSeeOther, res.Code)

	// retrieve the first Widget from the database
	err = as.DB.First(d)
	as.NoError(err)
	as.NotZero(d.ID)
	// assert the Widget title was saved correctly
	as.Equal("T manufacturer", d.Manufacturer)
	as.Equal(users[0].ID, d.UserID)
	// assert the redirect was sent to the place expected
	as.Equal(fmt.Sprintf("/devices/%v/details", d.ID), res.Location())
}

func (as *ActionSuite) Test_Create_Device_With_Valid_Values_And_Non_Existing_User() {

	fakeID := uuid.Must(uuid.NewV4())

	d := &models.Device{}

	form := url.Values{
		"Manufacturer": []string{"T manufacturer"},
		"Make":         []string{"T make"},
		"Model":        []string{"T model"},
		"Storage":      []string{"T storage"},
		"Cost":         []string{"120"},
		"OS":           []string{"Android"},
		"ImageURL":     []string{"T image url"},
		"IsNew":        []string{"true"},
		"UserID":       []string{fakeID.String()},
	}
	res := saveDevice(as, form, 0)

	// assert that the response status code was 302 as.Equal(201, res.Code)
	as.Equal(http.StatusNotFound, res.Code)

	// retrieve the first Widget from the database
	err := as.DB.First(d)
	as.Error(err)

	as.Empty(d)
}

func (as *ActionSuite) Test_Create_Device_With_Valid_Values_Without_User() {

	d := &models.Device{}

	form := url.Values{
		"Manufacturer": []string{"T manufacturer"},
		"Make":         []string{"T make"},
		"Model":        []string{"T model"},
		"Storage":      []string{"T storage"},
		"Cost":         []string{"120"},
		"OS":           []string{"Android"},
		"ImageURL":     []string{"T image url"},
		"IsNew":        []string{"true"},
		"UserID":       []string{""},
	}
	res := saveDevice(as, form, 0)

	// assert that the response status code was 302 as.Equal(201, res.Code)
	as.Equal(http.StatusSeeOther, res.Code)

	// retrieve the first Widget from the database
	err := as.DB.First(d)
	as.NoError(err)
	as.NotZero(d.ID)
	// assert the Widget title was saved correctly
	as.Equal("T manufacturer", d.Manufacturer)
	// assert the redirect was sent to the place expected
	as.Equal(fmt.Sprintf("/devices/%v/details", d.ID), res.Location())
}

func (as *ActionSuite) Test_Create_Device_With_Invalid_Values() {

	// Empty Values
	d := &models.Device{}

	form := url.Values{
		"manufacturer": []string{""},
		"make":         []string{""},
		"model":        []string{""},
		"storage":      []string{""},
		"cost":         []string{"0"},
		"os":           []string{""},
		"image_url":    []string{""},
		"is_new":       []string{"false"},
	}

	res := saveDevice(as, form, 0)
	invalidDeviceAssertions(as, res, d)
}

func (as *ActionSuite) Test_Update_Device_With_Valid_Values_And_Existing_User() {
	users := addUsers(as, 2)

	devices := addDevices(as, 1, users[0].ID)

	form := url.Values{
		"Manufacturer": []string{"Edited manufacturer"},
		"Make":         []string{"Edited make"},
		"Model":        []string{"Edited model"},
		"Storage":      []string{"Edited storage"},
		"Cost":         []string{"890"},
		"OS":           []string{"iOS"},
		"ImageURL":     []string{"Edited image url"},
		"IsNew":        []string{"false"},
		"UserID":       []string{users[1].ID.String()},
	}

	res := saveDevice(as, form, devices[0].ID)

	updatedDevice := &models.Device{}
	err := as.DB.First(updatedDevice)
	as.NoError(err)

	as.Equal(http.StatusSeeOther, res.Code)

	as.NotEqual(devices[0].Manufacturer, updatedDevice.Manufacturer)
	as.NotEqual(devices[0].UserID, updatedDevice.UserID)
}

func (as *ActionSuite) Test_Update_Device_With_Invalid_Values() {

	devices := addDevices(as, 1, nil)

	form := url.Values{
		"Manufacturer": []string{""},
		"Make":         []string{""},
		"Model":        []string{""},
		"Storage":      []string{""},
		"Cost":         []string{"0"},
		"OS":           []string{"Android"},
		"ImageURL":     []string{""},
		"IsNew":        []string{"false"},
	}

	res := saveDevice(as, form, devices[0].ID)

	updatedDevice := &models.Device{}
	err := as.DB.First(updatedDevice)
	as.NoError(err)

	as.Equal(http.StatusUnprocessableEntity, res.Code)

	as.Equal(devices[0].Manufacturer, updatedDevice.Manufacturer)
}

func (as *ActionSuite) Test_New_Device_View() {
	users := addUsers(as, 2)

	res := as.HTML("/devices/new").Get()

	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), users[0].ID.String())
}

func (as *ActionSuite) Test_Find_Existing_Device() {
	devices := addDevices(as, 1, nil)

	res := as.HTML("/devices/%v/details", devices[0].ID).Get()

	as.Equal(http.StatusOK, res.Code)

}

func (as *ActionSuite) Test_Find_Existing_Device_Without_User_Response() {

	devices := addDevices(as, 1, nil)

	res := as.HTML("/devices/%v/details", devices[0].ID).Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), devices[0].Manufacturer)
}
func (as *ActionSuite) Test_Find_Existing_Device_With_User_Response() {

	users := addUsers(as, 1)

	devices := addDevices(as, 1, users[0].ID)

	res := as.HTML("/devices/%v/details", devices[0].ID).Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), devices[0].Manufacturer)
	as.Contains(res.Body.String(), devices[0].ID.String())
}

func (as *ActionSuite) Test_Find_Non_Existing_Device() {
	fakeID := uuid.Must(uuid.NewV4())
	res := as.HTML("/devices/%v/details", fakeID).Get()

	as.Equal(http.StatusNotFound, res.Code)
}

func (as *ActionSuite) Test_Find_All_Devices() {

	devicesCreated := addDevices(as, 4, nil)

	devices := &[]models.Device{}

	res := as.HTML("/devices").Get()
	as.Equal(http.StatusOK, res.Code)

	err := as.DB.All(devices)

	as.NoError(err)

	as.Equal(len(devicesCreated), len(*devices))

}

func (as *ActionSuite) Test_Delete_Existing_Device() {
	devices := addDevices(as, 1, nil)
	res := as.HTML("/devices/%v", devices[0].ID).Delete()
	as.Equal(http.StatusSeeOther, res.Code)

	device := &models.Device{}

	err := as.DB.Find(device, devices[0].ID)

	as.Error(err)

	as.NotEqual(devices[0].ID, device.ID)
}

func (as *ActionSuite) Test_Delete_Non_Existing_Device() {
	fakeID := uuid.Must(uuid.NewV4())
	res := as.HTML("/devices/%v", fakeID).Delete()
	as.Equal(http.StatusNotFound, res.Code)
}

func (as *ActionSuite) Test_Edit_Existing_Device() {
	devices := addDevices(as, 1, nil)

	res := as.HTML("/devices/%v/edit", devices[0].ID).Get()

	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), devices[0].Manufacturer)
}

func (as *ActionSuite) Test_Edit_Non_Existing_Device() {
	fakeID := uuid.Must(uuid.NewV4())
	res := as.HTML("/devices/%v/edit", fakeID).Get()

	as.Equal(http.StatusNotFound, res.Code)

}

func (as *ActionSuite) Test_Search_Devices_By_Field() {

	devices := addDevices(as, 3, nil)

	name := devices[2].Manufacturer

	res := as.HTML("/devices/?search=" + name).Get()

	as.Equal(http.StatusOK, res.Code)

	as.Contains(res.Body.String(), name)

	as.Equal(0, strings.Count(res.Body.String(), devices[0].Manufacturer))
	as.Equal(0, strings.Count(res.Body.String(), devices[1].Manufacturer))

}

func saveDevice(as *ActionSuite, f url.Values, id interface{}) *httptest.Response {

	switch v := id.(type) {
	case uuid.UUID:
		return as.HTML("/devices/%v", v).Put(f)

	default:
		return as.HTML("/devices").Post(f)
	}
}

func invalidDeviceAssertions(as *ActionSuite, res *httptest.Response, d *models.Device) {
	// assert correct response code
	as.Equal(http.StatusUnprocessableEntity, res.Code)
}

func addDevices(as *ActionSuite, count int, id interface{}) []models.Device {

	var uid nulls.UUID

	switch v := id.(type) {
	case nulls.UUID:
		uid = v
	}

	var devices []models.Device

	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		devices = append(devices, models.Device{Manufacturer: "Manufacturer" + strconv.Itoa(i+1), Make: "Make #" + strconv.Itoa(i+1), Model: "Model #" + strconv.Itoa(i+1), Storage: "Storage #" + strconv.Itoa(i+1), Cost: int64(i + 1*10), OS: "Android", ImageURL: "www", IsNew: true, UserID: uid})

		as.DB.Create(&devices[i])

	}

	return devices
}
