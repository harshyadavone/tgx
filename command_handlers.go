package tgx

func (b *Bot) OnCommand(command string, handler Handler) {
	b.commandHandler[command] = handler
}
