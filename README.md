# htmltools
Tooling for frontend development using html/template

# Components
- [x] SHOULD add an rblog package that provides logging with request scoped logger
- [x] SHOULD add an rbview package that load templates
- [x] SHOULD add an rbi18n that read bundles from filesystem or embeds and provides template helper
- [x] SHOULD add a rbsess that does session saving (JIT?), CSRF and flash functionality (middleware, based)
- [x] SHOULD add an rbasset that has helpers for serving assets with cache busting/hashes. And contains a static filesystem http.handler that can be included in the handler (middleware). Maybe allow decoding of session as well.
- [ ] SHOULD add a package for form (binding?) decoding/validation, possibly with middleware that does it based on content-type. Use fx group functionality to allow configuration of arbitrary binding logic
- [ ] SHOULD extend the view package to support wider rendering, maybe middleware that looks at accept language and injects a request scoped render looking at accept headers. Configuration through fx's group feature to configure supported response renderings.
- [ ] COULD add struct based validation package
- [ ] COULD add a package rbsql for sql connections for di
- [ ] SHOULD have a package holds the routing/url generation helper. 
- [ ] COULD develop a formatter that formats html with go templating
- [ ] COULD develop a html validator that uses the checks if the html is valid
- [ ] COULD add a package that helps with bundling and compiling using esbuild. Build on embed