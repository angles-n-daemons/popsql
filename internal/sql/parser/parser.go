package parser

/*
statement       → insert | select | update | delete

select          → "SELECT" + expression_list +
                  (+ "FROM"  table_expr)?
                  (+ "WHERE" logic_or)?
                  (+ "GROUP BY" expression_list)?
                  (+ "OFFSET" expression)?
                  (+ "LIMIT" expression)?

insert          → "INSERT INTO " + table + parameters
                  ("VALUES" + tuple + ("," + tuple)*) |
                  select

expression      → or
or              → and + "OR" + and
and             → equality + "AND" + equality
equality
comparison

expression_list → expression + ("," + expression)*
parameters      → identifier + ("," + identifier)*
tuple           → "(" + expression_list + ")"
*/
