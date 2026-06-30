package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
)

type Config struct {
	DSN    string
	Schema string
	Dest   string
}

func ParseFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.DSN, "dsn", "", "PostgreSQL DSN (required)")
	flag.StringVar(&cfg.Schema, "schema", "./schema", "Path to schema directory")
	flag.StringVar(&cfg.Dest, "dest", "./gen", "Output directory for generated Jet code")

	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	if cfg.DSN == "" {
		fmt.Println("error: --dsn is required")
		flag.Usage()
		panic("missing dsn")
	}

	if *verbose {
		fmt.Printf("config: %+v\n", cfg)
	}

	return cfg
}

func main() {
	cfg := ParseFlags()
	err := postgres.GenerateDSN(cfg.DSN, cfg.Schema, cfg.Dest,
		template.Default(postgres2.Dialect).UseSchema(func(schemaMetadata metadata.Schema) template.Schema {
			return template.DefaultSchema(schemaMetadata).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						tableModel := template.DefaultTableModel(table)
						if table.Name == "spatial_ref_sys" {
							tableModel.Skip = true
						}
						return tableModel.
							UseField(func(column metadata.Column) template.TableModelField {
								if strings.ToLower(column.DataType.Name) == "citext" {
									column.DataType.Name = "text"
								}
								if column.DataType.Name == "geometry" {
									return template.TableModelField{
										Skip: true,
									}
								}

								return template.DefaultTableModelField(column)
							})
					}).
					UseView(func(table metadata.Table) template.TableModel {
						return template.ViewModel{Skip: true}
					}),
				).
				UseSQLBuilder(template.DefaultSQLBuilder().
					UseTable(func(table metadata.Table) template.TableSQLBuilder {
						tableModel := template.DefaultTableSQLBuilder(table)
						if table.Name == "spatial_ref_sys" {
							tableModel.Skip = true
						}
						return tableModel.
							UseColumn(func(column metadata.Column) template.TableSQLBuilderColumn {
								if strings.ToLower(column.DataType.Name) == "citext" {
									column.DataType.Name = "text"
								}
								if column.DataType.Name == "geometry" {
									return template.TableSQLBuilderColumn{
										Skip: true,
									}
								}

								return template.DefaultTableSQLBuilderColumn(column)
							})
					}).
					UseView(func(table metadata.Table) template.TableSQLBuilder {
						return template.TableSQLBuilder{Skip: true}
					}),
				)
		}),
	)

	if err != nil {
		panic(err)
	}
}
