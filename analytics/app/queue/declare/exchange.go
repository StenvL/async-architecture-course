package declare

func (c Client) DeclareExchanges() error {
	if err := c.declare("users.created"); err != nil {
		return err
	}
	if err := c.declare("users.updated"); err != nil {
		return err
	}
	if err := c.declare("users.deleted"); err != nil {
		return err
	}
	if err := c.declare("users.role_changed"); err != nil {
		return err
	}

	if err := c.declare("tasks.estimated"); err != nil {
		return err
	}

	if err := c.declare("balance.changed"); err != nil {
		return err
	}

	return nil
}

func (c Client) declare(exchange string) error {
	err := c.mq.Channel.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}
