# htmltools
Tooling for frontend development using html/template

# Components
- [x] SHOULD add an rblog package that provides logging with request scoped logger
- [x] SHOULD add an rbview package that load templates
- [x] SHOULD add an rbi18n that read bundles from filesystem or embeds and provides template helper
- [x] SHOULD add a rbsess that does session saving (JIT?), CSRF and flash functionality (middleware, based)
- [x] SHOULD add an rbasset that has helpers for serving assets with cache busting/hashes. And contains a static filesystem http.handler that can be included in the handler (middleware). Maybe allow decoding of session as well.
- [x] SHOULD add a package for form binding, possibly with middleware that does it based on content-type. Use fx group functionality to allow configuration of arbitrary binding logic
- [x] SHOULD extend the view package to support wider rendering, maybe middleware that looks at accept language and injects a request scoped render looking at accept headers. Configuration through fx's group feature to configure supported response renderings.
- [x] SHOULD have a package holds the routing/url generation helper. 
- [ ] SHOULD have a package that makes it easy to create handlers that return errors. It may also render errors (and request logs) very nicely in development. Possibly provide debug information about all components.
- [ ] COULD add a package that provides live-reload experience for assets, templates and go code
- [ ] COULD add a package that manages form rendering/creation
- [ ] COULD add struct based validation package
- [ ] COULD add a package rbsql for sql connections for di
- [ ] COULD develop a formatter that formats html with go templating
- [ ] COULD develop a html validator that uses the checks if the html is valid
- [ ] COULD add a package that helps with bundling and compiling using esbuild. Build on embed
- [ ] COULD create a type that represents a framework context and can be embedded in page types
- [ ] SHOULD have a mode where any error that is logged should be clearly visible on the response
- [ ] COULD allow all helpers to have there name prefixed with a configurable prefix in rbview conf
- [ ] COULD explore a pattern of server side "components" that make use of the new embed feature to include template files close to go code. It should be easy to build up page structs from these components
which are basically structs that render.