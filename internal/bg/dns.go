package bg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

/************* args & small DTOs *************/

type DNSReconcileArgs struct {
	MaxDomains int `json:"max_domains,omitempty"`
	MaxRecords int `json:"max_records,omitempty"`
	IntervalS  int `json:"interval_seconds,omitempty"`
}

// TXT marker content (compact)
type ownershipMarker struct {
	Ver string `json:"v"`   // "ag1"
	Org string `json:"org"` // org UUID
	Rec string `json:"rec"` // record UUID
	Fp  string `json:"fp"`  // short fp (first 16 of sha256)
}

// ExternalDNS poison owner id – MUST NOT match any real external-dns --txt-owner-id
const externalDNSPoisonOwner = "autoglue-lock"

// ExternalDNS poison content – fake owner so real external-dns skips it.
const externalDNSPoisonValue = "heritage=external-dns,external-dns/owner=" + externalDNSPoisonOwner + ",external-dns/resource=manual/autoglue"

/************* entrypoint worker *************/

func DNSReconsileWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := DNSReconcileArgs{MaxDomains: 25, MaxRecords: 100, IntervalS: 30}
		_ = j.ParseArguments(&args)

		if args.MaxDomains <= 0 {
			args.MaxDomains = 25
		}
		if args.MaxRecords <= 0 {
			args.MaxRecords = 100
		}
		if args.IntervalS <= 0 {
			args.IntervalS = 30
		}

		processedDomains, processedRecords, err := reconcileDNSOnce(ctx, db, args)
		if err != nil {
			log.Error().Err(err).Msg("[dns] reconcile tick failed")
		} else {
			log.Debug().
				Int("domains", processedDomains).
				Int("records", processedRecords).
				Msg("[dns] reconcile tick ok")
		}

		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(ctx, uuid.NewString(), "dns_reconcile", args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)

		return map[string]any{
			"domains_processed": processedDomains,
			"records_processed": processedRecords,
		}, nil
	}
}

/************* core tick *************/

func reconcileDNSOnce(ctx context.Context, db *gorm.DB, args DNSReconcileArgs) (int, int, error) {
	var domains []models.Domain

	// 1) validate/backfill pending domains
	if err := db.
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(args.MaxDomains).
		Find(&domains).Error; err != nil {
		return 0, 0, err
	}

	domainsProcessed := 0
	for i := range domains {
		if err := processDomain(ctx, db, &domains[i]); err != nil {
			log.Error().Err(err).Str("domain", domains[i].DomainName).Msg("[dns] domain processing failed")
		} else {
			domainsProcessed++
		}
	}

	// 2) apply pending record sets for ready domains
	var readyDomains []models.Domain
	if err := db.Where("status = ?", "ready").Find(&readyDomains).Error; err != nil {
		return domainsProcessed, 0, err
	}

	recordsProcessed := 0
	for i := range readyDomains {
		n, err := processPendingRecordsForDomain(ctx, db, &readyDomains[i], args.MaxRecords)
		if err != nil {
			log.Error().Err(err).Str("domain", readyDomains[i].DomainName).Msg("[dns] record processing failed")
			continue
		}
		recordsProcessed += n
	}

	return domainsProcessed, recordsProcessed, nil
}

/************* domain processing *************/

func processDomain(ctx context.Context, db *gorm.DB, d *models.Domain) error {
	orgID := d.OrganizationID

	// 1) Load credential (org-guarded)
	var cred models.Credential
	if err := db.Where("id = ? AND organization_id = ?", d.CredentialID, orgID).First(&cred).Error; err != nil {
		return setDomainFailed(db, d, fmt.Errorf("credential not found: %w", err))
	}

	// 2) Decrypt → dto.AWSCredential
	secret, err := utils.DecryptForOrg(orgID, cred.EncryptedData, cred.IV, cred.Tag, db)
	if err != nil {
		return setDomainFailed(db, d, fmt.Errorf("decrypt: %w", err))
	}
	var awsCred dto.AWSCredential
	if err := jsonUnmarshalStrict([]byte(secret), &awsCred); err != nil {
		return setDomainFailed(db, d, fmt.Errorf("secret decode: %w", err))
	}

	// 3) Client
	r53c, _, err := newRoute53Client(ctx, awsCred)
	if err != nil {
		return setDomainFailed(db, d, err)
	}

	// 4) Backfill zone id if missing
	zoneID := strings.TrimSpace(d.ZoneID)
	if zoneID == "" {
		zid, err := findHostedZoneID(ctx, r53c, d.DomainName)
		if err != nil {
			return setDomainFailed(db, d, fmt.Errorf("discover zone id: %w", err))
		}
		zoneID = zid
		d.ZoneID = zoneID
	}

	// 5) Sanity: can fetch zone
	if _, err := r53c.GetHostedZone(ctx, &r53.GetHostedZoneInput{Id: aws.String(zoneID)}); err != nil {
		return setDomainFailed(db, d, fmt.Errorf("get hosted zone: %w", err))
	}

	// 6) Mark ready
	d.Status = "ready"
	d.LastError = ""
	if err := db.Save(d).Error; err != nil {
		return err
	}
	return nil
}

func setDomainFailed(db *gorm.DB, d *models.Domain, cause error) error {
	d.Status = "failed"
	d.LastError = truncateErr(cause.Error())
	_ = db.Save(d).Error
	return cause
}

/************* record processing *************/

func processPendingRecordsForDomain(ctx context.Context, db *gorm.DB, d *models.Domain, max int) (int, error) {
	orgID := d.OrganizationID

	// reload credential
	var cred models.Credential
	if err := db.Where("id = ? AND organization_id = ?", d.CredentialID, orgID).First(&cred).Error; err != nil {
		return 0, err
	}

	secret, err := utils.DecryptForOrg(orgID, cred.EncryptedData, cred.IV, cred.Tag, db)
	if err != nil {
		return 0, err
	}

	var awsCred dto.AWSCredential
	if err := jsonUnmarshalStrict([]byte(secret), &awsCred); err != nil {
		return 0, err
	}
	r53c, _, err := newRoute53Client(ctx, awsCred)
	if err != nil {
		return 0, err
	}

	var records []models.RecordSet
	if err := db.
		Where("domain_id = ? AND status = ?", d.ID, "pending").
		Order("created_at ASC").
		Limit(max).
		Find(&records).Error; err != nil {
		return 0, err
	}

	applied := 0
	for i := range records {
		if err := applyRecord(ctx, db, r53c, d, &records[i]); err != nil {
			log.Error().Err(err).Str("rr", records[i].Name).Msg("[dns] apply record failed")
			_ = setRecordFailed(db, &records[i], err)
			continue
		}
		applied++
	}
	return applied, nil
}

// core write + ownership + external-dns hardening

func applyRecord(ctx context.Context, db *gorm.DB, r53c *r53.Client, d *models.Domain, r *models.RecordSet) error {
	zoneID := strings.TrimSpace(d.ZoneID)
	if zoneID == "" {
		return errors.New("domain has no zone_id")
	}

	rt := strings.ToUpper(r.Type)

	// FQDN & marker
	fq := recordFQDN(r.Name, d.DomainName) // ends with "."
	mname := markerName(fq)
	expected := buildMarkerValue(d.OrganizationID.String(), r.ID.String(), r.Fingerprint)

	// ---- ExternalDNS preflight ----
	extOwned, err := hasExternalDNSOwnership(ctx, r53c, zoneID, fq, rt)
	if err != nil {
		return fmt.Errorf("external_dns_lookup: %w", err)
	}
	if extOwned {
		r.Owner = "external"
		_ = db.Save(r).Error
		return fmt.Errorf("ownership_conflict: external-dns claims %s; refusing to modify", strings.TrimSuffix(fq, "."))
	}

	// ---- Autoglue ownership preflight via _autoglue.<fqdn> TXT ----
	markerVals, err := getMarkerTXTValues(ctx, r53c, zoneID, mname)
	if err != nil {
		return fmt.Errorf("marker lookup: %w", err)
	}
	hasForeignOwner := false
	hasOurExact := false
	for _, v := range markerVals {
		mk, ok := parseMarkerValue(v)
		if !ok {
			continue
		}
		switch {
		case mk.Org == d.OrganizationID.String() && mk.Rec == r.ID.String() && mk.Fp == shortFP(r.Fingerprint):
			hasOurExact = true
		case mk.Org != d.OrganizationID.String() || mk.Rec != r.ID.String():
			hasForeignOwner = true
		}
	}
	if hasForeignOwner {
		r.Owner = "external"
		_ = db.Save(r).Error
		return fmt.Errorf("ownership_conflict: marker for %s is owned by another controller; refusing to modify", strings.TrimSuffix(fq, "."))
	}

	// Build RR change (UPSERT)
	rrChange := r53types.Change{
		Action: r53types.ChangeActionUpsert,
		ResourceRecordSet: &r53types.ResourceRecordSet{
			Name: aws.String(fq),
			Type: r53types.RRType(rt),
		},
	}

	// Decode user values
	var userVals []string
	if len(r.Values) > 0 {
		if err := jsonUnmarshalStrict([]byte(r.Values), &userVals); err != nil {
			return fmt.Errorf("values decode: %w", err)
		}
	}

	// Quote TXT values as required by Route53
	recs := make([]r53types.ResourceRecord, 0, len(userVals))
	for _, v := range userVals {
		v = strings.TrimSpace(v)
		if rt == "TXT" && !(strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) {
			v = strconv.Quote(v)
		}
		recs = append(recs, r53types.ResourceRecord{Value: aws.String(v)})
	}
	rrChange.ResourceRecordSet.ResourceRecords = recs
	if r.TTL != nil {
		ttl := int64(*r.TTL)
		rrChange.ResourceRecordSet.TTL = aws.Int64(ttl)
	}

	// Build marker TXT change (UPSERT)
	markerChange := r53types.Change{
		Action: r53types.ChangeActionUpsert,
		ResourceRecordSet: &r53types.ResourceRecordSet{
			Name: aws.String(mname),
			Type: r53types.RRTypeTxt,
			TTL:  aws.Int64(300),
			ResourceRecords: []r53types.ResourceRecord{
				{Value: aws.String(strconv.Quote(expected))},
			},
		},
	}

	// Build external-dns poison TXT changes
	poisonChanges := buildExternalDNSPoisonTXTChanges(fq, rt)

	// Apply all in one batch (atomic-ish)
	changes := []r53types.Change{rrChange, markerChange}
	changes = append(changes, poisonChanges...)

	_, err = r53c.ChangeResourceRecordSets(ctx, &r53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch:  &r53types.ChangeBatch{Changes: changes},
	})
	if err != nil {
		return err
	}

	// Success → mark ready & ownership
	r.Status = "ready"
	r.LastError = ""
	r.Owner = "autoglue"
	if err := db.Save(r).Error; err != nil {
		return err
	}
	_ = hasOurExact // could be used to skip marker write in future
	return nil
}

func setRecordFailed(db *gorm.DB, r *models.RecordSet, cause error) error {
	msg := truncateErr(cause.Error())
	r.Status = "failed"
	r.LastError = msg
	// classify ownership on conflict
	if strings.HasPrefix(msg, "ownership_conflict") {
		r.Owner = "external"
	} else if r.Owner == "" || r.Owner == "unknown" {
		r.Owner = "unknown"
	}
	_ = db.Save(r).Error
	return cause
}

/************* AWS helpers *************/

func newRoute53Client(ctx context.Context, cred dto.AWSCredential) (*r53.Client, *aws.Config, error) {
	// Route53 is global, but config still wants a region
	region := strings.TrimSpace(cred.Region)
	if region == "" {
		region = "us-east-1"
	}
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cred.AccessKeyID, cred.SecretAccessKey, "",
		)),
	)
	if err != nil {
		return nil, nil, err
	}
	return r53.NewFromConfig(cfg), &cfg, nil
}

func findHostedZoneID(ctx context.Context, c *r53.Client, domain string) (string, error) {
	d := normalizeDomain(domain)
	out, err := c.ListHostedZonesByName(ctx, &r53.ListHostedZonesByNameInput{
		DNSName: aws.String(d),
	})
	if err != nil {
		return "", err
	}
	for _, hz := range out.HostedZones {
		if strings.TrimSuffix(aws.ToString(hz.Name), ".") == d {
			return trimZoneID(aws.ToString(hz.Id)), nil
		}
	}
	return "", fmt.Errorf("hosted zone not found for %q", d)
}

func trimZoneID(id string) string {
	return strings.TrimPrefix(id, "/hostedzone/")
}

func normalizeDomain(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return strings.TrimSuffix(s, ".")
}

func recordFQDN(name, domain string) string {
	name = strings.TrimSpace(name)
	if name == "" || name == "@" {
		return normalizeDomain(domain) + "."
	}
	if strings.HasSuffix(name, ".") {
		return name
	}
	return fmt.Sprintf("%s.%s.", name, normalizeDomain(domain))
}

/************* TXT marker / external-dns helpers *************/

func markerName(fqdn string) string {
	trimmed := strings.TrimSuffix(fqdn, ".")
	return "_autoglue." + trimmed + "."
}

func shortFP(full string) string {
	if len(full) > 16 {
		return full[:16]
	}
	return full
}

func buildMarkerValue(orgID, recID, fp string) string {
	return "v=ag1 org=" + orgID + " rec=" + recID + " fp=" + shortFP(fp)
}

func parseMarkerValue(s string) (ownershipMarker, bool) {
	out := ownershipMarker{}
	fields := strings.Fields(s)
	if len(fields) < 4 {
		return out, false
	}
	kv := map[string]string{}
	for _, f := range fields {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			kv[parts[0]] = parts[1]
		}
	}
	if kv["v"] == "" || kv["org"] == "" || kv["rec"] == "" || kv["fp"] == "" {
		return out, false
	}
	out.Ver, out.Org, out.Rec, out.Fp = kv["v"], kv["org"], kv["rec"], kv["fp"]
	return out, true
}

func getMarkerTXTValues(ctx context.Context, c *r53.Client, zoneID, marker string) ([]string, error) {
	return getTXTValues(ctx, c, zoneID, marker)
}

// generic TXT fetcher
func getTXTValues(ctx context.Context, c *r53.Client, zoneID, name string) ([]string, error) {
	out, err := c.ListResourceRecordSets(ctx, &r53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(name),
		StartRecordType: r53types.RRTypeTxt,
		MaxItems:        aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.ResourceRecordSets) == 0 {
		return nil, nil
	}
	rrset := out.ResourceRecordSets[0]
	if aws.ToString(rrset.Name) != name || rrset.Type != r53types.RRTypeTxt {
		return nil, nil
	}
	vals := make([]string, 0, len(rrset.ResourceRecords))
	for _, rr := range rrset.ResourceRecords {
		vals = append(vals, aws.ToString(rr.Value))
	}
	return vals, nil
}

// detect external-dns-style ownership for this fqdn/type
func hasExternalDNSOwnership(ctx context.Context, c *r53.Client, zoneID, fqdn, rrType string) (bool, error) {
	base := strings.TrimSuffix(fqdn, ".")
	candidates := []string{
		// with txtPrefix=extdns-, external-dns writes both:
		// extdns-<fqdn> and extdns-<rrtype-lc>-<fqdn>
		"extdns-" + base + ".",
		"extdns-" + strings.ToLower(rrType) + "-" + base + ".",
	}
	for _, name := range candidates {
		vals, err := getTXTValues(ctx, c, zoneID, name)
		if err != nil {
			return false, err
		}
		for _, raw := range vals {
			v := strings.TrimSpace(raw)
			// strip surrounding quotes if present
			if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
				if unq, err := strconv.Unquote(v); err == nil {
					v = unq
				}
			}
			meta := parseExternalDNSMeta(v)
			if meta == nil {
				continue
			}
			if meta["heritage"] == "external-dns" &&
				meta["external-dns/owner"] != "" &&
				meta["external-dns/owner"] != externalDNSPoisonOwner {
				return true, nil
			}
		}
	}
	return false, nil
}

// parseExternalDNSMeta parses the comma-separated external-dns TXT format into a small map.
func parseExternalDNSMeta(v string) map[string]string {
	parts := strings.Split(v, ",")
	if len(parts) == 0 {
		return nil
	}
	meta := make(map[string]string, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		meta[kv[0]] = kv[1]
	}
	if len(meta) == 0 {
		return nil
	}
	return meta
}

// build poison TXT records so external-dns thinks some *other* owner manages this
func buildExternalDNSPoisonTXTChanges(fqdn, rrType string) []r53types.Change {
	base := strings.TrimSuffix(fqdn, ".")
	names := []string{
		"extdns-" + base + ".",
		"extdns-" + strings.ToLower(rrType) + "-" + base + ".",
	}
	val := strconv.Quote(externalDNSPoisonValue)
	changes := make([]r53types.Change, 0, len(names))
	for _, n := range names {
		changes = append(changes, r53types.Change{
			Action: r53types.ChangeActionUpsert,
			ResourceRecordSet: &r53types.ResourceRecordSet{
				Name: aws.String(n),
				Type: r53types.RRTypeTxt,
				TTL:  aws.Int64(300),
				ResourceRecords: []r53types.ResourceRecord{
					{Value: aws.String(val)},
				},
			},
		})
	}
	return changes
}

/************* misc utils *************/

func truncateErr(s string) string {
	const max = 2000
	if len(s) > max {
		return s[:max]
	}
	return s
}

// Strict unmarshal that treats "null" -> zero value correctly.
func jsonUnmarshalStrict(b []byte, dst any) error {
	if len(b) == 0 {
		return errors.New("empty json")
	}
	return json.Unmarshal(b, dst)
}
