# Backlog
- [x] COULD add request id middleware and take it into account while logging
- [x] SHOULD only add request id as field, not other parts. Those should be logged once
- [x] SHOULD make stdlib logging for fx
- [ ] SHOULD consider not using fx.Annotated with name, but just have a NamedType for the middleware
- [ ] COULD log every request with duration
- [ ] COULD make common zap options logging options configurable through env
- [ ] COULD configure global logger by default