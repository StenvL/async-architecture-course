package declare

func (c Client) DeclareTaskQueues() error {
	if err := c.declare("tasks", "tasks.created", "created"); err != nil {
		return err
	}
	if err := c.declare("tasks", "tasks.shuffled", "shuffled"); err != nil {
		return err
	}
	if err := c.declare("tasks", "tasks.completed", "completed"); err != nil {
		return err
	}

	return nil
}

func (c Client) declare(exchange, queue, key string) error {
	err := c.mq.Channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = c.mq.Channel.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return err
	}

	return nil
}
