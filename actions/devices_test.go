package actions

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DanierJ/div_manager/models"
	"github.com/gobuffalo/httptest"
	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_Create_Device_With_Valid_Values() {

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
	}
	res := saveDevice(as, form, 0)

	// assert that the response status code was 302 as.Equal(201, res.Code)
	as.Equal(http.StatusCreated, res.Code)

	// retrieve the first Widget from the database
	err := as.DB.First(d)
	as.NoError(err)
	as.NotZero(d.ID)
	// assert the Widget title was saved correctly
	as.Equal("T manufacturer", d.Manufacturer)
	// assert the redirect was sent to the place expected
	as.Equal(fmt.Sprintf("/devices/new"), res.Location())
}

func (as *ActionSuite) Test_Update_Device_With_Valid_Values() {

	addDevices(as, 1)

	d := &models.Device{}
	err := as.DB.First(d)
	as.NoError(err)

	form := url.Values{
		"Manufacturer": []string{"Edited manufacturer"},
		"Make":         []string{"Edited make"},
		"Model":        []string{"Edited model"},
		"Storage":      []string{"Edited storage"},
		"Cost":         []string{"890"},
		"OS":           []string{"iOS"},
		"ImageURL":     []string{"Edited image url"},
		"IsNew":        []string{"false"},
	}

	res := saveDevice(as, form, d.ID)

	updatedDevice := &models.Device{}
	err = as.DB.First(updatedDevice)
	as.NoError(err)

	as.Equal(http.StatusOK, res.Code)

	as.NotEqual(d.Manufacturer, updatedDevice.Manufacturer)
}
func (as *ActionSuite) Test_Update_Device_With_Invalid_Values() {

	addDevices(as, 1)

	d := &models.Device{}
	err := as.DB.First(d)
	as.NoError(err)

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

	res := saveDevice(as, form, d.ID)

	updatedDevice := &models.Device{}
	err = as.DB.First(updatedDevice)
	as.NoError(err)

	as.Equal(http.StatusUnprocessableEntity, res.Code)

	as.Equal(d.Manufacturer, updatedDevice.Manufacturer)
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

func (as *ActionSuite) Test_New_Device_View() {
	res := as.HTML("/devices/new").Get()

	as.Equal(http.StatusOK, res.Code)
}

func (as *ActionSuite) Test_Find_Existing_Device() {
	devices := addDevices(as, 1)

	res := as.HTML("/devices/%v/details", devices[0].ID).Get()

	as.Equal(http.StatusOK, res.Code)

}

func (as *ActionSuite) Test_Find_Existing_Device_Response() {

	devices := addDevices(as, 1)

	res := as.HTML("/devices/%v/details", devices[0].ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), devices[0].Manufacturer)
}

func (as *ActionSuite) Test_Find_Non_Existing_Device() {
	fakeID := uuid.Must(uuid.NewV4())
	res := as.HTML("/devices/%v/details", fakeID).Get()

	as.Equal(http.StatusNotFound, res.Code)
}

func (as *ActionSuite) Test_Find_All_Devices() {

	devicesCreated := addDevices(as, 4)

	devices := &[]models.Device{}

	res := as.HTML("/devices").Get()
	as.Equal(http.StatusOK, res.Code)

	err := as.DB.All(devices)

	as.NoError(err)

	as.Equal(len(devicesCreated), len(*devices))

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

func addDevices(as *ActionSuite, count int) []models.Device {
	var devices []models.Device

	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		devices = append(devices, models.Device{Manufacturer: "Manufacturer #" + strconv.Itoa(i+1), Make: "Make #" + strconv.Itoa(i+1), Model: "Model #" + strconv.Itoa(i+1), Storage: "Storage #" + strconv.Itoa(i+1), Cost: int64(i + 1*10), OS: "Android", ImageURL: "www", IsNew: true})

		as.DB.Create(&devices[i])

	}

	return devices
}
