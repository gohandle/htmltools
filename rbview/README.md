# Backlogs
- [ ] COULD add a top level Render that either uses the view configured in the context, or a global view if
            non is available. BUT 1: view is not request scoped, BUT 2: globals are an anti-pattern
- [ ] SHOULD add a renderer for the "session" that "renders" by saving the session to the response. It must
            be placed before regular renderers.
- [ ] SHOULD move all template stuff to rbtemplate