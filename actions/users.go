package actions

import (
	"errors"
	"fmt"

	"net/http"
	"to_do_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

func ShowUsers(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	users := []models.Users{}
	ages := models.SliceUsers{}
	searchs := &models.Users{}
	if err := tx.Select(" avg(date_part('year',age(birthdate))) ").All(&ages); err != nil {
		return err
	}

	if err := c.Bind(searchs); err != nil {
		return err
	}

	s := fmt.Sprintf(`%%%v%%`, searchs.InformationPerson)

	
	q := tx.Where("first_name like ? OR last_name like ? OR email like ? ", s, s, s)
	if err := q.Select(" * , date_part('year',age(birthdate)) as edad").All(&users); err != nil {
		return err
	}

	fmt.Println(searchs.InformationPerson)

	c.Set("avg", ages.AvgAge(ages))
	c.Set("count", len(users))
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("users.plush.html"))
}
