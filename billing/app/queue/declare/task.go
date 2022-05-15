package declare

func (c Client) DeclareTaskQueues() error {
	_, err := c.mq.Channel.QueueDeclare("billing.tasks.created", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("billing.tasks.created", "", "tasks.created", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("billing.tasks.created.v2", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("billing.tasks.created.v2", "", "tasks.created", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("billing.tasks.shuffled", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("billing.tasks.shuffled", "", "tasks.shuffled", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("billing.tasks.completed", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("billing.tasks.completed", "", "tasks.completed", false, nil)
	if err != nil {
		return err
	}

	return nil
}
