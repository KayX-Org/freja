package component

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/kayx-org/freja/env"
)

type OptionArango func(*Arango)
type IndexType int

const (
	ErrDuplicate = 1207
)

type Index interface {
	Name() string
	Fields() []string
}

type GeoIndex struct {
	IxName   string
	IxFields []string
}

func (g GeoIndex) Name() string {
	return g.IxName
}

func (g GeoIndex) Fields() []string {
	return g.IxFields
}

type HashIndex struct {
	IxName   string
	IxFields []string
	Unique   bool
	Sparse   bool
}

func (g HashIndex) Name() string {
	return g.IxName
}

func (g HashIndex) Fields() []string {
	return g.IxFields
}

type Collection struct {
	Name    string
	Indexes []Index
}

type Edge struct {
	Name    string
	Indexes []Index
	// From: contains the names of one or more vertex collections that can contain source vertices.
	// To: contains the names of one or more edge collections that can contain target vertices
	From []string
	To   []string
}

type Graph struct {
	Name     string
	Edges    []Edge
	Vertexes []Collection
}

type Arango struct {
	db          string
	user        string
	password    string
	endpoints   []string
	graph       *Graph
	collections []Collection
	clientDB    driver.Database
	client      driver.Client
}

func NewArango(db string, endpoints []string, user, password string, options ...OptionArango) *Arango {
	a := &Arango{
		db:          db,
		user:        user,
		password:    password,
		endpoints:   endpoints,
		collections: make([]Collection, 0),
	}
	for _, o := range options {
		o(a)
	}

	return a
}

func OptionGraph(graph Graph) OptionArango {
	return func(a *Arango) {
		a.graph = &graph
	}
}

func OptionCollections(c []Collection) OptionArango {
	return func(a *Arango) {
		if c != nil {
			a.collections = c
		}
	}
}

func (a *Arango) ClientDB(ctx context.Context) (driver.Database, error) {
	if a.clientDB != nil {
		return a.clientDB, nil
	}

	if a.client == nil {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: a.endpoints,
			ConnLimit: env.GetEnvAsInt("ARANGODB_CONN_LIMIT", 32),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to create connection: %w", err)
		}

		a.client, err = driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(a.user, a.password),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to connect: %w", err)
		}
	}

	return a.createDB(ctx)
}

func (a *Arango) createDB(ctx context.Context) (driver.Database, error) {
	true := true
	connDB, err := a.client.CreateDatabase(ctx, a.db, &driver.CreateDatabaseOptions{
		Users: []driver.CreateDatabaseUserOptions{
			{UserName: a.user, Password: a.password, Active: &true},
		},
		Options: driver.CreateDatabaseDefaultOptions{},
	})
	a.clientDB, err = a.processCreateDBError(ctx, connDB, err)
	if err != nil {
		return nil, err
	}

	if err := a.createGraph(ctx); err != nil {
		return nil, err
	}

	if err := a.createCollections(ctx); err != nil {
		return nil, err
	}

	return a.clientDB, nil
}

func (a *Arango) processCreateDBError(ctx context.Context, dbClient driver.Database, err error) (driver.Database, error) {
	if err != nil {
		return dbClient, nil
	}

	if driver.IsArangoErrorWithErrorNum(err, ErrDuplicate) {
		if c, err := a.client.Database(ctx, a.db); err != nil {
			return nil, fmt.Errorf("unable to get db: %w", err)
		} else {
			return c, nil
		}
	}

	return nil, fmt.Errorf("unable to create db: %w", err)
}

func (a *Arango) createGraph(ctx context.Context) error {
	if a.graph == nil {
		return nil
	}

	g, err := a.clientDB.CreateGraph(ctx, a.graph.Name, &driver.CreateGraphOptions{ReplicationFactor: 2})
	g, err = a.processCreateGraphError(ctx, g, err)
	if err != nil {
		return err
	}

	for _, v := range a.graph.Vertexes {
		c, err := g.CreateVertexCollection(ctx, v.Name)
		if err != nil {
			return fmt.Errorf("unable to create vertex collection '%s', :%w", v.Name, err)
		}
		if err := a.ensureIndexes(ctx, c, v.Indexes); err != nil {
			return err
		}
	}

	for _, e := range a.graph.Edges {
		c, err := g.CreateEdgeCollection(ctx, e.Name, driver.VertexConstraints{})
		if err != nil {
			return fmt.Errorf("unable to create edge collection '%s', :%w", e.Name, err)
		}
		if err := a.ensureIndexes(ctx, c, e.Indexes); err != nil {
			return err
		}
		if err := g.SetVertexConstraints(ctx, e.Name, driver.VertexConstraints{
			From: e.From,
			To:   e.To,
		}); err != nil {
			return fmt.Errorf("unable to set vertex constrainse for edge '%s', :%w", e.Name, err)
		}
	}

	return nil
}

func (a *Arango) createCollections(ctx context.Context) error {
	if a.collections == nil {
		return nil
	}

	for _, col := range a.collections {
		c, err := a.clientDB.CreateCollection(ctx, col.Name, &driver.CreateCollectionOptions{})
		if err != nil {
			return fmt.Errorf("unable to create collection '%s', :%w", c.Name, err)
		}
		if err := a.ensureIndexes(ctx, c, col.Indexes); err != nil {
			return err
		}
	}

	return nil
}

func (a *Arango) ensureIndexes(ctx context.Context, c driver.Collection, indexes []Index) error {
	for _, ix := range indexes {
		switch index := ix.(type) {
		case GeoIndex:
			_, _, err := c.EnsureGeoIndex(ctx, ix.Fields(), &driver.EnsureGeoIndexOptions{
				// Long and then latitude
				GeoJSON:      true,
				InBackground: true,
				Name:         ix.Name(),
			})

			if err != nil {
				return fmt.Errorf("unable to create geo index '%s': %w", ix.Name, err)
			}
		case HashIndex:
			_, _, err := c.EnsureHashIndex(ctx, ix.Fields(), &driver.EnsureHashIndexOptions{
				Unique:        index.Unique,
				Sparse:        index.Sparse,
				NoDeduplicate: false,
				InBackground:  true,
				Name:          ix.Name(),
			})

			if err != nil {
				return fmt.Errorf("unable to create geo index '%s': %w", ix.Name, err)
			}
		default:
			return fmt.Errorf("unhandled index type")
		}
	}

	return nil
}

func (a *Arango) processCreateGraphError(ctx context.Context, graph driver.Graph, err error) (driver.Graph, error) {
	if err != nil {
		return graph, nil
	}

	if driver.IsArangoErrorWithErrorNum(err, ErrDuplicate) {
		if c, err := a.clientDB.Graph(ctx, a.graph.Name); err != nil {
			return nil, fmt.Errorf("unable to get graph: %w", err)
		} else {
			return c, nil
		}
	}

	return nil, fmt.Errorf("unable to create graph: %w", err)
}
