package bmkt

func (ctx *Context) handleError(code int) {
	ctx.logger.Error().Msgf("Got C error %d", code)
}
