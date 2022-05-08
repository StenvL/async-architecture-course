class Producer
  def call(event, topic:, key:)
    case ENV['BROKER_ADAPTER']
    when 'rabbitmq'
      send_rabbitmq_event(event, topic, key)
    end
  end

private

  def send_rabbitmq_event(event, topic, key)
    require 'bunny'

    connection = Bunny.new(ENV['AMQP_URL'])
    connection.start

    channel = connection.create_channel

    # It shouldn't be here, but I'm not a Ruby developer and don't know where to place it.
    exchange = channel.direct('users')

    created = channel.queue('users.created', durable: true)
    created.bind(exchange, routing_key: 'created')

    updated = channel.queue('users.updated', durable: true)
    updated.bind(exchange, routing_key: 'updated')

    deleted = channel.queue('users.deleted', durable: true)
    deleted.bind(exchange, routing_key: 'deleted')

    role_changed = channel.queue('users.role_changed', durable: true)
    role_changed.bind(exchange, routing_key: 'role_changed')

    exchange = get_rabbitmq_exchange(channel, topic).publish(event.to_json, routing_key: key)

    connection.close
  end

  def get_rabbitmq_exchange(channel, topic)
    channel.find_exchange(topic) || channel.direct(topic)
  end

end
