package bg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DbBackupArgs struct {
	IntervalS int `json:"interval_seconds,omitempty"`
}

type s3Scope struct {
	Service string `json:"service"`
	Region  string `json:"region"`
}

type encAWS struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func DbBackupWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := DbBackupArgs{IntervalS: 3600}
		_ = j.ParseArguments(&args)

		if args.IntervalS <= 0 {
			args.IntervalS = 3600
		}

		if err := DbBackup(ctx, db); err != nil {
			return nil, err
		}

		queue := j.QueueName
		if strings.TrimSpace(queue) == "" {
			queue = "db_backup_s3"
		}

		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)

		payload := DbBackupArgs{}

		opts := []archer.FnOptions{
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		}

		if _, err := jobs.Enqueue(ctx, uuid.NewString(), queue, payload, opts...); err != nil {
			log.Error().Err(err).Str("queue", queue).Time("next", next).Msg("failed to enqueue next db backup")
		} else {
			log.Info().Str("queue", queue).Time("next", next).Msg("scheduled next db backup")
		}
		return nil, nil
	}
}

func DbBackup(ctx context.Context, db *gorm.DB) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	cred, sc, err := loadS3Credential(ctx, db)
	if err != nil {
		return fmt.Errorf("load credential: %w", err)
	}

	ak, sk, err := decryptAwsAccessKeys(ctx, db, cred)
	if err != nil {
		return fmt.Errorf("decrypt aws keys: %w", err)
	}

	region := sc.Region

	if strings.TrimSpace(region) == "" {
		region = cred.Region
		if strings.TrimSpace(region) == "" {
			region = "us-west-1"
		}
	}

	bucket := strings.ToLower(fmt.Sprintf("%s-autoglue-backups-%s", cred.OrganizationID, region))

	s3cli, err := makeS3Client(ctx, ak, sk, region)
	if err != nil {
		return err
	}

	if err := ensureBucket(ctx, s3cli, bucket, region); err != nil {
		return fmt.Errorf("ensure bucket: %w", err)
	}

	tmpDir := os.TempDir()
	now := time.Now().UTC()
	key := fmt.Sprintf("%04d/%02d/%02d/backup-%02d.sql", now.Year(), now.Month(), now.Day(), now.Hour())
	outPath := filepath.Join(tmpDir, "autoglue-backup-"+now.Format("20060102T150405Z")+".sql")

	if err := runPgDump(ctx, cfg.DbURL, outPath); err != nil {
		return fmt.Errorf("pg_dump: %w", err)
	}
	defer os.Remove(outPath)

	if err := uploadFileToS3(ctx, s3cli, bucket, key, outPath); err != nil {
		return fmt.Errorf("s3 upload: %w", err)
	}

	log.Info().Str("bucket", bucket).Str("key", key).Msg("backup uploaded")

	return nil
}

// --- Helpers

func loadS3Credential(ctx context.Context, db *gorm.DB) (models.Credential, s3Scope, error) {
	var c models.Credential
	err := db.
		WithContext(ctx).
		Where("provider = ? AND kind = ? AND scope_kind = ?", "aws", "aws_access_key", "service").
		Where("scope ->> 'service' = ?", "s3").
		Order("created_at DESC").
		First(&c).Error
	if err != nil {
		return models.Credential{}, s3Scope{}, fmt.Errorf("load credential: %w", err)
	}

	var sc s3Scope
	_ = json.Unmarshal(c.Scope, &sc)
	return c, sc, nil
}

func decryptAwsAccessKeys(ctx context.Context, db *gorm.DB, c models.Credential) (string, string, error) {
	plain, err := utils.DecryptForOrg(c.OrganizationID, c.EncryptedData, c.IV, c.Tag, db)
	if err != nil {
		return "", "", err
	}

	var payload encAWS
	if err := json.Unmarshal([]byte(plain), &payload); err != nil {
		return "", "", fmt.Errorf("parse decrypted payload: %w", err)
	}

	if payload.AccessKeyID == "" || payload.SecretAccessKey == "" {
		return "", "", errors.New("decrypted payload missing keys")
	}
	return payload.AccessKeyID, payload.SecretAccessKey, nil
}

func makeS3Client(ctx context.Context, accessKey, secret, region string) (*s3.Client, error) {
	staticCredentialsProvider := credentials.NewStaticCredentialsProvider(accessKey, secret, "")
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithCredentialsProvider(staticCredentialsProvider), awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("aws config: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func ensureBucket(ctx context.Context, s3cli *s3.Client, bucket, region string) error {
	_, err := s3cli.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err == nil {
		return nil
	}

	if out, err := s3cli.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: aws.String(bucket)}); err == nil {
		existing := string(out.LocationConstraint)
		if existing == "" {
			existing = "us-east-1"
		}
		if existing != region {
			return fmt.Errorf("bucket %q already exists in region %q (requested %q)", bucket, existing, region)
		}
	}

	// Create; LocationConstraint except us-east-1
	in := &s3.CreateBucketInput{Bucket: aws.String(bucket)}
	if region != "us-east-1" {
		in.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(region),
		}
	}
	if _, err := s3cli.CreateBucket(ctx, in); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}

	// default SSE (best-effort)
	_, _ = s3cli.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucket),
		ServerSideEncryptionConfiguration: &s3types.ServerSideEncryptionConfiguration{
			Rules: []s3types.ServerSideEncryptionRule{
				{ApplyServerSideEncryptionByDefault: &s3types.ServerSideEncryptionByDefault{
					SSEAlgorithm: s3types.ServerSideEncryptionAes256,
				}},
			},
		},
	})
	return nil
}

func runPgDump(ctx context.Context, dbURL, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	args := []string{
		"--no-owner",
		"--no-privileges",
		"--format=plain",
		"--file", outPath,
		dbURL,
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %v | %s", err, stderr.String())
	}

	return nil
}

func uploadFileToS3(ctx context.Context, s3cli *s3.Client, bucket, key, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	info, _ := f.Stat()
	_, err = s3cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		Body:                 f,
		ContentLength:        aws.Int64(info.Size()),
		ContentType:          aws.String(mime.TypeByExtension(filepath.Ext(path))),
		ServerSideEncryption: s3types.ServerSideEncryptionAes256,
	})

	return err
}
