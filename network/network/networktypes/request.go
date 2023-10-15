package networktypes

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/pkg/xtime"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// Request action descriptions
const (
	RequestActionAddAccount        = "add account to the network"
	RequestActionAddValidator      = "join the network as a validator"
	RequestActionAddVestingAccount = "add vesting account to the network"
	RequestActionRemoveAccount     = "remove account from the network"
	RequestActionRemoveValidator   = "remove validator from the network"
	RequestActionChangeParams      = "change param on the network"

	RequestActionResultAddAccount        = "account added to the network"
	RequestActionResultAddValidator      = "Validator added to the network"
	RequestActionResultAddVestingAccount = "vesting account added to the network"
	RequestActionResultRemoveAccount     = "account removed from network"
	RequestActionResultRemoveValidator   = "validator removed from network"
	RequestActionResultChangeParams      = "param changed on network"

	RequestActionUnrecognized = "<unrecognized request>"
)

type (
	// Request represents the launch Request of a chain on SPN.
	Request struct {
		LaunchID  uint64                     `json:"LaunchID"`
		RequestID uint64                     `json:"RequestID"`
		Creator   string                     `json:"Creator"`
		CreatedAt string                     `json:"CreatedAt"`
		Content   launchtypes.RequestContent `json:"Content"`
		Status    string                     `json:"Status"`
	}
)

// ToRequest converts a request data from SPN and returns a Request object.
func ToRequest(request launchtypes.Request) Request {
	return Request{
		LaunchID:  request.LaunchID,
		RequestID: request.RequestID,
		Creator:   request.Creator,
		CreatedAt: xtime.FormatUnixInt(request.CreatedAt),
		Content:   request.Content,
		Status:    launchtypes.Request_Status_name[int32(request.Status)],
	}
}

// RequestsFromRequestContents creates a list of requests from a list request contents to simulate requests that have not been sent to request pool yet
// The request ID is set to 0 for the first request and incremented for each request, other values are not set
func RequestsFromRequestContents(launchID uint64, contents []launchtypes.RequestContent) []Request {
	requests := make([]Request, len(contents))
	for i, content := range contents {
		requests[i] = Request{
			LaunchID:  launchID,
			RequestID: uint64(i),
			Content:   content,
		}
	}
	return requests
}

// RequestActionDescriptionFromContent describes the action of the request from its content
func RequestActionDescriptionFromContent(content launchtypes.RequestContent) string {
	switch content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		return RequestActionAddAccount
	case *launchtypes.RequestContent_GenesisValidator:
		return RequestActionAddValidator
	case *launchtypes.RequestContent_VestingAccount:
		return RequestActionAddVestingAccount
	case *launchtypes.RequestContent_AccountRemoval:
		return RequestActionRemoveAccount
	case *launchtypes.RequestContent_ValidatorRemoval:
		return RequestActionRemoveValidator
	case *launchtypes.RequestContent_ParamChange:
		return RequestActionChangeParams
	default:
		return RequestActionUnrecognized
	}
}

// RequestActionResultDescriptionFromContent describe the result of the action of the request from its content
func RequestActionResultDescriptionFromContent(content launchtypes.RequestContent) string {
	switch content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		return RequestActionResultAddAccount
	case *launchtypes.RequestContent_GenesisValidator:
		return RequestActionResultAddValidator
	case *launchtypes.RequestContent_VestingAccount:
		return RequestActionResultAddVestingAccount
	case *launchtypes.RequestContent_AccountRemoval:
		return RequestActionResultRemoveAccount
	case *launchtypes.RequestContent_ValidatorRemoval:
		return RequestActionResultRemoveValidator
	case *launchtypes.RequestContent_ParamChange:
		return RequestActionResultChangeParams
	default:
		return RequestActionUnrecognized
	}
}

// VerifyRequest verifies the validity of the request from its content (static check).
func VerifyRequest(request Request) error {
	req, ok := request.Content.Content.(*launchtypes.RequestContent_GenesisValidator)
	if ok {
		err := VerifyAddValidatorRequest(req)
		if err != nil {
			return NewWrappedErrInvalidRequest(request.RequestID, err.Error())
		}
	}

	return nil
}

// VerifyAddValidatorRequest verifies the validator request parameters.
func VerifyAddValidatorRequest(req *launchtypes.RequestContent_GenesisValidator) error {
	// If this is an add validator request
	var (
		peer           = req.GenesisValidator.Peer
		valAddress     = req.GenesisValidator.Address
		consPubKey     = req.GenesisValidator.ConsPubKey
		selfDelegation = req.GenesisValidator.SelfDelegation
	)

	// Check values inside the gentx are correct
	info, err := cosmosutil.ParseGentx(req.GenesisValidator.GenTx)
	if err != nil {
		return fmt.Errorf("cannot parse gentx %w", err)
	}

	// Change the address prefix fetched from the gentx to the one used on SPN
	// Because all on-chain stored address on SPN uses the SPN prefix
	spnFetchedAddress, err := cosmosutil.ChangeAddressPrefix(info.DelegatorAddress, SPN)
	if err != nil {
		return err
	}

	// Check validator address
	if valAddress != spnFetchedAddress {
		return fmt.Errorf(
			"the validator address %s doesn't match the one inside the gentx %s",
			valAddress,
			spnFetchedAddress,
		)
	}

	// Check validator address
	if !info.PubKey.Equals(ed25519.PubKey(consPubKey)) {
		return fmt.Errorf(
			"the consensus pub key %s doesn't match the one inside the gentx %s",
			ed25519.PubKey(consPubKey).String(),
			info.PubKey.String(),
		)
	}

	// Check self delegation
	if selfDelegation.Denom != info.SelfDelegation.Denom ||
		!selfDelegation.IsEqual(info.SelfDelegation) {
		return fmt.Errorf(
			"the self delegation %s doesn't match the one inside the gentx %s",
			selfDelegation.String(),
			info.SelfDelegation.String(),
		)
	}

	// Check the format of the peer
	if !VerifyPeerFormat(peer) {
		return fmt.Errorf(
			"the peer address %s doesn't match the peer format <host>:<port>",
			peer.String(),
		)
	}
	return nil
}
