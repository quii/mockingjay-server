# MJ frontend todo

## make the structure of the front end suck less
- Learn webpack :/
- dont commit generated code, obviously

## Further improvements
- The UI currently implies you can have both a request body and a form, but that doesn't really make sense. Make the UI force user to fill one or the other.
- Proptypes for key areas, such as the endpoint data passed into the renderer
- This? https://github.com/yelouafi/redux-saga

## Tech debt
- UI isn't fully tested
  - Proved basic unit test so now lets test the UI layer in terms of stuff appearing and its interactions with the API first
- CSS is a mess, needs to be tidied. Is there a tool to remove unused things because there will be loads
