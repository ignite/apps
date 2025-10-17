package app

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	gnovmante "github.com/ignite/gnovm/x/gnovm/ante"
)

// setAnteHandler sets the ante handler for the application.
func (app *App) setAnteHandler(options ante.HandlerOptions) error {
	if options.AccountKeeper == nil {
		return errors.New("account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return errors.New("bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return errors.New("sign mode handler is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		gnovmante.NewAnteHandler(),
	}

	app.SetAnteHandler(sdk.ChainAnteDecorators(anteDecorators...))
	return nil
}
