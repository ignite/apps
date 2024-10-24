package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	launchtypes "github.com/ignite/network/x/launch/types"

	"github.com/ignite/apps/network/network/networktypes"
)

// Reviewal keeps a request's reviewal.
type Reviewal struct {
	RequestID  uint64
	IsApproved bool
}

// ApproveRequest returns approval for a request with id.
func ApproveRequest(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: true,
	}
}

// RejectRequest returns rejection for a request with id.
func RejectRequest(requestID uint64) Reviewal {
	return Reviewal{
		RequestID:  requestID,
		IsApproved: false,
	}
}

// Requests fetches all the chain requests from SPN by launch id.
func (n Network) Requests(ctx context.Context, launchID uint64) ([]networktypes.Request, error) {
	res, err := n.launchQuery.ListRequest(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchId: launchID,
	})
	if err != nil {
		return nil, err
	}
	requests := make([]networktypes.Request, len(res.Request))
	for i, req := range res.Request {
		requests[i] = networktypes.ToRequest(req)
	}
	return requests, nil
}

// Request fetches the chain request from SPN by launch and request id.
func (n Network) Request(ctx context.Context, launchID, requestID uint64) (networktypes.Request, error) {
	res, err := n.launchQuery.GetRequest(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchId:  launchID,
		RequestId: requestID,
	})
	if err != nil {
		return networktypes.Request{}, err
	}
	return networktypes.ToRequest(res.Request), nil
}

// RequestFromIDs fetches the chain requested from SPN by launch and provided request IDs
// TODO: once implemented, use the SPN query from https://github.com/ignite/network/issues/420
func (n Network) RequestFromIDs(ctx context.Context, launchID uint64, requestIDs ...uint64) (reqs []networktypes.Request, err error) {
	for _, id := range requestIDs {
		req, err := n.Request(ctx, launchID, id)
		if err != nil {
			return reqs, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

// SubmitRequestReviewals submits reviewals for proposals in batch for chain.
func (n Network) SubmitRequestReviewals(ctx context.Context, launchID uint64, reviewal ...Reviewal) error {
	n.ev.Send("Submitting requests...", events.ProgressStart())

	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	messages := make([]sdk.Msg, len(reviewal))
	for i, reviewal := range reviewal {
		messages[i] = launchtypes.NewMsgSettleRequest(
			addr,
			launchID,
			reviewal.RequestID,
			reviewal.IsApproved,
		)
	}

	res, err := n.cosmos.BroadcastTx(ctx, n.account, messages...)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSettleRequestResponse
	return res.Decode(&requestRes)
}

// SendRequest creates and sends the Request message to SPN.
func (n Network) SendRequest(
	ctx context.Context,
	launchID uint64,
	content launchtypes.RequestContent,
) error {
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msg := launchtypes.NewMsgSendRequest(
		addr,
		launchID,
		content,
	)

	n.ev.Send("Broadcasting transaction", events.ProgressStart())

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgSendRequestResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	if requestRes.AutoApproved {
		n.ev.Send(fmt.Sprintf(
			"%s by the coordinator!", networktypes.RequestActionResultDescriptionFromContent(content),
		),
			events.ProgressFinish())
	} else {
		n.ev.Send(
			fmt.Sprintf(
				"Request %d to %s has been submitted!",
				requestRes.RequestId,
				networktypes.RequestActionDescriptionFromContent(content),
			),
			events.ProgressFinish(),
		)
	}
	return nil
}

// SendRequests creates and sends the Request message to SPN.
func (n Network) SendRequests(
	ctx context.Context,
	launchID uint64,
	contents []launchtypes.RequestContent,
) error {
	for _, content := range contents {
		if err := n.SendRequest(ctx, launchID, content); err != nil {
			return err
		}
	}
	return nil
}
