package actions

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"to_do_app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

func ShowNewTask(c buffalo.Context) error {
	tasks := models.Tasks{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	c.Set("count", tasks.Count(tx))

	return c.Render(http.StatusOK, r.HTML("new_task.plush.html"))
}

func ShowTableIncomplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := models.Tasks{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}

	q := tx.Where("finished_at IS null")

	if search.Name != "" {
		s := fmt.Sprintf("%%%v%%", search.Name)
		q.Where("name like  ? ", s)

	} else if !search.Date.IsZero() {
		q.Where("created_at::date = ?::date", search.Date)
	}

	if err := q.All(&tasks); err != nil {
		return err
	}

	countInfoIncomplete := fmt.Sprintf("%v Incomplete Tasks", len(tasks))
	taskInfo := tasks.InfoTask(uuid.FromStringOrNil(c.Param("task_id")))

	c.Set("taskInfo", taskInfo)
	c.Set("count", countInfoIncomplete)
	c.Set("tasksIncomplete", tasks)

	return c.Render(http.StatusOK, r.HTML("incomplete_table_tasks.plush.html"))
}

func ShowTableComplete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	tasks := models.Tasks{}
	search := &models.Search{}

	if err := c.Bind(search); err != nil {
		return err
	}

	q := tx.Where("finished_at IS not null")

	if search.Name != "" {
		s := fmt.Sprintf("%%%v%%", search.Name)
		q.Where("name like  ?", s)

	} else if !search.Date.IsZero() {

		q.Where("created_at::date = ?::date", search.Date)

	}
	if err := q.All(&tasks); err != nil {
		return err
	}

	countTask := fmt.Sprintf("%v Complete Tasks", len(tasks))
	taskInfo := tasks.InfoTask(uuid.FromStringOrNil(c.Param("task_id")))

	c.Set("taskInfo", taskInfo)
	c.Set("count", countTask)
	c.Set("tasksComplete", tasks)

	return c.Render(http.StatusOK, r.HTML("complete_table_tasks.plush.html"))
}

func Create(c buffalo.Context) error {

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	task := models.Task{}
	tasks := models.Tasks{}
	if err := c.Bind(&task); err != nil {
		return err
	}

	c.Set("count", tasks.Count(tx))

	verrs, err := tx.ValidateAndCreate(&task)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("error", verrs.Get("name"))
		return c.Render(http.StatusUnprocessableEntity, r.HTML("new_task.plush.html"))
	}

	return c.Redirect(http.StatusSeeOther, "/table-incomplete")
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

	task := &models.Task{}
	tasks := models.Tasks{}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	c.Set("count", tasks.Count(tx))

	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("task", task)

	return c.Render(http.StatusOK, r.HTML("edit_task.plush.html"))
}

func Update(c buffalo.Context) error {

	task := models.Task{}
	tasks := models.Tasks{}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	if err := tx.Find(&task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(&task); err != nil {
		return err
	}
	

	c.Set("count", tasks.Count(tx))
	verrs, err := tx.ValidateAndUpdate(&task)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("error", verrs.Get("name"))
		return c.Render(http.StatusUnprocessableEntity, r.HTML("edit_task.plush.html"))
	}

	return c.Redirect(http.StatusSeeOther, "/table-incomplete")
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
		return c.Redirect(http.StatusSeeOther, "/table-incomplete")
	}

	c.Flash().Add("danger text-center", "This task is already complete!")
	return c.Redirect(http.StatusSeeOther, "/table-incomplete")
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
		return c.Redirect(http.StatusSeeOther, "/table-incomplete")
	}

	c.Flash().Add("danger", "This task is already incomplete!")
	return c.Redirect(http.StatusSeeOther, "/table-incomplete")
}
