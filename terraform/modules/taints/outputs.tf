output "ids" {
  description = "Map of taint IDs by key."
  value       = { for k, r in autoglue_taint.this : k => r.id }
}
