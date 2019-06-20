;;;; ale core: builtin

(def-builtin -)
(def-builtin !=)
(def-builtin *)
(def-builtin /)
(def-builtin +)
(def-builtin <)
(def-builtin <=)
(def-builtin =)
(def-builtin >)
(def-builtin >=)

(def-builtin append)
(def-builtin apply)
(def-builtin car)
(def-builtin cdr)
(def-builtin chan)
(def-builtin cons)
(def-builtin current-time)
(def-builtin defer)
(def-builtin deque)
(def-builtin eq)
(def-builtin first)
(def-builtin gensym)
(def-builtin get)
(def-builtin go*)
(def-builtin lazy-seq*)
(def-builtin length)
(def-builtin list)
(def-builtin macro)
(def-builtin mod)
(def-builtin nth)
(def-builtin object)
(def-builtin promise)
(def-builtin raise)
(def-builtin read)
(def-builtin recover)
(def-builtin rest)
(def-builtin reverse)
(def-builtin seq)
(def-builtin str!)
(def-builtin str)
(def-builtin sym)
(def-builtin vector)

;; base types
(def-builtin is-apply)
(def-builtin is-boolean)
(def-builtin is-list)
(def-builtin is-number)
(def-builtin is-pair)
(def-builtin is-string)
(def-builtin is-symbol)
(def-builtin is-vector)

(def-builtin is-appender)
(def-builtin is-atom)
(def-builtin is-counted)
(def-builtin is-delivered)
(def-builtin is-deque)
(def-builtin is-empty)
(def-builtin is-indexed)
(def-builtin is-keyword)
(def-builtin is-local)
(def-builtin is-macro)
(def-builtin is-mapped)
(def-builtin is-nan)
(def-builtin is-neg-inf)
(def-builtin is-object)
(def-builtin is-pos-inf)
(def-builtin is-promise)
(def-builtin is-qualified)
(def-builtin is-reversible)
(def-builtin is-seq)
(def-builtin is-special)

(def-macro syntax-quote)

(def-special declare)
(def-special def)
(def-special do)
(def-special eval)
(def-special if)
(def-special lambda)
(def-special let)
(def-special letrec)
(def-special macroexpand-1)
(def-special macroexpand)
(def-special quote)
