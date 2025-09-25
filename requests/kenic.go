package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Response structures for frontend
type DomainCheckResponse struct {
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

type DomainInfoResponse struct {
	Domain       string        `json:"domain"`
	ROID         string        `json:"roid"`
	Status       []string      `json:"status"`
	Registrant   string        `json:"registrant"`
	Contacts     []ContactInfo `json:"contacts"`
	NameServers  []string      `json:"nameservers"`
	ClientID     string        `json:"client_id"`
	CreatorID    string        `json:"creator_id"`
	CreatedDate  string        `json:"created_date"`
	UpdaterID    string        `json:"updater_id,omitempty"`
	UpdatedDate  string        `json:"updated_date,omitempty"`
	ExpiryDate   string        `json:"expiry_date"`
	TransferDate string        `json:"transfer_date,omitempty"`
}

type ContactInfo struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type DomainSuggestion struct {
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Type      string `json:"type"` // "alternative", "similar", "tld_variant"
}

// XML parsing structures (simplified)
type EPPResponse struct {
	Result struct {
		Code int    `xml:"code,attr"`
		Msg  string `xml:"msg"`
	} `xml:"response>result"`

	ResData struct {
		CheckData struct {
			Names []struct {
				Name struct {
					Avail string `xml:"avail,attr"`
					Value string `xml:",chardata"`
				} `xml:"name"`
				Reason string `xml:"reason,omitempty"`
			} `xml:"cd"`
		} `xml:"chkData"`

		InfoData struct {
			Name   string `xml:"name"`
			ROID   string `xml:"roid"`
			Status []struct {
				S string `xml:"s,attr"`
			} `xml:"status"`
			Registrant string `xml:"registrant"`
			Contact    []struct {
				Type string `xml:"type,attr"`
				ID   string `xml:",chardata"`
			} `xml:"contact"`
			NS struct {
				HostObj []string `xml:"hostObj"`
			} `xml:"ns"`
			ClID   string `xml:"clID"`
			CrID   string `xml:"crID"`
			CrDate string `xml:"crDate"`
			UpID   string `xml:"upID"`
			UpDate string `xml:"upDate"`
			ExDate string `xml:"exDate"`
			TrDate string `xml:"trDate"`
		} `xml:"infData"`
	} `xml:"response>resData"`
	Extension struct {
		FeeCheckData struct {
			Currency string `xml:"currency"`
			CD       []struct {
				Avail   string `xml:"avail,attr"`
				ObjID   string `xml:"objID"`
				Class   string `xml:"class"`
				Command struct {
					Name   string `xml:"name,attr"`
					Period struct {
						Unit  string `xml:"unit,attr"`
						Value string `xml:",chardata"`
					} `xml:"period"`
					Fee struct {
						Description string  `xml:"description,attr"`
						Refundable  string  `xml:"refundable,attr"`
						Amount      float64 `xml:",chardata"`
					} `xml:"fee"`
				} `xml:"command"`
			} `xml:"cd"`
		} `xml:"chkData"`
	} `xml:"response>extension"`
}

// KENIC EPP Client
type KenicClient struct {
	conn     net.Conn
	host     string
	username string
	password string
	mutex    sync.Mutex
	loggedIn bool
}

func NewKenicClient(host, username, password string) *KenicClient {
	return &KenicClient{
		host:     host,
		username: username,
		password: password,
	}
}

func (c *KenicClient) eppFrame(xml string) []byte {
	length := uint32(len(xml) + 4)
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, length)
	buf.WriteString(xml)
	return buf.Bytes()
}

func (c *KenicClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil && c.loggedIn {
		return nil // Already connected
	}

	conn, err := tls.Dial("tcp", c.host, &tls.Config{
		InsecureSkipVerify: true, // Use proper certs in production
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn

	// Read greeting
	buf := make([]byte, 8192)
	_, err = c.conn.Read(buf)
	if err != nil {
		c.conn.Close()
		c.conn = nil
		return fmt.Errorf("failed to read greeting: %v", err)
	}

	return c.login()
}

func (c *KenicClient) login() error {
	login := fmt.Sprintf(`
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <login>
      <clID>%s</clID>
      <pw>%s</pw>
      <options>
        <version>1.0</version>
        <lang>en</lang>
      </options>
      <services>
        <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
        <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
        <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
      </services>
    </login>
    <clTRID>login-%d</clTRID>
  </command>
</epp>`, c.username, c.password, time.Now().UnixNano())

	_, err := c.conn.Write(c.eppFrame(login))
	if err != nil {
		return fmt.Errorf("failed to send login: %v", err)
	}

	buf := make([]byte, 8192)
	n, err := c.conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read login response: %v", err)
	}

	var resp EPPResponse
	if err := xml.Unmarshal(buf[4:n], &resp); err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}
	log.Println("Check Raw Response: ", string(buf[4:n]))
	log.Println("Login Response: ", resp)

	if resp.Result.Code != 1000 {
		return fmt.Errorf("login failed: %s", resp.Result.Msg)
	}

	c.loggedIn = true
	return nil
}

func (c *KenicClient) CheckDomains(domains []string) ([]DomainCheckResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.loggedIn {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	// Build domain list for XML
	var domainList strings.Builder
	for _, domain := range domains {
		domainList.WriteString(fmt.Sprintf("        <domain:name>%s</domain:name>\n", domain))
	}

	check := fmt.Sprintf(`
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <check>
      <domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
%s      </domain:check>
    </check>
	<extension>
		<fee:check xmlns:fee="urn:ietf:params:xml:ns:epp:fee-1.0">
			<fee:currency>KES</fee:currency>
			<fee:command name="create">
				<fee:period unit="y">1</fee:period>
			</fee:command>
		</fee:check>
	</extension>
    <clTRID>CHECK-%d</clTRID>
  </command>
</epp>`, domainList.String(), time.Now().UnixNano())

	_, err := c.conn.Write(c.eppFrame(check))
	if err != nil {
		return nil, fmt.Errorf("failed to send check command: %v", err)
	}

	buf := make([]byte, 8192)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read check response: %v", err)
	}

	var resp EPPResponse
	if err := xml.Unmarshal(buf[4:n], &resp); err != nil {
		return nil, fmt.Errorf("failed to parse check response: %v", err)
	}
	log.Println("Check String: ", check)
	log.Println("Check Raw Response: ", string(buf[4:n]))
	log.Println("Check Response: ", resp)

	if resp.Result.Code != 1000 {
		return nil, fmt.Errorf("check failed: %s", resp.Result.Msg)
	}

	var results []DomainCheckResponse
	for _, cd := range resp.ResData.CheckData.Names {
		results = append(results, DomainCheckResponse{
			Domain:    cd.Name.Value,
			Available: cd.Name.Avail == "1",
			Reason:    cd.Reason,
		})
	}

	return results, nil
}

func (c *KenicClient) GetDomainInfo(domain string) (*DomainInfoResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.loggedIn {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	info := fmt.Sprintf(`
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <info>
      <domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name hosts="all">%s</domain:name>
      </domain:info>
    </info>
    <clTRID>INFO-%d</clTRID>
  </command>
</epp>`, domain, time.Now().UnixNano())

	_, err := c.conn.Write(c.eppFrame(info))
	if err != nil {
		return nil, fmt.Errorf("failed to send info command: %v", err)
	}

	buf := make([]byte, 8192)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read info response: %v", err)
	}

	var resp EPPResponse
	if err := xml.Unmarshal(buf[4:n], &resp); err != nil {
		return nil, fmt.Errorf("failed to parse info response: %v", err)
	}
	log.Println("Check Raw Response: ", string(buf[4:n]))
	log.Println("Info Response: ", resp)

	if resp.Result.Code != 1000 {
		return nil, fmt.Errorf("info failed: %s", resp.Result.Msg)
	}

	data := resp.ResData.InfoData

	// Parse status
	var statuses []string
	for _, status := range data.Status {
		statuses = append(statuses, status.S)
	}

	// Parse contacts
	var contacts []ContactInfo
	for _, contact := range data.Contact {
		contacts = append(contacts, ContactInfo{
			Type: contact.Type,
			ID:   contact.ID,
		})
	}

	return &DomainInfoResponse{
		Domain:       data.Name,
		ROID:         data.ROID,
		Status:       statuses,
		Registrant:   data.Registrant,
		Contacts:     contacts,
		NameServers:  data.NS.HostObj,
		ClientID:     data.ClID,
		CreatorID:    data.CrID,
		CreatedDate:  data.CrDate,
		UpdaterID:    data.UpID,
		UpdatedDate:  data.UpDate,
		ExpiryDate:   data.ExDate,
		TransferDate: data.TrDate,
	}, nil
}

func (c *KenicClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil && c.loggedIn {
		logout := fmt.Sprintf(`
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <logout/>
    <clTRID>logout-%d</clTRID>
  </command>
</epp>`, time.Now().UnixNano())

		c.conn.Write(c.eppFrame(logout))

		buf := make([]byte, 8192)
		c.conn.Read(buf) // Read logout response
		log.Println("Logout Response: ", string(buf))

		c.conn.Close()
		c.conn = nil
		c.loggedIn = false
	}
	return nil
}

// Domain suggestion generator
func GenerateDomainSuggestions(domain string, takenDomains []string) []DomainSuggestion {
	var suggestions []DomainSuggestion
	baseDomain := strings.TrimSuffix(domain, ".ke")

	// Common prefixes and suffixes
	prefixes := []string{"get", "my", "the", "best", "top", "new", "kenya"}
	suffixes := []string{"kenya", "ke", "online", "shop", "hub", "pro", "plus", "digital"}

	// Alternative spellings
	variations := []string{
		baseDomain + "s",
		baseDomain + "kenya",
		"ke" + baseDomain,
		"my" + baseDomain,
	}

	// Add prefix suggestions
	for _, prefix := range prefixes {
		suggestion := prefix + baseDomain + ".ke"
		if !contains(takenDomains, suggestion) && len(suggestions) < 10 {
			suggestions = append(suggestions, DomainSuggestion{
				Domain:    suggestion,
				Available: true,
				Type:      "alternative",
			})
		}
	}

	// Add suffix suggestions
	for _, suffix := range suffixes {
		suggestion := baseDomain + suffix + ".ke"
		if !contains(takenDomains, suggestion) && len(suggestions) < 10 {
			suggestions = append(suggestions, DomainSuggestion{
				Domain:    suggestion,
				Available: true,
				Type:      "alternative",
			})
		}
	}

	// Add variations
	for _, variation := range variations {
		if variation != baseDomain {
			suggestion := variation + ".ke"
			if !contains(takenDomains, suggestion) && len(suggestions) < 10 {
				suggestions = append(suggestions, DomainSuggestion{
					Domain:    suggestion,
					Available: true,
					Type:      "similar",
				})
			}
		}
	}

	return suggestions
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// WHOIS lookup (simplified)
func LookupWhois(domain string) (*WhoisInfo, error) {
	conn, err := net.DialTimeout("tcp", "whois.kenic.or.ke:43", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WHOIS server: %v", err)
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\r\n", domain)
	if err != nil {
		return nil, fmt.Errorf("failed to send WHOIS query: %v", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read WHOIS response: %v", err)
	}
	log.Println("WHOIS Raw Response: ", string(response))

	return parseWhoisResponse(domain, string(response))
}

type WhoisInfo struct {
	Domain                              string    `json:"domain"`
	RegistryDomainID                    string    `json:"registry_domain_id,omitempty"`
	Registrar                           string    `json:"registrar,omitempty"`
	RegistrarStreetAddress              string    `json:"registrar_street_address,omitempty"`
	RegistrarPhone                      string    `json:"registrar_phone,omitempty"`
	RegistrarEmail                      string    `json:"registrar_email,omitempty"`
	RegistrarAbuseContactEmail          string    `json:"registrar_abuse_contact_email,omitempty"`
	RegistrarAbuseContactPhone          string    `json:"registrar_abuse_contact_phone,omitempty"`
	CreatedDate                         time.Time `json:"created_date,omitempty"`
	UpdatedDate                         time.Time `json:"updated_date,omitempty"`
	RegistryExpiryDate                  time.Time `json:"registry_expiry_date,omitempty"`
	RegistrarRegistrationExpirationDate time.Time `json:"registrar_registration_expiration_date,omitempty"`
	DomainStatus                        []string  `json:"domain_status,omitempty"`
	NameServers                         []string  `json:"name_servers,omitempty"`
	DNSSEC                              string    `json:"dnssec,omitempty"`
	LastWHOISUpdate                     time.Time `json:"last_whois_update,omitempty"`

	// Optional: keep unparsed/extra lines (handy for debugging)
	Extras map[string]string `json:"extras,omitempty"`
}

// parseWhoisResponse parses a raw WHOIS text block into WhoisInfo.
// It is tolerant to casing/spacing differences, multi-line addresses,
// and common synonyms across WHOIS servers.
func parseWhoisResponse(domain, response string) (*WhoisInfo, error) {
	if strings.TrimSpace(response) == "" {
		return nil, errors.New("empty WHOIS response")
	}

	w := &WhoisInfo{
		Domain: strings.TrimSpace(domain),
		Extras: map[string]string{},
	}

	// Common key normalizer: lowercase, collapse spaces
	normKey := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
		return s
	}

	// Date parser that tries multiple layouts
	parseTime := func(s string) (time.Time, bool) {
		val := strings.TrimSpace(s)
		if val == "" {
			return time.Time{}, false
		}
		layouts := []string{
			time.RFC3339,               // 2006-01-02T15:04:05Z07:00
			"2006-01-02T15:04:05Z",     // 2006-01-02T15:04:05Z
			"2006-01-02 15:04:05Z",     // 2006-01-02 15:04:05Z
			"2006-01-02",               // 2006-01-02
			"2006.01.02 15:04:05",      // some WHOIS styles
			"2006-01-02T15:04:05.000Z", // with milliseconds
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, val); err == nil {
				return t.UTC(), true
			}
		}
		return time.Time{}, false
	}

	// Detect ">>> Last update of WHOIS database: <date> <<<"
	lastUpdateRe := regexp.MustCompile(`(?i)last update of whois database:\s*(.+?)\s*(?:<<<|$)`)

	// Generic "Key: Value" splitter; returns ok=false if no colon.
	splitKV := func(line string) (key, val string, ok bool) {
		i := strings.Index(line, ":")
		if i < 0 {
			return "", "", false
		}
		key = strings.TrimSpace(line[:i])
		val = strings.TrimSpace(line[i+1:])
		return key, val, true
	}

	// Track if we're in a multi-line address continuation
	const addrKey = "registrar street address"
	inAddr := false

	lines := strings.Split(response, "\n")
	for _, raw := range lines {
		line := strings.TrimRightFunc(raw, func(r rune) bool { return r == '\r' || r == '\n' })
		trim := strings.TrimSpace(line)
		if trim == "" {
			inAddr = false // blank line ends any continuation
			continue
		}

		// Last update special case
		if strings.HasPrefix(trim, ">>>") || lastUpdateRe.MatchString(trim) {
			if m := lastUpdateRe.FindStringSubmatch(trim); len(m) == 2 {
				if t, ok := parseTime(m[1]); ok {
					w.LastWHOISUpdate = t
				}
			}
			inAddr = false
			continue
		}

		// Lines that are URLs or help text — ignore
		if strings.Contains(trim, "For more information on domain status codes") ||
			strings.HasPrefix(trim, "https://") || strings.HasPrefix(trim, "http://") {
			inAddr = false
			continue
		}

		// Continuation of address (no colon, but previous key was address)
		if inAddr && !strings.Contains(trim, ":") {
			if w.RegistrarStreetAddress == "" {
				w.RegistrarStreetAddress = trim
			} else {
				// join with a space, preserving commas
				if strings.HasSuffix(w.RegistrarStreetAddress, ",") {
					w.RegistrarStreetAddress += " " + trim
				} else {
					w.RegistrarStreetAddress += " " + trim
				}
			}
			continue
		}

		key, val, ok := splitKV(trim)
		if !ok {
			inAddr = false
			continue
		}
		nk := normKey(key)

		switch nk {
		case "domain name":
			if w.Domain == "" {
				w.Domain = val
			}
		case "registry domain id":
			w.RegistryDomainID = val

		case "registrar":
			// Some WHOIS servers include registrar name here
			w.Registrar = val

		case addrKey:
			inAddr = true
			// Start/reset address with first line’s value
			w.RegistrarStreetAddress = val

		case "registrar phone", "registrar telephone":
			w.RegistrarPhone = val

		case "registrar email", "registrar e-mail":
			w.RegistrarEmail = val

		case "registrar abuse contact email":
			w.RegistrarAbuseContactEmail = val

		case "registrar abuse contact phone":
			w.RegistrarAbuseContactPhone = val

		case "dnssec":
			w.DNSSEC = strings.ToLower(val)

		case "name server", "nameserver":
			ns := strings.ToLower(strings.TrimSpace(val))
			// normalize trailing dot and extra spaces
			ns = strings.TrimSuffix(ns, ".")
			if ns != "" {
				w.NameServers = append(w.NameServers, ns)
			}

		// Dates (handle common synonyms)
		case "creation date", "created", "created on":
			if t, ok := parseTime(val); ok {
				w.CreatedDate = t
			}
		case "updated date", "last updated", "updated on":
			if t, ok := parseTime(val); ok {
				w.UpdatedDate = t
			}
		case "registry expiry date", "registry expiration date", "expiry date", "expires":
			if t, ok := parseTime(val); ok {
				w.RegistryExpiryDate = t
			}
		case "registrar registration expiration date", "registrar expiry date":
			if t, ok := parseTime(val); ok {
				w.RegistrarRegistrationExpirationDate = t
			}

		case "domain status", "status":
			// Extract the status code before any URL
			// e.g., "clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited"
			code := val
			if i := strings.Index(code, " http"); i > 0 {
				code = code[:i]
			}
			code = strings.TrimSpace(code)
			if code != "" {
				w.DomainStatus = append(w.DomainStatus, code)
			}

		default:
			// Keep unknown fields around for debugging
			if _, exists := w.Extras[nk]; !exists {
				w.Extras[nk] = val
			}
			inAddr = false
		}
	}

	// De-duplicate name servers & statuses (order-preserving)
	w.NameServers = uniqPreserve(w.NameServers)
	w.DomainStatus = uniqPreserve(w.DomainStatus)

	// Basic sanity: if the WHOIS has Domain Name, prefer it
	if w.Domain == "" {
		w.Domain = strings.TrimSpace(domain)
	}

	return w, nil
}

func uniqPreserve(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// HTTP API handlers
type DomainServer struct {
	client *KenicClient
}

func NewDomainServer(host, username, password string) *DomainServer {
	return &DomainServer{
		client: NewKenicClient(host, username, password),
	}
}

type SearchRequest struct {
	Domain string `json:"domain"`
}

type SearchResponse struct {
	Domain      string              `json:"domain"`
	Available   bool                `json:"available"`
	Reason      string              `json:"reason,omitempty"`
	CheckedAt   string              `json:"checked_at"`
	TLD         string              `json:"tld"`
	SLD         string              `json:"sld"`
	Premium     bool                `json:"premium"`
	Price       string              `json:"price"`
	WhoisData   *WhoisInfo          `json:"whois_data,omitempty"`
	Suggestions []DomainSuggestion  `json:"suggestions,omitempty"`
	Info        *DomainInfoResponse `json:"info,omitempty"`
}

func (s *DomainServer) handleDomainSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, "Domain required", http.StatusBadRequest)
		return
	}

	// Ensure .ke domain
	// originalDomain := req.Domain
	if !strings.HasSuffix(req.Domain, ".ke") {
		req.Domain += ".ke"
	}

	// Extract SLD (Second Level Domain)
	sld := strings.TrimSuffix(req.Domain, ".ke")

	// Check domain availability
	results, err := s.client.CheckDomains([]string{req.Domain})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check domain: %v", err), http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		http.Error(w, "No results returned", http.StatusInternalServerError)
		return
	}

	result := results[0]

	// Always generate suggestions regardless of availability
	suggestions := GenerateDomainSuggestions(result.Domain, []string{})

	// Create response with enhanced data
	response := SearchResponse{
		Domain:      result.Domain,
		Available:   result.Available,
		Reason:      result.Reason,
		CheckedAt:   time.Now().Format("2006-01-02 15:04:05"),
		TLD:         ".ke",
		SLD:         sld,
		Premium:     isPremiumDomain(sld),
		Price:       getDomainPrice(sld, result.Available),
		Suggestions: suggestions, // Always present
	}

	// If domain is taken, get additional info
	if !result.Available {
		// Try to get domain info from EPP
		if info, err := s.client.GetDomainInfo(result.Domain); err == nil {
			response.Info = info
		}

		// Try to get WHOIS data
		if whoisInfo, err := LookupWhois(result.Domain); err == nil {
			response.WhoisData = whoisInfo
		}
	} else {
		// Even for available domains, we might want to show some info
		response.Info = &DomainInfoResponse{
			Domain: result.Domain,
			Status: []string{"Available for registration"},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(response)
	// json.NewEncoder(w).Encode(results)
}

// Helper functions
func isPremiumDomain(sld string) bool {
	premiumKeywords := []string{
		"sex", "porn", "casino", "bet", "loan", "insurance", "bank", "finance",
		"car", "hotel", "travel", "shop", "store", "business", "company",
	}

	sldLower := strings.ToLower(sld)
	for _, keyword := range premiumKeywords {
		if strings.Contains(sldLower, keyword) {
			return true
		}
	}

	// Short domains (3 chars or less) are often premium
	return len(sld) <= 3
}

func getDomainPrice(sld string, available bool) string {
	if !available {
		return "N/A - Domain taken"
	}

	// Basic pricing logic (you'd get this from KENIC's pricing)
	if isPremiumDomain(sld) {
		return "KES 5,000/year (Premium)"
	}

	if len(sld) <= 3 {
		return "KES 3,000/year (Short domain)"
	}

	return "KES 1,500/year (Standard)"
}

func (s *DomainServer) handleWhoisLookup(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Domain parameter required", http.StatusBadRequest)
		return
	}

	whoisInfo, err := LookupWhois(domain)
	if err != nil {
		http.Error(w, fmt.Sprintf("WHOIS lookup failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(whoisInfo)
}

func (s *DomainServer) Close() {
	s.client.Close()
}

func (c *KenicClient) Ping() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.loggedIn {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	// EPP "hello" is the standard way to keep connection alive
	hello := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <hello/>
</epp>`

	_, err := c.conn.Write(c.eppFrame(hello))
	if err != nil {
		return fmt.Errorf("failed to send hello: %v", err)
	}

	buf := make([]byte, 8192)
	n, err := c.conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read hello response: %v", err)
	}

	log.Println("Hello Raw Response: ", string(buf[4:n]))

	// Usually <hello/> response is just <greeting>, not <response>.
	// So no need to unmarshal into EPPResponse. Instead, you can just check if it has <greeting>.
	if !strings.Contains(string(buf[4:n]), "<greeting") {
		return fmt.Errorf("unexpected hello response: %s", string(buf[4:n]))
	}

	return nil
}

func (s *DomainServer) checkPing(w http.ResponseWriter, r *http.Request) {
	if err := s.client.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Ping failed: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pong"))
}

// Example usage
func maini() {
	// Initialize your domain server
	server := NewDomainServer("ote.kenic.or.ke:700", "hack-a-milli", "TpEjG99Qq69t")
	defer server.Close()

	// Test the connection
	if err := server.client.Connect(); err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	fmt.Println("✅ Connected to KENIC EPP server")

	// Set up routes
	http.HandleFunc("/api/domain/search", server.handleDomainSearch)
	http.HandleFunc("/api/whois", server.handleWhoisLookup)
	http.HandleFunc("/api/ping", server.checkPing)

	fmt.Println("🚀 Server starting on :8080")
	fmt.Println("POST /api/domain/search - Search domains")
	fmt.Println("GET  /api/whois?domain=example.ke - WHOIS lookup")
	fmt.Println("GET  /api/ping - Health check")

	http.ListenAndServe(":8080", nil)
}
