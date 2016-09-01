# MJ frontend todo

- Doesn't work when config has regex because it doesnt know how to parse it.

## make the structure of the front end suck less
- Learn webpack :/
- dont commit generated code, obviously

## Further improvements
- The UI currently implies you can have both a request body and a form, but that doesn't really make sense. Make the UI force user to fill one or the other.
- Proptypes for key areas, such as the endpoint data passed into the renderer
- This? https://github.com/yelouafi/redux-saga

## Tech debt
- No tests at all around the react stuff!
  - Using Jest right now but maybe https://github.com/airbnb/enzyme is better
- CSS is a mess, needs to be tidied. Is there a tool to remove unused things because there will be loads

### useful links

http://ricostacruz.com/cheatsheets/react.html
