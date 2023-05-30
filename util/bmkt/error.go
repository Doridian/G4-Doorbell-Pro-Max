package bmkt

func (ctx *BMKTContext) handleError(code int) {
	ctx.logger.Error().Msgf("Got C error %d", code)
}
