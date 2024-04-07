package testdata

type ScopedStruct struct {
	Name       string
	NoDB       bool `blaze:"no-db"`
	Read       bool `blaze:"read"`
	ReadCreate bool `blaze:"read.create"`
	ReadUpdate bool `blaze:"read.update"`
	Update     bool `blaze:"update"`
	Create     bool `blaze:"create"`
	Write      bool `blaze:"write"`

	NoClient         bool `blaze:"client:-"`
	ClientRead       bool `blaze:"client:read"`
	ClientReadCreate bool `blaze:"client:read.create"`
	ClientReadUpdate bool `blaze:"client:read.update"`
	ClientUpdate     bool `blaze:"client:update"`
	ClientCreate     bool `blaze:"client:create"`
	ClientWrite      bool `blaze:"client:write"`

	NoAdmin         bool `blaze:"admin:-"`
	AdminRead       bool `blaze:"admin:read"`
	AdminReadCreate bool `blaze:"admin:read.create"`
	AdminReadUpdate bool `blaze:"admin:read.update"`
	AdminUpdate     bool `blaze:"admin:update"`
	AdminCreate     bool `blaze:"admin:create"`
	AdminWrite      bool `blaze:"admin:write"`
}
