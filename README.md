# ReverseRegex

The goal here is to have a function, that given a regex pattern and some strings, will return a string that is populated with these strins, and is matched by given regex pattern.

So, for example:
```
rRx = newReverseRegexp(`/article/(?P<id>\d)-(?P<slug>.*)`)
println(rRx.getReversedString("15", "title-or-something"))
```
should print the string `/article/15-title-or-something`.
