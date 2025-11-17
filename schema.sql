-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "jobs" table
CREATE TABLE "public"."jobs" (
  "id" character varying NOT NULL,
  "queue_name" character varying NOT NULL,
  "status" character varying NOT NULL,
  "arguments" jsonb NOT NULL DEFAULT '{}',
  "result" jsonb NOT NULL DEFAULT '{}',
  "last_error" character varying NULL,
  "retry_count" bigint NOT NULL DEFAULT 0,
  "max_retry" bigint NOT NULL DEFAULT 0,
  "retry_interval" bigint NOT NULL DEFAULT 0,
  "scheduled_at" timestamptz NULL DEFAULT now(),
  "started_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_jobs_scheduled_at" to table: "jobs"
CREATE INDEX "idx_jobs_scheduled_at" ON "public"."jobs" ("scheduled_at");
-- Create index "idx_jobs_started_at" to table: "jobs"
CREATE INDEX "idx_jobs_started_at" ON "public"."jobs" ("started_at");
-- Create "api_keys" table
CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL DEFAULT '',
  "key_hash" text NOT NULL,
  "scope" text NOT NULL DEFAULT '',
  "user_id" text NULL,
  "org_id" text NULL,
  "secret_hash" text NULL,
  "expires_at" timestamptz NULL,
  "revoked" boolean NOT NULL DEFAULT false,
  "prefix" text NULL,
  "last_used_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_api_keys_key_hash" to table: "api_keys"
CREATE UNIQUE INDEX "idx_api_keys_key_hash" ON "public"."api_keys" ("key_hash");
-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" text NOT NULL,
  "family_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_refresh_tokens_family_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_family_id" ON "public"."refresh_tokens" ("family_id");
-- Create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token_hash" ON "public"."refresh_tokens" ("token_hash");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "display_name" text NULL,
  "primary_email" text NULL,
  "avatar_url" text NULL,
  "is_disabled" boolean NULL,
  "is_admin" boolean NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create "signing_keys" table
CREATE TABLE "public"."signing_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "kid" text NOT NULL,
  "alg" text NOT NULL,
  "use" text NOT NULL DEFAULT 'sig',
  "is_active" boolean NOT NULL DEFAULT true,
  "public_pem" text NOT NULL,
  "private_pem" text NOT NULL,
  "not_before" timestamptz NULL,
  "expires_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "rotated_from" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_signing_keys_kid" to table: "signing_keys"
CREATE UNIQUE INDEX "idx_signing_keys_kid" ON "public"."signing_keys" ("kid");
-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "provider" text NOT NULL,
  "subject" text NOT NULL,
  "email" text NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "profile" jsonb NOT NULL DEFAULT '{}',
  "secret_hash" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_accounts_user_id" to table: "accounts"
CREATE INDEX "idx_accounts_user_id" ON "public"."accounts" ("user_id");
-- Create "organizations" table
CREATE TABLE "public"."organizations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "domain" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_organizations_domain" to table: "organizations"
CREATE INDEX "idx_organizations_domain" ON "public"."organizations" ("domain");
-- Create "annotations" table
CREATE TABLE "public"."annotations" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_annotations_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_annotations_organization_id" to table: "annotations"
CREATE INDEX "idx_annotations_organization_id" ON "public"."annotations" ("organization_id");
-- Create "credentials" table
CREATE TABLE "public"."credentials" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "provider" character varying(50) NOT NULL,
  "kind" character varying(50) NOT NULL,
  "scope_kind" character varying(20) NOT NULL,
  "scope" jsonb NOT NULL DEFAULT '{}',
  "scope_fingerprint" character(64) NOT NULL,
  "schema_version" bigint NOT NULL DEFAULT 1,
  "name" character varying(100) NOT NULL DEFAULT '',
  "scope_version" bigint NOT NULL DEFAULT 1,
  "account_id" character varying(32) NULL,
  "region" character varying(32) NULL,
  "encrypted_data" text NOT NULL,
  "iv" text NOT NULL,
  "tag" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_credentials_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_credentials_organization_id" to table: "credentials"
CREATE INDEX "idx_credentials_organization_id" ON "public"."credentials" ("organization_id");
-- Create index "idx_credentials_scope_fingerprint" to table: "credentials"
CREATE INDEX "idx_credentials_scope_fingerprint" ON "public"."credentials" ("scope_fingerprint");
-- Create index "idx_kind_scope" to table: "credentials"
CREATE INDEX "idx_kind_scope" ON "public"."credentials" ("kind", "scope");
-- Create index "idx_provider_kind" to table: "credentials"
CREATE INDEX "idx_provider_kind" ON "public"."credentials" ("provider", "kind");
-- Create index "uniq_org_provider_scopekind_scope" to table: "credentials"
CREATE UNIQUE INDEX "uniq_org_provider_scopekind_scope" ON "public"."credentials" ("organization_id", "provider", "scope_kind", "scope_fingerprint");
-- Create "backups" table
CREATE TABLE "public"."backups" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "enabled" boolean NOT NULL DEFAULT false,
  "credential_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_backups_credential" FOREIGN KEY ("credential_id") REFERENCES "public"."credentials" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_backups_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_backups_organization_id" to table: "backups"
CREATE INDEX "idx_backups_organization_id" ON "public"."backups" ("organization_id");
-- Create index "uniq_org_credential" to table: "backups"
CREATE UNIQUE INDEX "uniq_org_credential" ON "public"."backups" ("organization_id", "credential_id");
-- Create "load_balancers" table
CREATE TABLE "public"."load_balancers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "name" text NOT NULL,
  "kind" text NOT NULL,
  "public_ip_address" text NOT NULL,
  "private_ip_address" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_load_balancers_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_load_balancers_organization_id" to table: "load_balancers"
CREATE INDEX "idx_load_balancers_organization_id" ON "public"."load_balancers" ("organization_id");
-- Create "ssh_keys" table
CREATE TABLE "public"."ssh_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" text NOT NULL,
  "public_key" text NOT NULL,
  "encrypted_private_key" text NOT NULL,
  "private_iv" text NOT NULL,
  "private_tag" text NOT NULL,
  "fingerprint" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ssh_keys_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_ssh_keys_fingerprint" to table: "ssh_keys"
CREATE INDEX "idx_ssh_keys_fingerprint" ON "public"."ssh_keys" ("fingerprint");
-- Create index "idx_ssh_keys_organization_id" to table: "ssh_keys"
CREATE INDEX "idx_ssh_keys_organization_id" ON "public"."ssh_keys" ("organization_id");
-- Create "servers" table
CREATE TABLE "public"."servers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "hostname" text NULL,
  "public_ip_address" text NULL,
  "private_ip_address" text NOT NULL,
  "ssh_user" text NOT NULL,
  "ssh_key_id" uuid NOT NULL,
  "role" text NOT NULL,
  "status" text NULL DEFAULT 'pending',
  "ssh_host_key" text NULL,
  "ssh_host_key_algo" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_servers_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_servers_ssh_key" FOREIGN KEY ("ssh_key_id") REFERENCES "public"."ssh_keys" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "domains" table
CREATE TABLE "public"."domains" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "domain_name" character varying(253) NOT NULL,
  "zone_id" character varying(128) NOT NULL DEFAULT '',
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "last_error" text NOT NULL DEFAULT '',
  "credential_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_domains_credential" FOREIGN KEY ("credential_id") REFERENCES "public"."credentials" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_domains_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_domains_organization_id" to table: "domains"
CREATE INDEX "idx_domains_organization_id" ON "public"."domains" ("organization_id");
-- Create index "uniq_org_domain" to table: "domains"
CREATE UNIQUE INDEX "uniq_org_domain" ON "public"."domains" ("organization_id", "domain_name");
-- Create "record_sets" table
CREATE TABLE "public"."record_sets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "domain_id" uuid NOT NULL,
  "name" character varying(253) NOT NULL,
  "type" character varying(10) NOT NULL,
  "ttl" bigint NULL,
  "values" jsonb NOT NULL DEFAULT '[]',
  "fingerprint" character(64) NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "owner" character varying(16) NOT NULL DEFAULT 'unknown',
  "last_error" text NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_record_sets_domain" FOREIGN KEY ("domain_id") REFERENCES "public"."domains" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_record_sets_domain_id" to table: "record_sets"
CREATE INDEX "idx_record_sets_domain_id" ON "public"."record_sets" ("domain_id");
-- Create index "idx_record_sets_fingerprint" to table: "record_sets"
CREATE INDEX "idx_record_sets_fingerprint" ON "public"."record_sets" ("fingerprint");
-- Create index "idx_record_sets_type" to table: "record_sets"
CREATE INDEX "idx_record_sets_type" ON "public"."record_sets" ("type");
-- Create "clusters" table
CREATE TABLE "public"."clusters" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "name" text NOT NULL,
  "provider" text NULL,
  "region" text NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pre_pending',
  "last_error" text NOT NULL DEFAULT '',
  "captain_domain_id" uuid NULL,
  "control_plane_record_set_id" uuid NULL,
  "apps_load_balancer_id" uuid NULL,
  "glue_ops_load_balancer_id" uuid NULL,
  "bastion_server_id" uuid NULL,
  "random_token" text NULL,
  "certificate_key" text NULL,
  "encrypted_kubeconfig" text NULL,
  "kube_iv" text NULL,
  "kube_tag" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clusters_apps_load_balancer" FOREIGN KEY ("apps_load_balancer_id") REFERENCES "public"."load_balancers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_clusters_bastion_server" FOREIGN KEY ("bastion_server_id") REFERENCES "public"."servers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_clusters_captain_domain" FOREIGN KEY ("captain_domain_id") REFERENCES "public"."domains" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_clusters_control_plane_record_set" FOREIGN KEY ("control_plane_record_set_id") REFERENCES "public"."record_sets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_clusters_glue_ops_load_balancer" FOREIGN KEY ("glue_ops_load_balancer_id") REFERENCES "public"."load_balancers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_clusters_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "node_pools" table
CREATE TABLE "public"."node_pools" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" text NOT NULL,
  "role" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_node_pools_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_node_pools_organization_id" to table: "node_pools"
CREATE INDEX "idx_node_pools_organization_id" ON "public"."node_pools" ("organization_id");
-- Create "cluster_node_pools" table
CREATE TABLE "public"."cluster_node_pools" (
  "node_pool_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "cluster_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("node_pool_id", "cluster_id"),
  CONSTRAINT "fk_cluster_node_pools_cluster" FOREIGN KEY ("cluster_id") REFERENCES "public"."clusters" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_cluster_node_pools_node_pool" FOREIGN KEY ("node_pool_id") REFERENCES "public"."node_pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "labels" table
CREATE TABLE "public"."labels" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_labels_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_labels_organization_id" to table: "labels"
CREATE INDEX "idx_labels_organization_id" ON "public"."labels" ("organization_id");
-- Create "memberships" table
CREATE TABLE "public"."memberships" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "organization_id" uuid NOT NULL,
  "role" text NOT NULL DEFAULT 'member',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_memberships_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_memberships_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_memberships_organization_id" to table: "memberships"
CREATE INDEX "idx_memberships_organization_id" ON "public"."memberships" ("organization_id");
-- Create index "idx_memberships_user_id" to table: "memberships"
CREATE INDEX "idx_memberships_user_id" ON "public"."memberships" ("user_id");
-- Create "node_annotations" table
CREATE TABLE "public"."node_annotations" (
  "node_pool_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "annotation_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("node_pool_id", "annotation_id"),
  CONSTRAINT "fk_node_annotations_annotation" FOREIGN KEY ("annotation_id") REFERENCES "public"."annotations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_node_annotations_node_pool" FOREIGN KEY ("node_pool_id") REFERENCES "public"."node_pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "node_labels" table
CREATE TABLE "public"."node_labels" (
  "node_pool_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "label_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("node_pool_id", "label_id"),
  CONSTRAINT "fk_node_labels_label" FOREIGN KEY ("label_id") REFERENCES "public"."labels" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_node_labels_node_pool" FOREIGN KEY ("node_pool_id") REFERENCES "public"."node_pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "node_servers" table
CREATE TABLE "public"."node_servers" (
  "server_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "node_pool_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("server_id", "node_pool_id"),
  CONSTRAINT "fk_node_servers_node_pool" FOREIGN KEY ("node_pool_id") REFERENCES "public"."node_pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_node_servers_server" FOREIGN KEY ("server_id") REFERENCES "public"."servers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "taints" table
CREATE TABLE "public"."taints" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  "effect" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_taints_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "node_taints" table
CREATE TABLE "public"."node_taints" (
  "taint_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "node_pool_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("taint_id", "node_pool_id"),
  CONSTRAINT "fk_node_taints_node_pool" FOREIGN KEY ("node_pool_id") REFERENCES "public"."node_pools" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_node_taints_taint" FOREIGN KEY ("taint_id") REFERENCES "public"."taints" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "master_keys" table
CREATE TABLE "public"."master_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "key" text NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create "organization_keys" table
CREATE TABLE "public"."organization_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "organization_id" uuid NOT NULL,
  "master_key_id" uuid NOT NULL,
  "encrypted_key" text NOT NULL,
  "iv" text NOT NULL,
  "tag" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_organization_keys_master_key" FOREIGN KEY ("master_key_id") REFERENCES "public"."master_keys" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_organization_keys_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "user_emails" table
CREATE TABLE "public"."user_emails" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "email" text NOT NULL,
  "is_verified" boolean NOT NULL DEFAULT false,
  "is_primary" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_emails_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_emails_user_id" to table: "user_emails"
CREATE INDEX "idx_user_emails_user_id" ON "public"."user_emails" ("user_id");
