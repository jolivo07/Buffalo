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
	if err := tx.Select(" avg(date_part('year',age(birthdate))) ").All(&ages); err != nil {
		return err
	}


	p := c.Param("infomation_person")
	s := fmt.Sprintf(`%%%v%%`, p)
	
	q := tx.Select(" * , date_part('year',age(birthdate)) as age")
	if c.Param("infomation_person") != "" {
		q.Where("first_name like ? OR last_name like ? OR email like ? ", s, s, s)
	}
	if err := q.All(&users); err != nil{
		return err
	}


	c.Set("avg", ages.AvgAge(ages))
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("users.plush.html"))
}
