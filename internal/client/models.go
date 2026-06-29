package client

import (
	"context"
	"net/url"
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

func (c *Client) CreateDomain(ctx context.Context, domain Domain) error {
	return c.Post(ctx, "/domain", domain, nil)
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
