package client

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

type Domain struct {
	Name           string   `json:"name"`
	Comment        string   `json:"comment,omitempty"`
	Managers       []string `json:"managers,omitempty"`
	MaxUsers       *int64   `json:"max_users,omitempty"`
	MaxAliases     *int64   `json:"max_aliases,omitempty"`
	MaxQuotaBytes  *int64   `json:"max_quota_bytes,omitempty"`
	SignupEnabled  *bool    `json:"signup_enabled,omitempty"`
	Alternatives   []string `json:"alternatives,omitempty"`
	DNSAutoconfig  []string `json:"dns_autoconfig,omitempty"`
	DNSMX          string   `json:"dns_mx,omitempty"`
	DNSSPF         string   `json:"dns_spf,omitempty"`
	DNSDKIM        string   `json:"dns_dkim,omitempty"`
	DNSDMARC       string   `json:"dns_dmarc,omitempty"`
	DNSDMARCReport string   `json:"dns_dmarc_report,omitempty"`
	DNSTLSA        []string `json:"dns_tlsa,omitempty"`
}

type DomainUpdate struct {
	Comment       string   `json:"comment,omitempty"`
	MaxUsers      *int64   `json:"max_users,omitempty"`
	MaxAliases    *int64   `json:"max_aliases,omitempty"`
	MaxQuotaBytes *int64   `json:"max_quota_bytes,omitempty"`
	SignupEnabled *bool    `json:"signup_enabled,omitempty"`
	Alternatives  []string `json:"alternatives,omitempty"`
}

type User struct {
	Email              string   `json:"email"`
	RawPassword        string   `json:"raw_password,omitempty"`
	Password           string   `json:"password,omitempty"`
	Comment            string   `json:"comment,omitempty"`
	QuotaBytes         *int64   `json:"quota_bytes,omitempty"`
	QuotaBytesUsed     *int64   `json:"quota_bytes_used,omitempty"`
	GlobalAdmin        *bool    `json:"global_admin,omitempty"`
	Enabled            *bool    `json:"enabled,omitempty"`
	ChangePWNextLogin  *bool    `json:"change_pw_next_login,omitempty"`
	EnableIMAP         *bool    `json:"enable_imap,omitempty"`
	EnablePOP          *bool    `json:"enable_pop,omitempty"`
	AllowSpoofing      *bool    `json:"allow_spoofing,omitempty"`
	ForwardEnabled     *bool    `json:"forward_enabled,omitempty"`
	ForwardDestination []string `json:"forward_destination,omitempty"`
	ForwardKeep        *bool    `json:"forward_keep,omitempty"`
	ReplyEnabled       *bool    `json:"reply_enabled,omitempty"`
	ReplySubject       string   `json:"reply_subject,omitempty"`
	ReplyBody          string   `json:"reply_body,omitempty"`
	ReplyStartDate     string   `json:"reply_startdate,omitempty"`
	ReplyEndDate       string   `json:"reply_enddate,omitempty"`
	DisplayedName      string   `json:"displayed_name,omitempty"`
	SpamEnabled        *bool    `json:"spam_enabled,omitempty"`
	SpamMarkAsRead     *bool    `json:"spam_mark_as_read,omitempty"`
	SpamThreshold      *int64   `json:"spam_threshold,omitempty"`
}

type Alias struct {
	Email       string   `json:"email"`
	Destination []string `json:"destination,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Wildcard    *bool    `json:"wildcard,omitempty"`
}

type AlternativeDomain struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type ManagerCreate struct {
	UserEmail string `json:"user_email"`
}

type Relay struct {
	Name    string `json:"name"`
	SMTP    string `json:"smtp,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type RelayUpdate struct {
	SMTP    string `json:"smtp,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type Token struct {
	ID            FlexibleString `json:"id,omitempty"`
	Token         string         `json:"token,omitempty"`
	Email         string         `json:"email,omitempty"`
	Comment       string         `json:"comment,omitempty"`
	AuthorizedIPs []string       `json:"AuthorizedIP,omitempty"`
	Created       string         `json:"Created,omitempty"`
	LastEdit      string         `json:"Last edit,omitempty"`
}

type TokenUpdate struct {
	Comment       string   `json:"comment,omitempty"`
	AuthorizedIPs []string `json:"AuthorizedIP,omitempty"`
}

type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err == nil {
		*s = FlexibleString(raw)
		return nil
	}

	var number float64
	if err := json.Unmarshal(data, &number); err == nil {
		*s = FlexibleString(strconv.FormatFloat(number, 'f', -1, 64))
		return nil
	}

	return nil
}

func (s FlexibleString) String() string {
	return string(s)
}

func (c *Client) CreateDomain(ctx context.Context, domain Domain) error {
	return c.Post(ctx, "/domain", domain, nil)
}

func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var domains []Domain
	if err := c.Get(ctx, "/domain", &domains); err != nil {
		return nil, err
	}

	return domains, nil
}

func (c *Client) GetDomain(ctx context.Context, name string) (*Domain, error) {
	var domain Domain
	err := c.Get(ctx, "/domain/"+url.PathEscape(name), &domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (c *Client) UpdateDomain(ctx context.Context, name string, domain DomainUpdate) error {
	return c.Patch(ctx, "/domain/"+url.PathEscape(name), domain, nil)
}

func (c *Client) DeleteDomain(ctx context.Context, name string) error {
	return c.Delete(ctx, "/domain/"+url.PathEscape(name))
}

func (c *Client) CreateUser(ctx context.Context, user User) error {
	return c.Post(ctx, "/user", user, nil)
}

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := c.Get(ctx, "/user", &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Client) GetUser(ctx context.Context, email string) (*User, error) {
	var user User
	err := c.Get(ctx, "/user/"+url.PathEscape(email), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) UpdateUser(ctx context.Context, email string, user User) error {
	return c.Patch(ctx, "/user/"+url.PathEscape(email), user, nil)
}

func (c *Client) DeleteUser(ctx context.Context, email string) error {
	return c.Delete(ctx, "/user/"+url.PathEscape(email))
}

func (c *Client) CreateAlias(ctx context.Context, alias Alias) error {
	return c.Post(ctx, "/alias", alias, nil)
}

func (c *Client) ListAliases(ctx context.Context) ([]Alias, error) {
	var aliases []Alias
	if err := c.Get(ctx, "/alias", &aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (c *Client) GetAlias(ctx context.Context, email string) (*Alias, error) {
	var alias Alias
	err := c.Get(ctx, "/alias/"+url.PathEscape(email), &alias)
	if err != nil {
		return nil, err
	}

	return &alias, nil
}

func (c *Client) UpdateAlias(ctx context.Context, email string, alias Alias) error {
	return c.Patch(ctx, "/alias/"+url.PathEscape(email), alias, nil)
}

func (c *Client) DeleteAlias(ctx context.Context, email string) error {
	return c.Delete(ctx, "/alias/"+url.PathEscape(email))
}

func (c *Client) CreateAlternativeDomain(ctx context.Context, alternative AlternativeDomain) error {
	return c.Post(ctx, "/alternative", alternative, nil)
}

func (c *Client) GetAlternativeDomain(ctx context.Context, name string) (*AlternativeDomain, error) {
	var alternative AlternativeDomain
	err := c.Get(ctx, "/alternative/"+url.PathEscape(name), &alternative)
	if err != nil {
		return nil, err
	}

	return &alternative, nil
}

func (c *Client) DeleteAlternativeDomain(ctx context.Context, name string) error {
	return c.Delete(ctx, "/alternative/"+url.PathEscape(name))
}

func (c *Client) CreateDomainManager(ctx context.Context, domain string, manager ManagerCreate) error {
	return c.Post(ctx, "/domain/"+url.PathEscape(domain)+"/manager", manager, nil)
}

func (c *Client) GetDomainManager(ctx context.Context, domain string, email string) error {
	var out map[string]any
	return c.Get(ctx, "/domain/"+url.PathEscape(domain)+"/manager/"+url.PathEscape(email), &out)
}

func (c *Client) DeleteDomainManager(ctx context.Context, domain string, email string) error {
	return c.Delete(ctx, "/domain/"+url.PathEscape(domain)+"/manager/"+url.PathEscape(email))
}

func (c *Client) CreateRelay(ctx context.Context, relay Relay) error {
	return c.Post(ctx, "/relay", relay, nil)
}

func (c *Client) GetRelay(ctx context.Context, name string) (*Relay, error) {
	var relay Relay
	err := c.Get(ctx, "/relay/"+url.PathEscape(name), &relay)
	if err != nil {
		return nil, err
	}

	return &relay, nil
}

func (c *Client) UpdateRelay(ctx context.Context, name string, relay RelayUpdate) error {
	return c.Patch(ctx, "/relay/"+url.PathEscape(name), relay, nil)
}

func (c *Client) DeleteRelay(ctx context.Context, name string) error {
	return c.Delete(ctx, "/relay/"+url.PathEscape(name))
}

func (c *Client) CreateToken(ctx context.Context, token Token) (*Token, error) {
	var created Token
	err := c.Post(ctx, "/token", token, &created)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (c *Client) GetToken(ctx context.Context, id string) (*Token, error) {
	var token Token
	err := c.Get(ctx, "/token/"+url.PathEscape(id), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Client) UpdateToken(ctx context.Context, id string, token TokenUpdate) error {
	return c.Patch(ctx, "/token/"+url.PathEscape(id), token, nil)
}

func (c *Client) DeleteToken(ctx context.Context, id string) error {
	return c.Delete(ctx, "/token/"+url.PathEscape(id))
}

func (c *Client) GenerateDKIM(ctx context.Context, domain string) error {
	return c.Post(ctx, "/domain/"+url.PathEscape(domain)+"/dkim", nil, nil)
}
