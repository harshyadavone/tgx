package tgx

func (b *Bot) OnError(handler ErrorHandler) {
	b.errorHandler = handler
}

func (b *Bot) handleError(ctx *Context, err error) {
	if b.errorHandler != nil {
		b.errorHandler(ctx, err)
	} else {
		b.logger.Error("Unhandled error: %v", err)
	}
}
