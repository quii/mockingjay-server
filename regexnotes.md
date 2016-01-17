# Implementing regex matching

## Done
- Parsing regex field in YAML
- If a request has a regexURI configured then it will use it to match if the URI does not match

## To do
- Check regex is valid against the defined URI.
- With multiple requests defined you could run into a scenario where regex "wins" over an exact match, depending on the other. That should never happen.
- With multiple regexURIs a given request might match more than one. That should result in a conflict.
- Acceptance test
- Update documentation

### Refactors
- The logic of "request matching" was already a bit too much for the "server", so it should be split into a dependency
