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

func (as *ActionSuite) Test_Create_User_With_Valid_Values() {
	u := &models.User{}

	form := url.Values{
		"Name":     []string{"John"},
		"Lastname": []string{"Doe"},
		"Email":    []string{"john@doe.com"},
	}

	res := saveUser(as, form, 0)
	as.Equal(http.StatusSeeOther, res.Code)

	err := as.DB.First(u)
	as.NoError(err)
	as.NotZero(u.ID)

	// assert the Widget title was saved correctly
	as.Equal("John", u.Name)
	// assert the redirect was sent to the place expected
	as.Equal(fmt.Sprintf("/users/%v/details", u.ID), res.Location())

}

func (as *ActionSuite) Test_Create_User_With_Invalid_Values() {

	user := &models.User{}

	form := url.Values{
		"Name":     []string{""},
		"Lastname": []string{""},
		"Email":    []string{""},
	}

	res := saveUser(as, form, 0)

	as.Equal(http.StatusUnprocessableEntity, res.Code)

	form = url.Values{
		"Name":     []string{"Dan"},
		"Lastname": []string{"Jav"},
		"Email":    []string{"NoValidEmail"},
	}

	res = saveUser(as, form, 0)

	err := as.DB.Find(user, user.ID)
	as.NoError(err)
	as.Empty(user)
}

func (as *ActionSuite) Test_Find_Existing_User() {
	users := addUsers(as, 1)

	res := as.HTML("/users/%v/details", users[0].ID).Get()

	as.Equal(http.StatusOK, res.Code)
}

func (as *ActionSuite) Test_Find_Non_Existing_User() {
	fakeID := uuid.Must(uuid.NewV4())

	res := as.HTML("/users/%v/details", fakeID).Get()

	as.Equal(http.StatusNotFound, res.Code)
}

func (as *ActionSuite) Test_Find_Existing_User_Response() {

	users := addUsers(as, 1)

	res := as.HTML("/users/%v/details", users[0].ID).Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), users[0].Name)
}

func (as *ActionSuite) Test_Update_User_With_Valid_Values() {

	users := addUsers(as, 1)

	form := url.Values{
		"Name":     []string{"Jane"},
		"Lastname": []string{"Doe"},
		"Email":    []string{"jane@doe.com"},
	}

	res := saveUser(as, form, users[0].ID)
	updatedUser := &models.User{}
	err := as.DB.Find(updatedUser, users[0].ID)

	as.NoError(err)
	as.Equal(http.StatusSeeOther, res.Code)
	as.NotEqual(users[0].Name, updatedUser.Name)

}

func (as *ActionSuite) Test_Update_User_With_Invalid_Values() {

	users := addUsers(as, 1)

	form := url.Values{
		"Name":     []string{""},
		"Lastname": []string{""},
		"Email":    []string{""},
	}

	res := saveUser(as, form, users[0].ID)
	updatedUser := &models.User{}
	err := as.DB.Find(updatedUser, users[0].ID)

	as.NoError(err)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Equal(users[0].Name, updatedUser.Name)
}

func (as *ActionSuite) Test_Find_All_Users() {
	users := addUsers(as, 5)

	res := as.HTML("/users").Get()

	as.Equal(http.StatusOK, res.Code)

	foundUsers := &[]models.User{}

	err := as.DB.All(foundUsers)

	as.NoError(err)
	as.Equal(len(users), len(*foundUsers))

}

func (as *ActionSuite) Test_Delete_Existing_User() {
	users := addUsers(as, 2)

	res := as.HTML("/users/%v", users[0].ID).Delete()
	as.Equal(http.StatusSeeOther, res.Code)

	user := &models.User{}

	err := as.DB.Find(user, users[0].ID)

	as.Error(err)

	as.NotEqual(users[0].ID, user.ID)

}

func (as *ActionSuite) Test_Delete_Non_Existing_User() {
	fakeID := uuid.Must(uuid.NewV4())
	res := as.HTML("/users/%v", fakeID).Delete()
	as.Equal(http.StatusNotFound, res.Code)
}

func saveUser(as *ActionSuite, f url.Values, id interface{}) *httptest.Response {

	switch v := id.(type) {
	case uuid.UUID:
		return as.HTML("/users/%v", v).Put(f)

	default:
		return as.HTML("/users").Post(f)
	}
}

func addUsers(as *ActionSuite, count int) []models.User {
	var users []models.User

	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		users = append(users, models.User{Name: "User" + strconv.Itoa(i+1), Lastname: "Last" + strconv.Itoa(i+1), Email: "user" + strconv.Itoa(i+1) + "@user.com"})

		as.DB.Create(&users[i])

	}

	return users
}
