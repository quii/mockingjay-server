# Implementing regex matching

## Done
- If a request has a regexURI configured then it will use it to match if the URI does not match

## To do
- Isn't working end-to-end! Need to parse YAML field into struct. Validation required:
  - Is a valid regex (regexp.Compile should cover this)
  - Is valid against the defined URI.
- With multiple requests defined you could run into a scenario where regex "wins" over an exact match, depending on the other. That should never happen.
- Acceptance test
- Update documentation

### Refactors
- The logic of "request matching" was already a bit too much for the "server", so it should be split into a dependency
