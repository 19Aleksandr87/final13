package db

import "strconv"

type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func ConvertTaskToTaskDTO(task *Task) TaskDTO {
	return TaskDTO{
		ID:      strconv.Itoa(task.ID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}

func ConvertToTaskDTO_Task(dto TaskDTO) (Task, error) {
	i, err := strconv.Atoi(dto.ID)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:      i,
		Date:    dto.Date,
		Title:   dto.Title,
		Comment: dto.Comment,
		Repeat:  dto.Repeat,
	}, nil
}

func convertTaskToTaskDTO(task Task) TaskDTO {
	return TaskDTO{
		ID:      strconv.Itoa(task.ID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
