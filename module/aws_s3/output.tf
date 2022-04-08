output "tag_region" {
  value = aws_s3_bucket.demo_bucket.tags.region
}
output "tag_deployment" {
  value = aws_s3_bucket.demo_bucket.tags.deployment
}
output "tag_enviornment" {
  value = aws_s3_bucket.demo_bucket.tags.enviornment
}
output "bucket_id" {
  value = aws_s3_bucket.demo_bucket.id
}
output "bucket_name" {
  value = aws_s3_bucket.demo_bucket.bucket
}