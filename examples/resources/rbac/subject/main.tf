resource "chronicle_rbac_subject" "subject" {
  name  = "name"
  type  = "SUBJECT_TYPE_ANALYST"
  roles = ["Editor"]
}
