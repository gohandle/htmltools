# Backlogs
- [ ] COULD add a top level Render that either uses the view configured in the context, or a global view if
            non is available. BUT 1: view is not request scoped, BUT 2: globals are an anti-pattern
- [ ] COULD add a renderer for the "session" that "renders" by saving the session to the response. It must
            be placed before regular renderers.
- [x] SHOULD move all template stuff to rbtemplate
- [ ] COULD allow for url path extensions to determine which encoder to use
- [ ] SHOULD do render logging on request scope if possible
- [ ] SHOULD test error paths, after having a nice error type to work with
- [ ] COULD add middleware that returns a request/response scoped render from the context