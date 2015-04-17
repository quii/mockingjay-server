# jsonequaliser

Checks that two json strings (A and B) are compatible or not.

See the tests for examples as to how this works.

## What is compatability?

If A has a field then the code expects B to have that same field with a value which corresponds to the same type. The code does not care about what actual data B or A has.

When it comes to arrays, B will need to have at least one item in it's own version so it can be proven it's compatible.

If B has extra fields they will still be seen as compatible.

## Use cases

For help with implementing CDCs.

You can run this against a fake server and the real API you are using to ensure your test code and downstream services produce the JSON you need.