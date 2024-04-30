---
title: "let"
description: "binds local values"
names: ["let", "let*"]
usage: "(let ([name expr]*) form*) (let* ([name expr]*) form*)"
tags: ["binding"]
---

Create a new local scope, evaluate the provided expressions, and then bind the resulting value to their respective names. It will then evaluate the specified forms within that scope and return the result of the last evaluation. The `let` form performs these bindings in parallel, whereas the `let*` form performs them sequentially.

#### An Example

```scheme
(let ([x '(1 2 3 4)]
      [y [5 6 7 8] ])
  (concat x y))
```

This example will create a list called _x_ and a vector called _y_ and return the lazy concatenation of those sequences. Note that the two names do not exist outside the `let` form.
