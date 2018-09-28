package server

import (
	"database/sql"

	"gopkg.in/gorp.v2"

	"github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Client is struct that holds database connection
type Client struct {
	db    *sql.DB
	dbmap *gorp.DbMap
}

type Conversions struct {
	Programs        map[string]Program
	DoorModels      map[string]DoorModel
	Hinges          map[byte]Hinge
	Handednesses    map[byte]Handedness
	Handles         map[byte]Handle
	HandlePositions map[byte]HandlePosition
}

type Program struct {
	ID            int64   `db:"id"`
	Name          string  `db:"name"`
	Program       string  `db:"program"`
	HingePosition float64 `db:"hinge_position"`
	SlateHinge    int     `db:"slate_hinge"`
}

type DoorModel struct {
	ID            int64   `db:"id"`
	Name          string  `db:"name"`
	Depth         float64 `db:"depth"`
	Stopper       int     `db:"stopper"`
	SlatePosition int     `db:"slate_position"`
}

type Hinge struct {
	ID      int64  `db:"id"`
	Barcode string `db:"barcode"`
	Var5    int    `db:"variable"`
}

type Handedness struct {
	ID         int64  `db:"id"`
	Barcode    string `db:"barcode"`
	Handedness string `db:"handedness"`
}

type Handle struct {
	ID      int64  `db:"id"`
	Barcode string `db:"barcode"`
	Handle  int    `db:"handle"`
}

type HandlePosition struct {
	ID       int64  `db:"id"`
	Barcode  string `db:"barcode"`
	Position string `db:"handle_position"`
}

// NewClient opens dtabase connection and returns client struct
// which wraps connections
func NewClient(dbUrl string) (c *Client, err error) {
	connStr, err := pq.ParseURL(dbUrl)
	if err != nil {
		return
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}

	// Ensure DB version and run needed migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return
	}
	err = nil

	c = &Client{
		db: db,
		dbmap: &gorp.DbMap{
			Db:      db,
			Dialect: gorp.PostgresDialect{},
		},
	}

	c.dbmap.AddTableWithName(Hinge{}, "hinges").SetKeys(true, "ID").SetKeys(false, "Barcode")
	c.dbmap.AddTableWithName(Handle{}, "handles").SetKeys(true, "ID").SetKeys(false, "Barcode")
	c.dbmap.AddTableWithName(Program{}, "programs").SetKeys(true, "ID").SetKeys(false, "Name")
	c.dbmap.AddTableWithName(DoorModel{}, "door_models").SetKeys(true, "ID").SetKeys(false, "Name")
	c.dbmap.AddTableWithName(Handedness{}, "handednesses").SetKeys(true, "ID").SetKeys(false, "Barcode")
	c.dbmap.AddTableWithName(HandlePosition{}, "handle_positions").SetKeys(true, "ID").SetKeys(false, "Barcode")
	return
}

func (c *Client) UpdateHandlePositions(rows []HandlePosition) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) UpdateHandednesses(rows []Handedness) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) UpdateProgram(rows []Program) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) UpdateDoorModels(rows []DoorModel) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) UpdateHinges(rows []Hinge) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) UpdateHandles(rows []Handle) (err error) {
	tx, err := c.dbmap.Begin()
	if err != nil {
		return err
	}

	for _, row := range rows {
		if _, err = tx.Update(&row); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *Client) AddHandlePositions(row *HandlePosition) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) AddHandednesses(row *Handedness) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) AddProgram(row *Program) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) AddDoorModels(row *DoorModel) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) AddHinges(row *Hinge) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) AddHandles(row *Handle) (err error) {
	return c.dbmap.Insert(row)
}

func (c *Client) GetPrograms() (programs map[string]Program, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT name, id, program, hinge_position, slate_hinge FROM programs WHERE program NOTNULL AND hinge_position NOTNULL ORDER BY name"); err != nil {
		return
	}
	defer rows.Close()

	programs = make(map[string]Program)
	for rows.Next() {
		name := ""
		pd := Program{}
		if err = rows.Scan(&name, &pd.ID, &pd.Program, &pd.HingePosition, &pd.SlateHinge); err != nil {
			return
		}
		programs[name] = pd
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (c *Client) GetDoorModels() (doorModels map[string]DoorModel, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT name, id, depth, stopper, slate_position FROM door_models WHERE depth NOTNULL AND stopper NOTNULL ORDER BY name"); err != nil {
		return
	}
	defer rows.Close()

	doorModels = make(map[string]DoorModel)
	for rows.Next() {
		name := ""
		dm := DoorModel{}
		if err = rows.Scan(&name, &dm.ID, &dm.Depth, &dm.Stopper, &dm.SlatePosition); err != nil {
			return
		}
		doorModels[name] = dm
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (c *Client) GetHinges() (hinges map[byte]Hinge, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT barcode, id, variable FROM hinges WHERE barcode NOTNULL AND variable NOTNULL ORDER BY barcode"); err != nil {
		return
	}
	defer rows.Close()

	hinges = make(map[byte]Hinge)
	for rows.Next() {
		barcode := make([]byte, 1)
		hinge := Hinge{}
		if err = rows.Scan(&barcode, &hinge.ID, &hinge.Var5); err != nil {
			return
		}
		hinges[barcode[0]] = hinge
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (c *Client) GetHandednesses() (handednesses map[byte]Handedness, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT barcode, id, handedness FROM handednesses WHERE barcode NOTNULL AND handedness NOTNULL ORDER BY barcode"); err != nil {
		return
	}
	defer rows.Close()

	handednesses = make(map[byte]Handedness)
	for rows.Next() {
		barcode := make([]byte, 1)
		handedness := Handedness{}
		if err = rows.Scan(&barcode, &handedness.ID, &handedness.Handedness); err != nil {
			return
		}
		handednesses[barcode[0]] = handedness
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (c *Client) GetHandles() (handles map[byte]Handle, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT barcode, id, handle FROM handles WHERE barcode NOTNULL AND handle NOTNULL ORDER BY barcode"); err != nil {
		return
	}
	defer rows.Close()

	handles = make(map[byte]Handle)
	for rows.Next() {
		barcode := make([]byte, 1)
		handle := Handle{}
		if err = rows.Scan(&barcode, &handle.ID, &handle.Handle); err != nil {
			return
		}
		handles[barcode[0]] = handle
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (c *Client) GetHandlePositions() (handlePositions map[byte]HandlePosition, err error) {
	var rows *sql.Rows
	if rows, err = c.db.Query("SELECT barcode, id, handle_position FROM handle_positions WHERE barcode NOTNULL AND handle_position NOTNULL ORDER BY barcode"); err != nil {
		return
	}
	defer rows.Close()

	handlePositions = make(map[byte]HandlePosition)
	for rows.Next() {
		barcode := make([]byte, 1)
		handlePosition := HandlePosition{}
		if err = rows.Scan(&barcode, &handlePosition.ID, &handlePosition.Position); err != nil {
			return
		}
		handlePositions[barcode[0]] = handlePosition
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (c *Client) GetAll() (all *Conversions, err error) {
	programs, err := c.GetPrograms()
	if err != nil {
		return
	}
	doorModels, err := c.GetDoorModels()
	if err != nil {
		return
	}
	hinges, err := c.GetHinges()
	if err != nil {
		return
	}
	handednesses, err := c.GetHandednesses()
	if err != nil {
		return
	}
	handles, err := c.GetHandles()
	if err != nil {
		return
	}
	handlePositions, err := c.GetHandlePositions()
	if err != nil {
		return
	}
	all = &Conversions{
		Programs:        programs,
		DoorModels:      doorModels,
		Hinges:          hinges,
		Handednesses:    handednesses,
		Handles:         handles,
		HandlePositions: handlePositions,
	}
	return
}
