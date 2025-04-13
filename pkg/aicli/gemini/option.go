package gemini

const (
	Model20Flash         = "gemini-2.0-flash"
	Model20FlashThinking = "gemini-2.0-flash-thinking-exp"
	Model20Pro           = "gemini-2.0-pro-exp"
	Model25Pro           = "gemini-2.5-pro-exp-03-25"

	DefaultModel = Model25Pro

	RoleUser  = "user"
	RoleModel = "model"
)

// ClientOption is a function that sets a Client option.
type ClientOption func(*Client)

func defaultClientOptions() *Client {
	return &Client{}
}

func (c *Client) apply(opts ...ClientOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithModel sets the model name
func WithModel(name string) ClientOption {
	return func(c *Client) {
		c.ModelName = name
	}
}

// WithEnableContext enable assistant context
func WithEnableContext() ClientOption {
	return func(c *Client) {
		c.enableContext = true
	}
}

// WithInitialContextMessages sets assistant initial context messages
func WithInitialContextMessages(messages ...*ContextMessage) ClientOption {
	return func(c *Client) {
		if len(messages) > 0 {
			c.enableContext = true
			for i, message := range messages {
				if message.Role == "" {
					messages[i].Role = RoleUser // default role is user
				}
			}
			c.contextMessages = messages
		}
	}
}
