package token

import (
	"github.com/gocql/gocql"
	"github.com/nalej/authx/internal/app/authx/entities"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

const table = "tokens"
const tablePK_1 = "username"
const tablePK_2 = "token_id"

const rowNotFound = "not found"

type ScyllaTokenProvider struct {
	Address string
	Port int
	KeySpace string
	Session *gocql.Session
}


func NewScyllaTokenProvider (address string, port int, keyspace string) *ScyllaTokenProvider{
	provider := ScyllaTokenProvider{Address:address, Port:port, KeySpace:keyspace}
	provider.Connect()
	return &provider
}

func (sp *ScyllaTokenProvider) Connect() derrors.Error {
	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.KeySpace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaTokeProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session
	return nil
}

func (sp *ScyllaTokenProvider) Disconnect()  {
	if sp != nil {
		sp.Session.Close()
	}
}

func (sp *ScyllaTokenProvider) CheckConnectionAndConnect () derrors.Error {

	if sp.Session != nil {
		return nil
	}
	log.Info().Str("provider", "ScyllaTokeProvider"). Msg("session not connected, trying to connect it!")
	err := sp.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (sp *ScyllaTokenProvider) Delete(username string, tokenID string) derrors.Error{
	if err := sp.CheckConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.Exist(username, tokenID)

	if err != nil {
		return err
	}
	if ! *exists {
		return derrors.NewNotFoundError("role").WithParams(username, tokenID)
	}

	stmt, _ := qb.Delete(table).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	cqlErr := sp.Session.Query(stmt, username, tokenID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete token")
	}

	return nil
}

// Add a token.
func (sp *ScyllaTokenProvider) Add(token *entities.TokenData) derrors.Error{

	if err := sp.CheckConnectionAndConnect(); err != nil {
		return err
	}

	exists, err := sp.Exist(token.Username, token.TokenID)
	if err != nil {
		return err
	}
	if *exists {
		return  derrors.NewAlreadyExistsError("token").WithParams(token.Username, token.TokenID)
	}

	// add new basic credential
	stmt, names := qb.Insert(table).Columns("username", "token_id", "refresh_token", "expiration_date").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(token)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add token")
	}

	return nil
}

// Get an existing token.
func (sp *ScyllaTokenProvider) Get(username string, tokenID string) (*entities.TokenData, derrors.Error){

	if err := sp.CheckConnectionAndConnect(); err != nil {
		return nil, err
	}

	var token entities.TokenData
	stmt, names := qb.Select(table).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: username,
		tablePK_2: tokenID})

	err := q.GetRelease(&token)
	if err != nil {
		if err.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError("token").WithParams(username, tokenID)
		}else{
			return nil, derrors.AsError(err, "cannot get token")
		}
	}

	return &token, nil
}

// Exist checks if the token was added.
func (sp *ScyllaTokenProvider) Exist(username string, tokenID string) (*bool, derrors.Error){

	ok := false

	if err := sp.CheckConnectionAndConnect(); err != nil {
		return &ok, err
	}

	var returnedId string

	stmt, names := qb.Select(table).Columns(tablePK_1).Where(qb.Eq(tablePK_1)).Where(qb.Eq(tablePK_2)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		tablePK_1: username,
		tablePK_2: tokenID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return &ok, nil
		}else{
			return &ok, derrors.AsError(err, "cannot determinate token exists")
		}
	}
	ok = true
	return &ok, nil
}

// Truncate cleans all data.
func (sp *ScyllaTokenProvider) Truncate() derrors.Error{
	if err := sp.CheckConnectionAndConnect(); err != nil {
		return err
	}

	err := sp.Session.Query("TRUNCATE TABLE tokens").Exec()
	if err != nil {
		log.Info().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the table")
		return derrors.AsError(err, "cannot truncate token table")
	}

	return nil
}