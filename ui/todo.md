# MJ frontend todo

## Make eslint pass

This will ensure a basic level of not completely awfulness.
- Classify all things that were `createClass`
- propTypes defined

## Further improvements
- `PUT` endpoint needs to validate incoming requests better
- Would be nice if saved yaml had newlines between endpoints
- The UI currently implies you can have both a request body and a form, but that doesn't really make sense. Make the UI force user to fill one or the other.


## Tech debt
- Find `todo` marked in code
- Not much tests around new endpoints created
- No tests at all around the react stuff!
- `index.jsx` is a mess. Need to break out functionality into smaller pieces, such as the left menu.
- CSS is a mess, needs to be tidied. Is there a tool to remove unused things because there will be loads

### useful links

http://ricostacruz.com/cheatsheets/react.html
