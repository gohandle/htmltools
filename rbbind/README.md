# Backlog
- [ ] SHOULD add xml decoder
- [ ] COULD add debug logging for binding, but that should be request scoped, which requires inclusion of other package


- [ ] COULD add a way to decode/bind data from the session into a request/page struct. But that might allow accidental writing of page fields from other request inputs if the fields have the same name. Instead, maybe add a rbbind.Target type that configures what (form, query, sess) will be decoded with what. Or maybe look at the content type

Binding targets
- Body: Based on content-type
- Query: Always form
- Session: implemented by rbsess
- Headers

## Bind v2
Features:
- Zero-alloc if the encoder supports zero allocations (such as form query decoding). This prevents us from
  using the option pattern. Maybe flags, or pair providing
- Should support decoders that decode the: body, query, session or headers
- Should support taking into account content-type