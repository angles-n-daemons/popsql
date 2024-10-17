package schema_test

import "testing"

func TestColumnBasic(t *testing.T) {
}

/*
Feature breakdown:
 - column with space / table embedded?



In the catalog, we have a schema, table and column def:
- Should the column have Space and Table strings?
- Should the schema / tables have table and col references keyed by ID or name?
    - Name prevents duplication
- Should the catalog be responsible for its own persistence?
- Should There be a register type? Should the table and columns implement it?

What questions need I answer:
- What does the column have?
- If it has a table, it can have a Key function
- but then it also has a circular ref table -> col -> table
- I'm not sure that I like it, so it seems like there are three options:
   1. Include both space and table strings on the col
   2. Include the table on the col
   3. Get rid of a notion of having the column itself act as a serializable thing

// I like having a new function, allows for validation of name
// I like also having the ability to consider this a register

// Alas a new  issue crops up
// to load a column, I require a table

// Columns themselves cannot have a FromBytes abstraction iff they store a table object on themselves
// they cannot recreate the table

Do I need an order this early?


So what's bad about having the space / table combo on the column.

Duplicated data, opportunity for error proneness on deserialization (eg the table is renamed)?

Rather what is the push and pull between runtime representation and serialization?

and then another core issue becomes that anything which needs to be keyed, then becomes

autoincrement is specifically challenging for a number of reasons:
- it needs to be locked, kept safe
- its most up to date value needs to be persisted to disk, it serves no benefit to read the largest value each time the service starts up


an autogen primary key will in itself be useful

what types of records will then theoretically be serialized?:
- column
- table
- table rows
- table statistics
- role assignments
- settings (hrm)

How to deserialize the above? Is a schema necessary?
Possibly to pull the appropriate columns out

if each of these had a key and a value func?

Could each of them be bundled with a record?

How should writes be handled by the sql engine?

Keyable object, Value

The key should be joined with the table prefix on write. In this way, it shouldn't need to know the table prefix for the key

And now here comes the second challenging question.

How should the schema manage tables?

map[ID]Table? (mix system and user tables alike?)
I am beginning to feel more amenable to this idea

Table Columns = map[string]Table?

What's the performance gain of a small map vs a small array?

And so tables and columns are glue-like runtime types. They are responsible not only for the functionality which the engine depends on, keying information, structure, naming, typing. But also for their own ability to transform to and from registers.

They are not responsible for persisting their own changes

Or should they be?

Should the catalog be responsible for storage persistence as well?

I feel no closer to an answer, at a philosophical standstill like old ryan pal. Can I at least summarize some of the decisions to be made next time?

If the catalog takes care of its own schema management, it lifts a lot of the burden off the engine.

























*/
