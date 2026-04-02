package beosin

import (
	"context"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// ProviderName is the identifier for the Beosin provider.
const ProviderName = "beosin"

// Provider implements kyt.Provider for Beosin.
type Provider struct {
	client   beosin.Client
	mapper   *Mapper
	mapperV4 *MapperV4
	useV4    bool
}

// Option configures a Provider.
type Option func(*Provider)

// WithV4 enables the V4 API for more detailed risk information.
func WithV4() Option {
	return func(p *Provider) {
		p.useV4 = true
	}
}

// New creates a new Beosin KYT provider with the given client.
func New(client beosin.Client, opts ...Option) *Provider {
	p := &Provider{
		client:   client,
		mapper:   NewMapper(),
		mapperV4: NewMapperV4(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return ProviderName
}

// AddressRisk performs risk assessment on an address.
func (p *Provider) AddressRisk(ctx context.Context, req *kyt.AddressRiskRequest) (*kyt.RiskResult, error) {
	if req == nil {
		return nil, kyt.NewValidationError("request is nil", nil)
	}
	if req.Address == "" {
		return nil, kyt.NewValidationError("address is required", nil)
	}
	if req.ChainID == "" {
		return nil, kyt.NewValidationError("chain_id is required", nil)
	}

	if p.useV4 {
		return p.addressRiskV4(ctx, req)
	}
	return p.addressRiskV3(ctx, req)
}

// addressRiskV3 uses the V3 API for address risk assessment.
func (p *Provider) addressRiskV3(ctx context.Context, req *kyt.AddressRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.AddressRiskRequest{
		ChainID: req.ChainID,
		Address: req.Address,
	}

	resp, err := p.client.EOAAddressRiskAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapper.MapAddressRisk(resp), nil
}

// addressRiskV4 uses the V4 API for address risk assessment.
func (p *Provider) addressRiskV4(ctx context.Context, req *kyt.AddressRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.AddressRiskRequest{
		ChainID: req.ChainID,
		Address: req.Address,
	}

	resp, err := p.client.V4EOAAddressRiskAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapperV4.MapAddressRisk(resp), nil
}

// DepositRisk performs risk assessment on a deposit transaction.
func (p *Provider) DepositRisk(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	if req == nil {
		return nil, kyt.NewValidationError("request is nil", nil)
	}
	if req.TxHash == "" {
		return nil, kyt.NewValidationError("tx_hash is required", nil)
	}
	if req.ChainID == "" {
		return nil, kyt.NewValidationError("chain_id is required", nil)
	}

	if p.useV4 {
		return p.depositRiskV4(ctx, req)
	}
	return p.depositRiskV2(ctx, req)
}

// depositRiskV2 uses the V2 API for deposit risk assessment.
func (p *Provider) depositRiskV2(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.DepositRequest{
		ChainID: req.ChainID,
		Hash:    req.TxHash,
	}
	if req.Token != nil {
		beosinReq.Token = *req.Token
	}

	resp, err := p.client.DepositTransactionAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapper.MapTransactionRisk(resp), nil
}

// depositRiskV4 uses the V4 API for deposit risk assessment.
func (p *Provider) depositRiskV4(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.DepositRequest{
		ChainID: req.ChainID,
		Hash:    req.TxHash,
	}
	if req.Token != nil {
		beosinReq.Token = *req.Token
	}

	resp, err := p.client.V4DepositTransactionAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapperV4.MapTransactionRisk(resp), nil
}

// WithdrawRisk performs risk assessment on a withdrawal transaction.
func (p *Provider) WithdrawRisk(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	if req == nil {
		return nil, kyt.NewValidationError("request is nil", nil)
	}
	if req.TxHash == "" {
		return nil, kyt.NewValidationError("tx_hash is required", nil)
	}
	if req.ChainID == "" {
		return nil, kyt.NewValidationError("chain_id is required", nil)
	}

	if p.useV4 {
		return p.withdrawRiskV4(ctx, req)
	}
	return p.withdrawRiskV2(ctx, req)
}

// withdrawRiskV2 uses the V2 API for withdrawal risk assessment.
func (p *Provider) withdrawRiskV2(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.WithdrawalRequest{
		ChainID: req.ChainID,
		Hash:    req.TxHash,
	}
	if req.Token != nil {
		beosinReq.Token = *req.Token
	}

	resp, err := p.client.WithdrawalTransactionAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapper.MapTransactionRisk(resp), nil
}

// withdrawRiskV4 uses the V4 API for withdrawal risk assessment.
func (p *Provider) withdrawRiskV4(ctx context.Context, req *kyt.TransactionRiskRequest) (*kyt.RiskResult, error) {
	beosinReq := beosin.WithdrawalRequest{
		ChainID: req.ChainID,
		Hash:    req.TxHash,
	}
	if req.Token != nil {
		beosinReq.Token = *req.Token
	}

	resp, err := p.client.V4WithdrawalTransactionAssessment(ctx, &beosinReq)
	if err != nil {
		return nil, p.wrapError(err)
	}

	if !resp.IsSuccess() {
		return nil, kyt.NewProviderError(ProviderName, resp.Msg, nil)
	}

	return p.mapperV4.MapTransactionRisk(resp), nil
}

// Close releases resources held by the provider.
func (p *Provider) Close() error {
	return nil
}

// wrapError converts Beosin errors to unified KYT errors.
func (p *Provider) wrapError(err error) error {
	if apiErr, ok := err.(*beosin.APIError); ok {
		switch {
		case apiErr.IsTaskExecuting():
			return kyt.NewRetryableError("assessment in progress", err)
		case apiErr.IsPlatformNotSupported():
			return kyt.NewValidationError("unsupported chain", err)
		case apiErr.IsAddressError():
			return kyt.NewValidationError("invalid address", err)
		case apiErr.IsTxHashError(), apiErr.IsTxHashNotExist():
			return kyt.NewValidationError("invalid or non-existent transaction hash", err)
		default:
			return kyt.NewProviderError(ProviderName, apiErr.Message, err)
		}
	}
	return kyt.NewProviderError(ProviderName, err.Error(), err)
}

// Info returns information about the Beosin provider.
func (p *Provider) Info() kyt.ProviderInfo {
	return kyt.ProviderInfo{
		Name:            ProviderName,
		DisplayName:     "Beosin KYT",
		Description:     "Beosin blockchain compliance and security platform",
		SupportedChains: kyt.FullSupportChainIDs(),
	}
}
