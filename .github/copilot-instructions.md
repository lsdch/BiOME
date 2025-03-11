We use Gel as our database manager, so when suggesting database queries, use EdgeQL syntax. When the query is embedded in Go code, always use backticks string delimiters and prefix with #edgeql.
Example:
`#edgeql
{ query goes here }
`

Database queries must use the database schema defined in the files db/schema/*.esdl