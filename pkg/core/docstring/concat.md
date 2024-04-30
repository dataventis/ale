---
title: "concat"
description: "concatenates sequences"
names: ["concat"]
usage: "(concat seq+) (concat! seq+)"
tags: ["sequence", "comprehension"]
---

Creates a lazy sequence whose content is the result of concatenating the elements of each provided sequence. To immediately materialize a complete concatenated sequence, use the `concat!` function.

#### An Example

```scheme
(seq->list (concat [1 2 3] '(4 5 6)))
```

This will return the list _(1 2 3 4 5 6)_
