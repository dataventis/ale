;;;; ale bootstrap: builtin

(def-builtin do)
(def-builtin if)
(def-builtin let)
(def-builtin fn)

(def-builtin read)
(def-builtin eval)
(def-builtin is-eq)

;; globals

(def-builtin declare)
(def-builtin def)

;; basic predicates

(def-builtin is-nil)
(def-builtin is-atom)
(def-builtin is-keyword)

;; macros

(def-builtin defmacro)
(def-builtin quote)
(def-builtin syntax-quote)
(def-builtin macroexpand-1)
(def-builtin macroexpand)
(def-builtin is-macro)

;; symbols

(def-builtin sym)
(def-builtin gensym)
(def-builtin is-symbol)
(def-builtin is-local)
(def-builtin is-qualified)

;; strings

(def-builtin str)
(def-builtin str!)
(def-builtin is-str)

;; sequences

(def-builtin seq)
(def-builtin first)
(def-builtin rest)
(def-builtin last)
(def-builtin cons)
(def-builtin conj)
(def-builtin len)
(def-builtin nth)
(def-builtin get)
(def-builtin assoc)
(def-builtin list)
(def-builtin vector)

(def-builtin is-seq)
(def-builtin is-len)
(def-builtin is-indexed)
(def-builtin is-assoc)
(def-builtin is-mapped)
(def-builtin is-list)
(def-builtin is-vector)

;; numeric

(def-builtin +)
(def-builtin -)
(def-builtin *)
(def-builtin /)
(def-builtin mod)

(def-builtin =)
(def-builtin !=)
(def-builtin >)
(def-builtin >=)
(def-builtin <)
(def-builtin <=)

(def-builtin is-pos-inf)
(def-builtin is-neg-inf)
(def-builtin is-nan)

;; functions

(def-builtin partial)
(def-builtin apply)
(def-builtin is-apply)
(def-builtin is-special)

;; concurrency

(def-builtin go*)
(def-builtin chan)
(def-builtin promise)
(def-builtin is-promise)

;; lazy sequences

(def-builtin lazy-seq*)
(def-builtin concat)
(def-builtin filter)
(def-builtin map)
(def-builtin take)
(def-builtin drop)
(def-builtin reduce)
(def-builtin for-each*)

;; raise and recover

(def-builtin raise)
(def-builtin recover)

;; current time

(def-builtin current-time)
