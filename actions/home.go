package actions

import (
	"errors"
	"fmt"
	"time"

	"net/http"

	"to_do_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"


	"github.com/gobuffalo/pop/v6"
)


func Home(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}
func ShowNewTask(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	countIncomplete := []models.Tasks{}

	qi := tx.Where("finish_at is null")
	if err := qi.All(&countIncomplete); err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", len(countIncomplete))

	c.Set("count", countInfoIncomplete)
	return c.Render(http.StatusOK, r.HTML("new_task.plush.html"))
}

func ShowTableIncomplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := []models.Tasks{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}
	


	if search.NameSearch != "" {
		if err := tx.Where("finish_at IS null AND name_task like  '%" + search.NameSearch + "%'").All(&tasks); err != nil {
		return err
		}
	} else if search.DateSearch != "" {
		if err := tx.Where("finish_at IS null AND extract(day from created_at)= EXTRACT(day FROM TIMESTAMP '"+search.DateSearch+"')").All(&tasks); err != nil {
			return err
		}
		
	} else {if err := tx.Where("finish_at IS null").All(&tasks); err != nil {
		return err
	  }
	}


	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", len(tasks))

	c.Set("count", countInfoIncomplete)
	c.Set("underlineIncomplete", tasks)
	c.Set("tasks", tasks)

	for _, v := range tasks {

		if v.ID.String() == c.Param("task_id") {

			c.Set("taskInfo", v)
			c.Set("finish_at", "this task is not completed")

		}

	}

	return c.Render(http.StatusOK, r.HTML("incomplete_table_tasks.plush.html"))
}

func ShowTableComplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := []models.Tasks{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}


	if search.NameSearch != "" {

		if err := tx.Where("finish_at IS not null AND name_task like  '%" + search.NameSearch + "%'").All(&tasks); err != nil {
		return err
		}

	} else if search.DateSearch != "" {
		if err := tx.Where("finish_at IS not null AND extract(day from created_at) = EXTRACT(day FROM TIMESTAMP '"+search.DateSearch+"')").All(&tasks); err != nil {
			return err
		}
	} else {
		if err := tx.Where("finish_at IS not null").All(&tasks); err != nil {
		return err
	  }
	}

	countInfoComplete := fmt.Sprintf("%v Complete Tasks", len(tasks))

	c.Set("count", countInfoComplete)
	c.Set("underline", tasks)
	c.Set("tasks", tasks)

	for _, v := range tasks {

		if v.ID.String() == c.Param("task_id") {

			c.Set("taskInfo", v)
			c.Set("finish_at", v.Finish_at)

		}

	}

	return c.Render(http.StatusOK, r.HTML("complete_table_tasks.plush.html"))
}

func SendNewTask(c buffalo.Context) error {


	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	u := models.Tasks{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	if  u.Name_task ==""{
		c.Flash().Add("danger", "Alert enter task name!")
	}else {
		err := tx.Create(&u)
		if err != nil {
			return err
		}

	}

	return c.Redirect(302, "/table-incomplete")
}

func Delete(c buffalo.Context) error {

	tasks := &models.Tasks{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(tasks); err != nil {
		return err
	}

	return c.Redirect(302, "/table-incomplete")
}
func ShowEditTask(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := &models.Tasks{}
	countIncomplete := []models.Tasks{}

	qi := tx.Where("finish_at is null")
	if err := qi.All(&countIncomplete); err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", len(countIncomplete))

	c.Set("count", countInfoIncomplete)

	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("tasks", tasks)

	return c.Render(http.StatusOK, r.HTML("edit_task.plush.html"))
}

func Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := &models.Tasks{}
	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(tasks); err != nil {
		return err
	}

	if tasks.Name_task =="" {
		c.Flash().Add("danger", "Alert enter task name!")
	}else{
		if err := tx.Update(tasks); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
	}

	

	return c.Redirect(302, "/table-incomplete")
}

func Check(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction")
	}

	task := &models.Tasks{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	task.Finish_at = nulls.NewTime(time.Now())

	if err := tx.Update(task); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Redirect(302, "/table-incomplete")
}

func UnCheck(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction")
	}

	task := &models.Tasks{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	task.Finish_at = nulls.Time{}

	if err := tx.Update(task); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return c.Redirect(302, "/table-incomplete")
}
