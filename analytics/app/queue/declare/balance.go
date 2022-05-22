package declare

func (c Client) DeclareBalanceQueues() error {
	_, err := c.mq.Channel.QueueDeclare("analytics.balance.changed", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = c.mq.Channel.QueueBind("analytics.balance.changed", "", "balance.changed", false, nil)
	if err != nil {
		return err
	}

	return nil
}
