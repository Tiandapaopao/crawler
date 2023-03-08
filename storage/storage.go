package storage

type DataCell struct {
	Data map[string]interface{}
}

func (c *DataCell) GetTableName() string {
	return c.Data["Task"].(string)
}

func (c *DataCell) GetTaskName() string {
	return c.Data["Task"].(string)
}

type Storage interface {
	Save(data ...*DataCell) error
}
