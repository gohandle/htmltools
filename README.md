# htmltools
Tooling for frontend development using html/template

# Components
- [x] SHOULD add an rblog package that provides logging with request scoped logger
- [x] SHOULD add an rbview package that load templates
- [x] SHOULD add an rbi18n that read bundles from filesystem or embeds and provides template helper
- [ ] SHOULD add a rbsess that does session saving (JIT?), CSRF and flash functionality (middleware, based)
- [ ] SHOULD add an rbasset that has helpers for serving assets with cache busting/hashes. And contains a static filesystem http.handler that can be included in the handler (middleware)
- [ ] SHOULD add a package for form (binding?) decoding/validation, possibly with middleware that does it based on content-type
- [ ] SHOULD extend the view package to support wider rendering, maybe middleware that looks at accept language and injects a request scoped render.
- [ ] SHOULD have a package holds the routing/url generation helper. 
- [ ] COULD add a package that helps with bundling and compiling using esbuild. Build on embed
- [ ] COULD develop a formatter that formats html with go templating
- [ ] COULD develop a html validator that uses the checks if the html is valid