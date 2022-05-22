package declare

func (c Client) DeclareTasksQueues() error {
	_, err := c.mq.Channel.QueueDeclare("analytics.tasks.estimated", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.tasks.estimated", "", "tasks.estimated", false, nil)
	if err != nil {
		return err
	}

	return nil
}
