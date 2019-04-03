;;;; ale bootstrap: basics

(defmacro defn
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defmacro error
  [& clauses]
  `(raise (assoc ~@clauses)))

(defmacro panic
  [& clauses]
  `(raise (error ~@clauses)))

(defmacro eq
  [value & comps]
  `(is-eq ~value ~@comps))

(defmacro !eq
  [value & comps]
  `(not (is-eq ~value ~@comps)))

(defn is-even
  [value]
  (= (mod value 2) 0))

(defn is-odd
  [value]
  (= (mod value 2) 1))

(defn inc
  [value]
  (+ value 1))

(defn dec
  [value]
  (- value 1))
