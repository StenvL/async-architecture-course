package declare

func (c Client) DeclareUserQueues() error {
	_, err := c.mq.Channel.QueueDeclare("analytics.users.created", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.users.created", "", "users.created", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("analytics.users.updated", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.users.updated", "", "users.updated", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("analytics.users.deleted", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.users.deleted", "", "users.deleted", false, nil)
	if err != nil {
		return err
	}

	_, err = c.mq.Channel.QueueDeclare("analytics.users.role_changed", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.users.role_changed", "", "users.role_changed", false, nil)
	if err != nil {
		return err
	}

	return nil
}
