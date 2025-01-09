package tgx

func (b *Bot) OnCommand(command string, handler handler) {
	b.commandHandler[command] = handler
}
