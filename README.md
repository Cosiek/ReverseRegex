# ReverseRegex

The goal here is to have a function, that given a regex pattern, and some values to populate it with, will return a string that is populated with hese values, and is matched by given regex pattern.

So, for example:
```
revRx := newReverseRegex("/api/((?P<id>\d+)/edit")
revRx.getReversedString([]string{"123"})
```
should give the string `/api/123/edit`
