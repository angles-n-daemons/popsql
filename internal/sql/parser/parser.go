package parser

/*
statement       → insert | select | update | delete | create

create          → "CREATE TABLE" table "("
		   column_spec ( "," column_spec)*
                   ")"

select          → "SELECT" expression_list
                  ( "FROM" table_expr)?
                  ( "WHERE" logic_or)?
                  ( "GROUP BY" expression_list)?
                  ( "OFFSET" expression)?
                  ( "LIMIT" expression)?

insert          → "INSERT INTO " table parameters
                  ("VALUES" tuple ("," tuple)*) | select

column          → identifier type
type            → "integer" | "varchar" | "boolean"

expression      → or
or              → and "OR" and
and             → equality "AND" equality
equality        → comparison ( ( "!=" | "==" ) comparison)*
comparison      → term ( ( ">" | ">=" | "<" | "<=" ) term)*
term            → factor ( ( "-" | "+" ) factor )*
factor          → unary ( ( "/" | "*" ) unary)*
primary         → "true" | "false" | "nil" |
                  NUMBER | STRING | IDENTIFIER | "(" expression ")"

expression_list → expression  (","  expression)*
parameters      → identifier  (","  identifier)*
tuple           → "("  expression_list  ")"
*/
