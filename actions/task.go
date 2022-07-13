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

func ShowNewTask(c buffalo.Context) error {
	tasks := []models.Task{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	q := tx.Where("finished_at is null")
	count, err := q.Count(tasks)
	if err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", count)

	c.Set("count", countInfoIncomplete)

	return c.Render(http.StatusOK, r.HTML("new_task.plush.html"))
}

func ShowTableIncomplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := []models.Task{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}

	q := tx.Where("finished_at IS null")

	if search.Name != "" {
		q.Where("name like  '%" + search.Name + "%'")

	} else if !search.Date.IsZero() {
		q.Where("created_at::date = ?::date", search.Date)
	}

	if err := q.All(&tasks); err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", len(tasks))

	c.Set("count", countInfoIncomplete)
	c.Set("underlineIncomplete", tasks)
	c.Set("tasks", tasks)

	for _, v := range tasks {
		if v.ID.String() == c.Param("task_id") {
			c.Set("taskInfo", v)
			c.Set("finishedAt", "this task is not completed")
		}
	}

	return c.Render(http.StatusOK, r.HTML("incomplete_table_tasks.plush.html"))
}

func ShowTableComplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := []models.Task{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}

	q := tx.Where("finished_at IS not null")

	if search.Name != "" {

		q.Where("name like  '%" + search.Name + "%'")

	} else if !search.Date.IsZero() {

		q.Where("created_at::date = ?::date", search.Date)

	}
	if err := q.All(&tasks); err != nil {
		return err
	}

	countInfoComplete := fmt.Sprintf("%v Complete Tasks", len(tasks))

	c.Set("count", countInfoComplete)
	c.Set("underline", tasks)
	c.Set("tasks", tasks)

	for _, v := range tasks {

		if v.ID.String() == c.Param("task_id") {

			c.Set("taskInfo", v)
			c.Set("finishedAt", v.FinishedAt)

		}

	}

	return c.Render(http.StatusOK, r.HTML("complete_table_tasks.plush.html"))
}

func Create(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	task := models.Task{}
	if err := c.Bind(&task); err != nil {
		return err
	}

	q := tx.Where("finished_at is null")
	count, err := q.Count(task)
	if err != nil {
		return err

	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", count)

	c.Set("count", countInfoIncomplete)

	verrs, err := tx.ValidateAndCreate(&task)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("error", verrs.Get("name"))
		return c.Render(422, r.HTML("new_task.plush.html"))
	}

	return c.Redirect(302, "/table-incomplete")
}

func Delete(c buffalo.Context) error {

	tasks := &models.Task{}
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
func ShowEdit(c buffalo.Context) error {
    
	tasks := &models.Task{}
	
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	q := tx.Where("finished_at is null")
	count, err := q.Count(tasks)
	if err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", count)

	c.Set("count", countInfoIncomplete)

	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("tasks", tasks)

	return c.Render(http.StatusOK, r.HTML("edit_task.plush.html"))
}

func Update(c buffalo.Context) error {
	
	tasks := models.Task{}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	if err := tx.Find(&tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(&tasks); err != nil {
		return err
	}
	q := tx.Where("finished_at is null")
	count, err := q.Count(tasks)
	if err != nil {
		return err

	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", count)

	c.Set("count", countInfoIncomplete)
	c.Set("tasks", &tasks)

	verrs, err := tx.ValidateAndUpdate(&tasks)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("error", verrs.Get("name"))
		return c.Render(422, r.HTML("edit_task.plush.html"))
	}

	// if tasks.Name_task == "" {
	// 	c.Flash().Add("danger", "Alert enter task name!")
	// } else {
	// 	if err := tx.Update(tasks); err != nil {
	// 		return c.Error(http.StatusNotFound, err)
	// 	}
	// }

	return c.Redirect(302, "/table-incomplete")
}

func Check(c buffalo.Context) error {
	
	tasks := &models.Task{}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction")
	}

	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if !tasks.FinishedAt.Valid {
		tasks.FinishedAt = nulls.NewTime(time.Now())
		if err := tx.Update(tasks); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
		return c.Redirect(302, "/table-incomplete")
	}

	c.Flash().Add("danger text-center", "This task is already complete!")
	return c.Redirect(302, "/table-incomplete")
}

func UnCheck(c buffalo.Context) error {

	tasks := &models.Task{}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction")
	}

	

	if err := tx.Find(tasks, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if tasks.FinishedAt.Valid {
		tasks.FinishedAt = nulls.Time{}

		if err := tx.Update(tasks); err != nil {
			return c.Error(http.StatusNotFound, err)
		}
		return c.Redirect(302, "/table-incomplete")
	}

	c.Flash().Add("danger", "This task is already incomplete!")
	return c.Redirect(302, "/table-incomplete")
}
