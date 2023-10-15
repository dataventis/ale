---
title: "cons"
description: "combines two values into a pair"
names: ["cons"]
usage: "(cons car cdr)"
tags: ["sequence"]
---

When `cdr` is an ordered sequence, such as a list or vector, the result is a new list or vector with the `car` value prepended to the original. With an unordered sequence, such as an object array, there is no guarantee regarding position. If `cdr` is not a sequence, then a new cons cell will be constructed.

The name _cons_ is a vestige of when Lisp implementations constructed new lists or cells by pairing a _car_ (contents of the address part of register) with a _cdr_ (contents of the decrement part of register).

#### An Example

```scheme
(define x '(3 4 5 6))
(define y (cons 2 x))
(cons 1 y)
```

This example will return _(1 2 3 4 5 6)_.
